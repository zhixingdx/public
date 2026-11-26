// Package utils 通用工具包
package utils

import (
	"sync"
	"time"

	"github.com/zhixingdx/public/pkg/logger/kratoslog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// GRPCConnectionManager gRPC 连接管理器
// 统一管理所有 gRPC 连接，提供连接池和手动重连功能
// 注意：自动重连由 keepalive 提供，此管理器仅提供连接池和健康监控
type GRPCConnectionManager struct {
	connections map[string]*managedConnection // 连接映射：endpoint -> connection
	mu          sync.RWMutex                  // 读写锁
	log         *log.Helper                   // 日志
	config      *ManagerConfig                // 配置
}

// ManagerConfig 连接管理器配置
type ManagerConfig struct {
	HealthCheckInterval time.Duration // 健康监控间隔（默认 30 秒，仅检测不重连）
	ReconnectInterval   time.Duration // 重连间隔（默认 5 秒，避免连接风暴）
	EnableHealthMonitor bool          // 是否启用健康监控（默认 true，仅检测不重连）
}

// DefaultManagerConfig 默认配置
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		HealthCheckInterval: 30 * time.Second,
		ReconnectInterval:   5 * time.Second,
		EnableHealthMonitor: true, // 仅监控，不自动重连（keepalive 已提供自动重连）
	}
}

// managedConnection 管理的连接
type managedConnection struct {
	conn           *grpc.ClientConn
	endpoint       string
	lastReconnect  time.Time
	reconnectCount int
	mu             sync.RWMutex
}

// 全局连接管理器实例（单例模式）
var (
	globalConnManager     *GRPCConnectionManager
	globalConnManagerOnce sync.Once
)

// GetGlobalGRPCConnectionManager 获取全局 gRPC 连接管理器（单例）
// 用于 Client 层统一管理所有连接
// 注意：自动重连由 keepalive 提供，此管理器仅提供连接池和手动重连功能
func GetGlobalGRPCConnectionManager(logger log.Logger) *GRPCConnectionManager {
	globalConnManagerOnce.Do(func() {
		config := DefaultManagerConfig()
		globalConnManager = &GRPCConnectionManager{
			connections: make(map[string]*managedConnection),
			log:         log.NewHelper(log.With(logger, "component", "grpc-manager")),
			config:      config,
		}

		// 启动健康监控（仅检测，不自动重连）
		if config.EnableHealthMonitor {
			go globalConnManager.startHealthMonitor()
		}
	})
	return globalConnManager
}

// NewGRPCConnectionManager 创建新的 gRPC 连接管理器
// 用于 Usecase 层创建独立的连接管理器实例
// 注意：自动重连由 keepalive 提供，此管理器仅提供连接池和手动重连功能
func NewGRPCConnectionManager(logger log.Logger, config *ManagerConfig) *GRPCConnectionManager {
	if config == nil {
		config = DefaultManagerConfig()
	}

	manager := &GRPCConnectionManager{
		connections: make(map[string]*managedConnection),
		log:         log.NewHelper(log.With(logger, "component", "grpc-manager")),
		config:      config,
	}

	// 启动健康监控（仅检测，不自动重连）
	if config.EnableHealthMonitor {
		go manager.startHealthMonitor()
	}

	return manager
}

// GetConnection 获取或创建连接
// 如果连接不存在或已断开，会创建新连接（检测失败才重连）
func (m *GRPCConnectionManager) GetConnection(endpoint string, dialFunc func(string) (*grpc.ClientConn, error)) (*grpc.ClientConn, error) {
	m.mu.RLock()
	managedConn, exists := m.connections[endpoint]
	m.mu.RUnlock()

	// 如果连接存在且正常，直接返回
	if exists {
		managedConn.mu.RLock()
		conn := managedConn.conn
		state := conn.GetState()
		managedConn.mu.RUnlock()

		if state == connectivity.Ready {
			return conn, nil
		}

		// 连接异常，需要重连
		m.log.Warn("连接状态异常，准备重连",
			"endpoint", endpoint,
			"state", state.String(),
		)
	}

	// 创建新连接
	return m.createConnection(endpoint, dialFunc)
}

// createConnection 创建新连接
func (m *GRPCConnectionManager) createConnection(endpoint string, dialFunc func(string) (*grpc.ClientConn, error)) (*grpc.ClientConn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if managedConn, exists := m.connections[endpoint]; exists {
		managedConn.mu.RLock()
		conn := managedConn.conn
		state := conn.GetState()
		managedConn.mu.RUnlock()

		if state == connectivity.Ready {
			return conn, nil
		}

		// 检查重连间隔
		managedConn.mu.Lock()
		if time.Since(managedConn.lastReconnect) < m.config.ReconnectInterval {
			managedConn.mu.Unlock()
			return conn, nil // 还在重连间隔内，返回旧连接
		}
		managedConn.mu.Unlock()

		// 关闭旧连接
		_ = conn.Close()
	}

	// 创建新连接
	conn, err := dialFunc(endpoint)
	if err != nil {
		return nil, err
	}

	// 保存连接
	managedConn := &managedConnection{
		conn:           conn,
		endpoint:       endpoint,
		lastReconnect:  time.Now(),
		reconnectCount: 1,
	}
	m.connections[endpoint] = managedConn

	m.log.Info("创建新连接成功",
		"endpoint", endpoint,
		"state", conn.GetState().String(),
	)

	return conn, nil
}

// ReconnectConnection 手动重连指定服务
// 返回重连是否成功
func (m *GRPCConnectionManager) ReconnectConnection(endpoint string, dialFunc func(string) (*grpc.ClientConn, error)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	managedConn, exists := m.connections[endpoint]
	if !exists {
		// 连接不存在，创建新连接
		_, err := m.createConnection(endpoint, dialFunc)
		return err
	}

	managedConn.mu.Lock()
	defer managedConn.mu.Unlock()

	// 检查重连间隔
	if time.Since(managedConn.lastReconnect) < m.config.ReconnectInterval {
		m.log.Warn("重连间隔未到，跳过重连",
			"endpoint", endpoint,
			"interval", m.config.ReconnectInterval,
			"lastReconnect", managedConn.lastReconnect,
		)
		return nil
	}

	// 关闭旧连接
	oldConn := managedConn.conn
	_ = oldConn.Close()

	// 创建新连接
	newConn, err := dialFunc(endpoint)
	if err != nil {
		m.log.Error("手动重连失败",
			"endpoint", endpoint,
			"error", err,
		)
		return err
	}

	// 更新连接
	managedConn.conn = newConn
	managedConn.lastReconnect = time.Now()
	managedConn.reconnectCount++

	m.log.Info("手动重连成功",
		"endpoint", endpoint,
		"reconnectCount", managedConn.reconnectCount,
		"state", newConn.GetState().String(),
	)

	return nil
}

// ReconnectAllConnections 手动重连所有连接
func (m *GRPCConnectionManager) ReconnectAllConnections(dialFuncMap map[string]func(string) (*grpc.ClientConn, error)) {
	m.mu.RLock()
	endpoints := make([]string, 0, len(m.connections))
	for endpoint := range m.connections {
		endpoints = append(endpoints, endpoint)
	}
	m.mu.RUnlock()

	for _, endpoint := range endpoints {
		if dialFunc, ok := dialFuncMap[endpoint]; ok {
			if err := m.ReconnectConnection(endpoint, dialFunc); err != nil {
				m.log.Error("重连失败",
					"endpoint", endpoint,
					"error", err,
				)
			}
		}
	}
}

// CheckConnectionHealth 检查连接健康状态
// 返回是否需要重连
func (m *GRPCConnectionManager) CheckConnectionHealth(endpoint string) (bool, error) {
	m.mu.RLock()
	managedConn, exists := m.connections[endpoint]
	m.mu.RUnlock()

	if !exists {
		return true, nil // 连接不存在，需要创建
	}

	managedConn.mu.RLock()
	conn := managedConn.conn
	managedConn.mu.RUnlock()

	if conn == nil {
		return true, nil // 连接为空，需要创建
	}

	needReconnect, _ := TriggerGRPCReconnect(conn)
	return needReconnect, nil
}

// startHealthMonitor 启动健康监控（仅检测，不自动重连）
// 定期检查连接状态，记录异常日志
// 注意：实际重连由 keepalive 自动处理，或通过 ReconnectConnection 手动触发
func (m *GRPCConnectionManager) startHealthMonitor() {
	ticker := time.NewTicker(m.config.HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.RLock()
		endpoints := make([]string, 0, len(m.connections))
		for endpoint := range m.connections {
			endpoints = append(endpoints, endpoint)
		}
		m.mu.RUnlock()

		for _, endpoint := range endpoints {
			needReconnect, err := m.CheckConnectionHealth(endpoint)
			if err != nil {
				m.log.Error("健康监控检查失败",
					"endpoint", endpoint,
					"error", err,
				)
				continue
			}

			if needReconnect {
				m.log.Warn("健康监控检测到连接异常",
					"endpoint", endpoint,
					"hint", "keepalive 会自动重连，或调用 ReconnectConnection 手动重连",
				)
				// 注意：不在此处自动重连，因为：
				// 1. keepalive 已经提供了自动重连
				// 2. 重连需要 dialFunc，但监控时可能没有
				// 3. 实际重连应该在业务层调用 ReconnectConnection
			}
		}
	}
}

// Close 关闭所有连接
func (m *GRPCConnectionManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for endpoint, managedConn := range m.connections {
		managedConn.mu.Lock()
		if managedConn.conn != nil {
			if err := managedConn.conn.Close(); err != nil {
				errs = append(errs, err)
				m.log.Error("关闭连接失败",
					"endpoint", endpoint,
					"error", err,
				)
			}
		}
		managedConn.mu.Unlock()
	}

	m.connections = make(map[string]*managedConnection)

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// GetConnectionStats 获取连接统计信息
func (m *GRPCConnectionManager) GetConnectionStats() map[string]ConnectionStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]ConnectionStats)
	for endpoint, managedConn := range m.connections {
		managedConn.mu.RLock()
		conn := managedConn.conn
		reconnectCount := managedConn.reconnectCount
		lastReconnect := managedConn.lastReconnect
		managedConn.mu.RUnlock()

		state := connectivity.Idle
		if conn != nil {
			state = conn.GetState()
		}

		stats[endpoint] = ConnectionStats{
			State:          state.String(),
			ReconnectCount: reconnectCount,
			LastReconnect:  lastReconnect,
		}
	}

	return stats
}

// ConnectionStats 连接统计信息
type ConnectionStats struct {
	State          string    // 连接状态
	ReconnectCount int       // 重连次数
	LastReconnect  time.Time // 最后重连时间
}
