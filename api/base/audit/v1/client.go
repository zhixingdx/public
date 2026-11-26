package audit

import "github.com/zhixingdx/public/api/rpcclient"

func GetRPC(endpoint string) BaseAuditManageClient {
	return NewBaseAuditManageClient(rpcclient.Connect(endpoint))
}

func GetPlusRPC(endpoint string) BaseAuditManagePlusClient {
	return NewBaseAuditManagePlusClient(rpcclient.Connect(endpoint))
}

func GetPortalRPC(endpoint string) BaseAuditPortalClient {
	return NewBaseAuditPortalClient(rpcclient.Connect(endpoint))
}

func GetPortalPlusRPC(endpoint string) BaseAuditPortalPlusClient {
	return NewBaseAuditPortalPlusClient(rpcclient.Connect(endpoint))
}
