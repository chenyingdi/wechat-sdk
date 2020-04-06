package WeChatSDK

// 客户端类
type Client struct {
	AppId          string
	AppSecret      string
	MchId          string
	EncodingAesKey string
	Token          string
}

func NewClient(appId, appSecret, mchId string) *Client {
	return &Client{
		AppId:     appId,
		AppSecret: appSecret,
		MchId:     mchId,
	}
}
