package WeChatSDK

import "encoding/xml"

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

type CloseOrderRequest struct {
	XMLName xml.Name `xml:"xml"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	NonceStr    string   `xml:"nonce_str"`
	Sign        string   `xml:"sign"`
	OutTradeNo  string   `xml:"out_trade_no"`
}
