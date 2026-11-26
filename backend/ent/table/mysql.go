// Package table 定义 Ent Schema 复用的物理表配置。
package table

import (
	"context"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

const (
	// MySQLCharset 是 Ent 创建 MySQL 表时统一使用的字符集。
	MySQLCharset = "utf8mb4"
	// MySQLCaseSensitiveCollation 是 Ent 创建 MySQL 表时统一使用的大小写敏感排序规则。
	MySQLCaseSensitiveCollation = "utf8mb4_0900_as_cs"
)

// MigrateOptions 为 MySQL Schema.Create 追加统一的表默认字符集和排序规则。
//
// Ent v0.14.6 会在表没有显式 annotation 时使用 utf8mb4_bin，并为 MySQL UUID
// 回退列显式设置 utf8mb4_bin。这里在迁移边界统一修正表默认值，并让未声明
// 物理类型的 UUID 列继承表默认值；这里不写入列级 collation，避免覆盖显式规则
// 或在启动时触发存量列批量转换。该选项也覆盖 Casbin 等外部 Ent Schema；
// PostgreSQL 等其他方言保持原迁移选项不变。
func MigrateOptions(driver string, options ...schema.MigrateOption) []schema.MigrateOption {
	if driver != dialect.MySQL {
		return options
	}
	return append(options, schema.WithHooks(mysqlTableDefaultsHook))
}

// mysqlTableDefaultsHook 在 Ent 比较目标 Schema 前冻结 MySQL 表默认值。
func mysqlTableDefaultsHook(next schema.Creator) schema.Creator {
	return schema.CreateFunc(func(ctx context.Context, tables ...*schema.Table) error {
		for _, table := range tables {
			if table == nil {
				continue
			}
			if table.Annotation == nil {
				table.Annotation = &entsql.Annotation{}
			}
			table.Annotation.Charset = MySQLCharset
			table.Annotation.Collation = MySQLCaseSensitiveCollation
			for _, column := range table.Columns {
				if column == nil || column.Type != field.TypeUUID {
					continue
				}
				// Ent 的 MySQL UUID 回退逻辑会在类型转换阶段强制写入
				// utf8mb4_bin。只为未声明物理类型的 UUID 指定等价 CHAR(36)，
				// 使新列继承表默认值，同时不覆盖 BINARY 或原生 UUID 等显式契约。
				if column.SchemaType == nil {
					column.SchemaType = make(map[string]string)
				}
				if column.SchemaType[dialect.MySQL] == "" {
					column.SchemaType[dialect.MySQL] = "char(36)"
				}
			}
		}
		return next.Create(ctx, tables...)
	})
}
