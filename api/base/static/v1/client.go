package static

import "github.com/zhixingdx/public/api/rpcclient"

func GetRPC(endpoint string) BaseStaticManageClient {
	return NewBaseStaticManageClient(rpcclient.Connect(endpoint))
}

func GetPlusRPC(endpoint string) BaseStaticManagePlusClient {
	return NewBaseStaticManagePlusClient(rpcclient.Connect(endpoint))
}

func GetPortalRPC(endpoint string) BaseStaticPortalClient {
	return NewBaseStaticPortalClient(rpcclient.Connect(endpoint))
}

func GetPortalPlusRPC(endpoint string) BaseStaticPortalPlusClient {
	return NewBaseStaticPortalPlusClient(rpcclient.Connect(endpoint))
}
