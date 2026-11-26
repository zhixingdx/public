package iam

import "github.com/zhixingdx/oipivot-api/rpcclient"

func GetRPC(endpoint string) BaseIamManageClient {
	return NewBaseIamManageClient(rpcclient.Connect(endpoint))
}

func GetPlusRPC(endpoint string) BaseIamManagePlusClient {
	return NewBaseIamManagePlusClient(rpcclient.Connect(endpoint))
}

func GetPortalRPC(endpoint string) BaseIamPortalClient {
	return NewBaseIamPortalClient(rpcclient.Connect(endpoint))
}

func GetPortalPlusRPC(endpoint string) BaseIamPortalPlusClient {
	return NewBaseIamPortalPlusClient(rpcclient.Connect(endpoint))
}
