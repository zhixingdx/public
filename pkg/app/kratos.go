// Package app kratos.go
package app

import (
	"encoding/json"
	"fmt"

	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/file"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"
	conf "github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot/backend/pkg/logger/kratoslog"
	"google.golang.org/protobuf/encoding/protojson"
)

// DepsInit 依赖初始化 当前仅用于占位 用于当依赖注入时使用
type DepsInit struct {
}

// NewApp http.Server 支持多个 可分别支持 manage 和 portal API
func NewApp(logger log.Logger, conf *conf.Config, _ *DepsInit, gs *grpc.Server, hss ...*http.Server) *kratos.App {
	servers := []transport.Server{}
	for _, s := range hss {
		if s != nil {
			servers = append(servers, s)
		}
	}
	if gs != nil {
		servers = append(servers, gs)
	}
	if len(servers) == 0 {
		panic("请检查配置文件，至少保证启动 GRPC 或 HTTP 中的一个服务")
	}
	server := conf.Server
	return kratos.New(
		kratos.ID(server.Id),
		kratos.Name(server.Name),
		kratos.Version(server.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(log.AsSlog(logger)),
		kratos.Server(servers...),
	)
}

// UpdateDefaultConfig 根据启动参数更新配置文件
func UpdateDefaultConfig(conf *conf.Config, id string, name string, version string) *conf.Config {
	serverName := name
	if conf.Server != nil {
		if conf.Server.Id == "" {
			conf.Server.Id = id
		}
		if conf.Server.Name == "" {
			conf.Server.Name = name
		}
		serverName = conf.Server.Name
		if conf.Server.Version == "" {
			conf.Server.Version = version
		}
	}
	// 将 server.env 作为可观测性信号环境标识的唯一来源
	serverEnv := ""
	if conf.Server != nil {
		serverEnv = conf.Server.Env
	}

	if conf.Trace != nil {
		if conf.Trace.Service == "" {
			conf.Trace.Service = serverName
		}
		if conf.Trace.Env == "" && serverEnv != "" {
			conf.Trace.Env = serverEnv
		}
	}
	if conf.Metrics != nil {
		if conf.Metrics.Service == "" {
			conf.Metrics.Service = serverName
		}
		if conf.Metrics.Env == "" && serverEnv != "" {
			conf.Metrics.Env = serverEnv
		}
	}

	return conf
}

// LoadFileConf 加载文件配置
func LoadFileConf(flagconf string) (config.Config, *conf.Config) {
	c := config.New(config.WithSource(file.NewSource(flagconf)))
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Config
	if err := c.Scan(&bc); err != nil {
		_ = c.Close()
		panic(err)
	}
	// Kratos 默认忽略未知字段。严格读取监听器，以免证书配置拼错后启用明文。
	if listener := bc.GetServer().GetGrpc(); listener != nil {
		raw := c.Value("server.grpc").Load()
		data, err := json.Marshal(raw)
		if err == nil {
			err = protojson.Unmarshal(data, listener)
		}
		if err != nil {
			_ = c.Close()
			panic(fmt.Errorf("server.grpc: %w", err))
		}
		if fields, ok := raw.(map[string]any); ok {
			if _, configured := fields["mtls"]; configured && listener.Mtls == nil {
				_ = c.Close()
				panic("server.grpc.mtls 须为对象")
			}
		}
	}
	return c, &bc
}
