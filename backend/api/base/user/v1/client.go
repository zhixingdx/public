// Package user client
package user

import "github.com/zhixingdx/oipivot-api/rpcclient"

// GetRPC 获取 RPC 客户端
func GetRPC(endpoint string) BaseUserManageClient {
	return NewBaseUserManageClient(rpcclient.Connect(endpoint))
}

// GetPlusRPC 获取 PlusRPC 客户端
func GetPlusRPC(endpoint string) BaseUserManagePlusClient {
	return NewBaseUserManagePlusClient(rpcclient.Connect(endpoint))
}

// GetPortalRPC 获取 PortalRPC 客户端
func GetPortalRPC(endpoint string) BaseUserPortalClient {
	return NewBaseUserPortalClient(rpcclient.Connect(endpoint))
}

// GetPortalPlusRPC 获取 PortalPlusRPC 客户端
func GetPortalPlusRPC(endpoint string) BaseUserPortalPlusClient {
	return NewBaseUserPortalPlusClient(rpcclient.Connect(endpoint))
}
