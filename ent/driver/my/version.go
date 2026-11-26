package my

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

const (
	// MinimumVersion 是 qiyun 业务数据库使用 VALUES 表值构造器的最低 MySQL 版本。
	MinimumVersion = "8.0.19"
	minimumSemver  = "v" + MinimumVersion
)

// ValidateVersion 校验 MySQL 核心版本，同时允许云厂商或商业发行版追加连字符后缀。
func ValidateVersion(raw string) error {
	version := strings.TrimSpace(raw)
	if strings.Contains(strings.ToLower(version), "mariadb") {
		return fmt.Errorf("qiyun 业务数据库要求 MySQL >= %s，不支持 MariaDB 版本 %s", MinimumVersion, raw)
	}
	// VERSION() 的连字符后缀是发行渠道，不按 SemVer 预发布版本参与能力比较。
	if suffix := strings.IndexAny(version, "-+"); suffix >= 0 {
		version = version[:suffix]
	}
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	if !semver.IsValid(version) {
		return fmt.Errorf("无法识别 MySQL 版本 %q", raw)
	}
	if semver.Compare(version, minimumSemver) < 0 {
		return fmt.Errorf("qiyun 业务数据库要求 MySQL >= %s，实际为 %s", MinimumVersion, raw)
	}
	return nil
}

// ValidateServer 同时校验版本号和 VALUES ROW 能力，防止兼容数据库仅伪装版本但缺少目标语法。
func ValidateServer(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("MySQL 连接池不能为空")
	}
	var version string
	if err := db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		return fmt.Errorf("读取 MySQL 版本失败: %w", err)
	}
	if err := ValidateVersion(version); err != nil {
		return err
	}
	rows, err := db.QueryContext(ctx, "SELECT `value` FROM (VALUES ROW(1)) AS `q` (`value`)")
	if err != nil {
		return fmt.Errorf("MySQL %s 缺少 VALUES ROW 表值构造器能力: %w", version, err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return fmt.Errorf("验证 MySQL VALUES ROW 能力失败: %w", err)
		}
		return fmt.Errorf("验证 MySQL VALUES ROW 能力未返回结果")
	}
	var value int
	if err := rows.Scan(&value); err != nil {
		return fmt.Errorf("读取 MySQL VALUES ROW 能力结果失败: %w", err)
	}
	if value != 1 {
		return fmt.Errorf("MySQL VALUES ROW 能力返回异常值 %d", value)
	}
	return nil
}
