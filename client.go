package WeChatSDK

// 客户端类
type Client struct {
	AppId     string
	AppSecret string
	MchId     string
}

func NewClient(appId, appSecret string) *Client {
	return &Client{
		AppId:     appId,
		AppSecret: appSecret,
	}
}
