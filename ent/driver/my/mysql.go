// Package my ent mysql driver
package my

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect/sql"

	"entgo.io/ent/dialect"
	"github.com/zhixingdx/public/api/config"
	"github.com/zhixingdx/public/ent/driver/pool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// GetMySQLDriver 创建采用公共连接预算和可选 tracing 的 MySQL Ent driver。
func GetMySQLDriver(conf *config.Database) dialect.Driver {
	drv, err := sql.Open(conf.Driver, conf.Source)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	if err := pool.Apply(drv.DB(), conf); err != nil {
		log.Fatalf("数据库连接池配置失败: %v", err)
	}
	// 公共 driver 在所有业务服务启动时执行同一能力门禁，避免低版本实例运行到批量写才暴露语法错误。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ValidateServer(ctx, drv.DB()); err != nil {
		_ = drv.Close()
		log.Fatalf("数据库能力检查失败: %v", err)
	}
	if !conf.IsTrace {
		return dialect.Driver(drv)
	}
	return dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		tracer := otel.Tracer("ent")
		kind := trace.SpanKindServer
		_, span := tracer.Start(ctx,
			"mysql",
			trace.WithAttributes(
				attribute.String("db.statement", fmt.Sprint(i...)),
				attribute.String("db.system", "mysql"),
			),
			trace.WithSpanKind(kind),
		)
		span.End()
	})
}
