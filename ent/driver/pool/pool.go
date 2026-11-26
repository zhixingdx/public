// Package pool 提供所有 SQL driver 共用的有界连接池配置。
package pool

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zhixingdx/oipivot-api/config"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	// DefaultMaxOpenConns 是单进程数据库池的默认最大连接数。
	DefaultMaxOpenConns = 64
	// DefaultMaxIdleConns 是单进程数据库池的默认最大空闲连接数。
	DefaultMaxIdleConns = 32
	// DefaultConnMaxLifetime 是连接被轮换前的默认最大复用时长。
	DefaultConnMaxLifetime = 30 * time.Minute
	// DefaultConnMaxIdleTime 是空闲连接的默认最大保留时长。
	DefaultConnMaxIdleTime = 5 * time.Minute

	maxOpenConnsLimit = 2000
	maxConnLifetime   = 24 * time.Hour
	maxConnIdleTime   = time.Hour
)

// Settings 是已应用默认值并通过边界校验的 SQL 连接池参数。
type Settings struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// Resolve 将公共 Database 配置解析为安全的连接池参数。
func Resolve(database *config.Database) (Settings, error) {
	if database == nil {
		return Settings{}, fmt.Errorf("database 配置不能为空")
	}
	maxOpen := int(database.MaxOpenConns)
	if maxOpen == 0 {
		maxOpen = DefaultMaxOpenConns
	}
	if maxOpen < 1 || maxOpen > maxOpenConnsLimit {
		return Settings{}, fmt.Errorf("database.maxOpenConns 必须在 1 到 %d 之间", maxOpenConnsLimit)
	}
	maxIdle := int(database.MaxIdleConns)
	if maxIdle == 0 {
		maxIdle = min(DefaultMaxIdleConns, maxOpen)
	}
	if maxIdle < 1 || maxIdle > maxOpen {
		return Settings{}, fmt.Errorf("database.maxIdleConns 必须在 1 到 maxOpenConns 之间")
	}
	lifetime, err := resolveDuration(
		database.ConnMaxLifetime,
		DefaultConnMaxLifetime,
		maxConnLifetime,
		"database.connMaxLifetime",
	)
	if err != nil {
		return Settings{}, err
	}
	idleTime, err := resolveDuration(
		database.ConnMaxIdleTime,
		DefaultConnMaxIdleTime,
		maxConnIdleTime,
		"database.connMaxIdleTime",
	)
	if err != nil {
		return Settings{}, err
	}
	if idleTime > lifetime {
		return Settings{}, fmt.Errorf("database.connMaxIdleTime 不能超过 connMaxLifetime")
	}
	return Settings{
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: lifetime,
		ConnMaxIdleTime: idleTime,
	}, nil
}

// Apply 对长生命周期 sql.DB 应用统一的有界连接预算。
func Apply(db *sql.DB, database *config.Database) error {
	if db == nil {
		return fmt.Errorf("sql.DB 不能为空")
	}
	settings, err := Resolve(database)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(settings.MaxOpenConns)
	db.SetMaxIdleConns(settings.MaxIdleConns)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime)
	db.SetConnMaxIdleTime(settings.ConnMaxIdleTime)
	return nil
}

// resolveDuration 解析可选 protobuf Duration，并拒绝零值、负数与不合理极值。
func resolveDuration(
	value *durationpb.Duration,
	fallback time.Duration,
	maximum time.Duration,
	name string,
) (time.Duration, error) {
	if value == nil {
		return fallback, nil
	}
	if err := value.CheckValid(); err != nil {
		return 0, fmt.Errorf("%s 非法: %w", name, err)
	}
	duration := value.AsDuration()
	if duration <= 0 || duration > maximum {
		return 0, fmt.Errorf("%s 必须在 0 到 %s 之间", name, maximum)
	}
	return duration, nil
}
