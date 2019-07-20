package wxwork

import "github.com/lanceryou/wxwork/http_client"

type CompanyInfo struct {
	CompanyID       string
	ApplicationName string
	ApplicationID   string
	AgentID         int64
}

type WxWorkMessageOptions struct {
	client          *http_client.HttpClient
	companyInfoList []CompanyInfo
}

type WxWorkMessageOption func(*WxWorkMessageOptions)

func WithHttpClient(client *http_client.HttpClient) WxWorkMessageOption {
	return func(o *WxWorkMessageOptions) {
		o.client = client
	}
}

func WithWxWorkMessages(msgs []CompanyInfo) WxWorkMessageOption {
	return func(o *WxWorkMessageOptions) {
		o.companyInfoList = msgs
	}
}
