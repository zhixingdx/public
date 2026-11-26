// Package mixin 提供 Ent 表的公共字段。
package mixin

import (
	"database/sql"
	"database/sql/driver"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
	"go.temporal.io/api/workflowservice/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// WorkflowOptions mixin 工作流配置
type WorkflowOptions struct {
	mixin.Schema
}

// Fields _
func (WorkflowOptions) Fields() []ent.Field {
	return []ent.Field{
		field.String("workflowOptions").
			GoType(&workflowservice.StartWorkflowExecutionRequest{}).
			SchemaType(map[string]string{dialect.MySQL: "json", dialect.Postgres: "jsonb", dialect.SQLite: "json"}).
			ValueScanner(workflowRequestJSON).
			Optional().
			DefaultFunc(func() *workflowservice.StartWorkflowExecutionRequest {
				return &workflowservice.StartWorkflowExecutionRequest{}
			}).
			Comment("workflow 配置").
			Annotations(ann.Field{
				PbIndex: 224,
				Query:   ann.FieldQuery{Disabled: true},
				PB: ann.FieldPB{
					Package: "temporal.api.workflowservice.v1",
					Import:  "temporal/api/workflowservice/v1/request_response.proto",
				},
				NotPortalAdd:         true,
				NotPortalEdit:        true,
				NotPortalQuery:       true,
				NotPortalQueryReply:  true,
				NotPortalDetailReply: true,
			}),
	}
}

// Ent 的默认 JSON 编码无法回读 Protobuf oneof；数据库使用官方 ProtoJSON。
var workflowRequestJSON = field.ValueScannerFunc[*workflowservice.StartWorkflowExecutionRequest, *sql.NullString]{
	V: func(req *workflowservice.StartWorkflowExecutionRequest) (driver.Value, error) {
		if req == nil {
			return nil, nil
		}
		data, err := protojson.Marshal(req)
		return string(data), err
	},
	S: func(value *sql.NullString) (*workflowservice.StartWorkflowExecutionRequest, error) {
		if !value.Valid {
			return nil, nil
		}
		req := new(workflowservice.StartWorkflowExecutionRequest)
		if err := protojson.Unmarshal([]byte(value.String), req); err != nil {
			return nil, err
		}
		return req, nil
	},
}

// workflowSchema
type WorkflowSchema struct {
	mixin.Schema
}

// Fields _
func (WorkflowSchema) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("workflowSchema", &structpb.Struct{}).
			Optional().
			Default(&structpb.Struct{}).
			Comment("workflow Schema").
			Annotations(
				ann.Field{
					PbIndex:              225,
					Query:                ann.FieldQuery{Disabled: true},
					NotPortalAdd:         true,
					NotPortalEdit:        true,
					NotPortalQuery:       true,
					NotPortalQueryReply:  true,
					NotPortalDetailReply: true,
				},
			),
	}
}

// workflowData
type WorkflowData struct {
	mixin.Schema
}

// Fields _
func (WorkflowData) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("workflowData", &structpb.Struct{}).
			Optional().
			Default(&structpb.Struct{}).Comment("workflow 数据").
			Annotations(
				ann.Field{
					PbIndex:              226,
					Query:                ann.FieldQuery{Disabled: true},
					NotPortalAdd:         true,
					NotPortalEdit:        true,
					NotPortalQuery:       true,
					NotPortalQueryReply:  true,
					NotPortalDetailReply: true,
				},
			),
	}
}
