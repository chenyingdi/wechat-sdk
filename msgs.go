package WeChatSDK

import (
	"encoding/xml"
	"strconv"
)

type ReqArgs struct {
	AppId     string // appId
	AppSecret string // app密钥
	MchId     string // 商户ID
	Fee       int    // 价格
	MchName   string // 商户名
	GoodType  string // 商品类目
	OrderSn   string // 订单流水号
	NotifyUrl string // 回调信息通知地址
	PayKey    string // 支付密钥
}

type GetPrepayIdRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid"`
	MchId          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIp string   `xml:"spbill_create_ip"`
	NotifyUrl      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
}

// 初始化请求
// 参数：
// mchName: 商家名称
// goodType: 商品类目
// nonceStr
func (g *GetPrepayIdRequest) Init(args *ReqArgs) error {
	var err error

	g.AppId = args.AppId
	g.MchId = args.MchId
	g.NonceStr = GeneNonceStr(32)
	g.TotalFee = args.Fee
	g.Body = args.MchName + "-" + args.GoodType
	g.TradeType = "JSAPI"
	g.OutTradeNo = args.OrderSn

	g.SpbillCreateIp, err = GetIp()
	if err != nil {
		return err
	}

	g.NotifyUrl = args.NotifyUrl

	g.Sign = GeneSign(map[string]string{
		"appid":            g.AppId,
		"mch_id":           g.MchId,
		"nonce_str":        g.NonceStr,
		"body":             g.Body,
		"out_trade_no":     g.OutTradeNo,
		"total_fee":        strconv.Itoa(g.TotalFee),
		"spbill_create_ip": g.SpbillCreateIp,
		"notify_url":       g.NotifyUrl,
		"trade_type":       g.TradeType,
	}, args.PayKey)

	return nil
}

type RefundRequest struct {
	XMLName     xml.Name `xml:"xml"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	OutTradeNo  string   `xml:"out_trade_no"`
	OutRefundNo string   `xml:"out_refund_no"`
	TotalFee    int      `xml:"total_fee"`
}

func (r *RefundRequest) Init(args ReqArgs) {
	r.AppId = args.AppId
	r.MchId = args.MchId
	r.NonceStr = GeneNonceStr(32)
	r.OutTradeNo = args.OrderSn
	r.OutRefundNo = args.OrderSn
	r.TotalFee = args.Fee

	r.Sign = GeneSign(map[string]string{
		"appid":         r.AppId,
		"mch_id":        r.MchId,
		"nonce_str":     r.NonceStr,
		"out_trade_no":  r.OutTradeNo,
		"out_refund_no": r.OutRefundNo,
		"total_fee":     strconv.Itoa(r.TotalFee),
	}, args.PayKey)
}

type CloseOrderRequest struct {
	XMLName    xml.Name `xml:"xml"`
	AppId      string   `xml:"appid"`
	MchId      string   `xml:"mch_id"`
	NonceStr   string   `xml:"nonce_str"`
	Sign       string   `xml:"sign"`
	OutTradeNo string   `xml:"out_trade_no"`
}

func (c *CloseOrderRequest) Init(args ReqArgs)  {
	c.AppId = args.AppId
	c.MchId = args.MchId
	c.NonceStr = GeneNonceStr(32)
	c.OutTradeNo = args.OrderSn

	c.Sign = GeneSign(map[string]string{
		"appid": c.AppId,
		"mch_id": c.MchId,
		"nonce_str": c.NonceStr,
		"out_trade_no": c.OutTradeNo,
	}, args.PayKey)
}
