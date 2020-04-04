package WeChatSDK

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

// 获取用户openid
func (c *Client) GetOpenid(code string) (string, error) {
	url := "https://api.weixin.qq.com/sns/jscode2session" +
		"?appid=" + c.AppId + "&secret=" + c.AppSecret + "&js_code=" + code +
		"&grant_type=authorization_code"

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if data["openid"] == nil {
		return "", ResponseError
	}

	return data["openid"].(string), nil
}

// 获取预支付ID
func (c *Client) GetPrepayId(r GetPrepayIdRequest) (interface{}, error) {
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"

	reqBody, err := xml.Marshal(&r)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(url, "application/xml", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	err = xml.Unmarshal(resBody, &data)
	if err != nil {
		return nil, err
	}

	if data["return_code"] == nil || data["return_msg"] == nil {
		return nil, ResponseError
	}

	if data["return_code"].(string) != "SUCCESS" {
		return nil, errors.New("get prepay id error: " + data["return_msg"].(string))
	}

	if data["result_code"] == nil {
		return nil, ResponseError
	}

	if data["result_code"].(string) != "SUCCESS" {
		return nil, errors.New("get prepay id error: " +
			"[" + data["err_code"].(string) + "] " + data["err_code_des"].(string))
	}

	if data["prepay_id"] == nil {
		return nil, ResponseError
	}

	return data["prepay_id"], nil
}

// 退款
func (c *Client) Refund(r RefundRequest) error {
	url := "https://api.mch.weixin.qq.com/secapi/pay/refund"

	req, err := xml.Marshal(&r)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/xml", bytes.NewReader(req))
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})

	err = xml.Unmarshal(resBody, &data)
	if err != nil {
		return err
	}

	if data["return_code"] == nil || data["return_msg"] == nil {
		return ResponseError
	}

	if data["return_code"].(string) != "SUCCESS" {
		return errors.New("refund error: " + data["return_msg"].(string))
	}

	if data["result_code"] == nil {
		return ResponseError
	}

	if data["result_code"] != "SUCCESS" {
		return errors.New("refund error: " +
			"[" + data["err_code"].(string) + "] " + data["err_code_des"].(string))
	}

	return nil
}

// 关闭订单
func (c *Client) CloseOrder(r CloseOrderRequest) error {
	url := "https://api.mch.weixin.qq.com/pay/closeorder"

	req, err := xml.Marshal(&r)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "application/xml", bytes.NewReader(req))
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})

	err = xml.Unmarshal(resBody, &data)
	if err != nil {
		return err
	}

	if data["return_code"] == nil || data["return_msg"] == nil {
		return ResponseError
	}

	if data["return_code"].(string) != "SUCCESS" {
		return errors.New("refund error: " + data["return_msg"].(string))
	}

	if data["result_code"] == nil {
		return ResponseError
	}

	if data["result_code"] != "SUCCESS" {
		return errors.New("refund error: " +
			"[" + data["err_code"].(string) + "] " + data["err_code_des"].(string))
	}

	return nil
}

