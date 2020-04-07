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

// 获取接口调用凭证
func (c *Client) GetAccessToken() (string, error) {
	var (
		err  error
		resp *http.Response
		body []byte
		data map[string]interface{}
		url  = "https://api.weixin.qq.com/cgi-bin/token" +
			"?grant_type=client_credential" +
			"&appid=" + c.AppId +
			"&secret=" + c.AppSecret
	)

	resp, err = http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	if data["access_token"] == nil {
		return "", errors.New("there is not access_token in response body... ")
	}

	c.accessToken = data["access_token"].(string)
	return data["access_token"].(string), nil
}

// 获取日访问留存
func (c *Client) GetDailyRetain(date string) (map[string]interface{}, error) {
	var (
		err     error
		reqBody []byte
		resBody []byte
		r       *bytes.Reader
		res     *http.Response
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappiddailyretaininfo" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{
			date,
			date,
		}
	)

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取月访问留存
func (c *Client) GetMonthlyRetain(year int, month int) (map[string]interface{}, error) {
	var (
		err     error
		begin   string
		end     string
		reqBody []byte
		r       *bytes.Reader
		res     *http.Response
		resBody []byte
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappidmonthlyretaininfo" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{}
	)

	begin, end = GetBeginAndEndByMonth(year, month)
	req.BeginDate = begin
	req.EndDate = end

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取周访问留存
func (c *Client) GetWeeklyRetain() (map[string]interface{}, error) {
	var (
		err     error
		begin   string
		end     string
		reqBody []byte
		r       *bytes.Reader
		res     *http.Response
		resBody []byte
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappidweeklyretaininfo" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{}
	)

	begin, end = GetBeginAndEndByWeek()

	req.BeginDate = begin
	req.EndDate = end

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取日统计
func (c *Client) GetDailySummary(date string) (map[string]interface{}, error) {
	var (
		err     error
		reqBody []byte
		resBody []byte
		r       *bytes.Reader
		res     *http.Response
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappiddailysummarytrend" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{
			date,
			date,
		}
	)

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取日趋势
func (c *Client) GetDailyVisitTrend(date string) (map[string]interface{}, error) {
	var (
		err     error
		reqBody []byte
		resBody []byte
		r       *bytes.Reader
		res     *http.Response
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappiddailyvisittrend" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{
			date,
			date,
		}
	)

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取周趋势
func (c *Client) GetWeeklyVisitTrend() (map[string]interface{}, error) {
	var (
		err     error
		begin   string
		end     string
		reqBody []byte
		r       *bytes.Reader
		res     *http.Response
		resBody []byte
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappidweeklyvisittrend" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{}
	)

	begin, end = GetBeginAndEndByWeek()

	req.BeginDate = begin
	req.EndDate = end

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取月趋势
func (c *Client) GetMonthlyVisitTrend(year, month int) (map[string]interface{}, error) {
	var (
		err     error
		begin   string
		end     string
		reqBody []byte
		r       *bytes.Reader
		res     *http.Response
		resBody []byte
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappidmonthlyvisittrend" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{}
	)

	begin, end = GetBeginAndEndByMonth(year, month)
	req.BeginDate = begin
	req.EndDate = end

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取用户画像
func (c *Client) GetUserPortrait(date string) (map[string]interface{}, error) {
	var (
		err     error
		reqBody []byte
		resBody []byte
		r       *bytes.Reader
		res     *http.Response
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappiduserportrait" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{
			date,
			date,
		}
	)

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil

}

// 获取用户分布
func (c *Client) GetVisitDistribution(date string) (map[string]interface{}, error) {
	var (
		err     error
		reqBody []byte
		resBody []byte
		r       *bytes.Reader
		res     *http.Response
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappidvisitdistribution" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{
			date,
			date,
		}
	)

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// 获取页面数据
func (c *Client) GetVisitPage(date string) (map[string]interface{}, error) {
	var (
		err     error
		reqBody []byte
		resBody []byte
		r       *bytes.Reader
		res     *http.Response
		resData map[string]interface{}
		url     = "https://api.weixin.qq.com/datacube/getweanalysisappidvisitpage" +
			"?access_token=" + c.accessToken
		req = struct {
			BeginDate string `json:"begin_date"`
			EndDate   string `json:"end_date"`
		}{
			date,
			date,
		}
	)

	reqBody, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(reqBody)

	res, err = http.Post(url, "application/json", r)
	if err != nil {
		return nil, err
	}

	resBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return nil, err
	}

	return resData, nil
}
