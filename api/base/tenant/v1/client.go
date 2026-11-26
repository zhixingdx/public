// Package tenant client
package tenant

import "github.com/zhixingdx/oipivot-api/rpcclient"

// GetRPC 获取 RPC 客户端
func GetRPC(endpoint string) BaseTenantManageClient {
	return NewBaseTenantManageClient(rpcclient.Connect(endpoint))
}

// GetPlusRPC 获取 PlusRPC 客户端
func GetPlusRPC(endpoint string) BaseTenantManagePlusClient {
	return NewBaseTenantManagePlusClient(rpcclient.Connect(endpoint))
}

// GetPortalRPC 获取 PortalRPC 客户端
func GetPortalRPC(endpoint string) BaseTenantPortalClient {
	return NewBaseTenantPortalClient(rpcclient.Connect(endpoint))
}

// GetPortalPlusRPC 获取 PortalPlusRPC 客户端
func GetPortalPlusRPC(endpoint string) BaseTenantPortalPlusClient {
	return NewBaseTenantPortalPlusClient(rpcclient.Connect(endpoint))
}
