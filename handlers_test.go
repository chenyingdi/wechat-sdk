package WeChatSDK

import (
	"testing"
	"time"
)

func TestClient_GetAccessToken(t *testing.T) {
	c := NewClient(
		"wxe9284cefae4934d0",
		"0deff095cc4da391b51887cb1906e8cd",
		"",
	)

	accessToken, err := c.GetAccessToken()
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(accessToken)
}

func TestClient_GetDailyRetain(t *testing.T) {
	c := NewClient(
		"wxe9284cefae4934d0",
		"0deff095cc4da391b51887cb1906e8cd",
		"",
	)

	_, err := c.GetAccessToken()
	if err != nil {
		t.Log(err)
		return
	}

	date := time.Now().AddDate(0, 0, -1).Format("20060102")

	res, err := c.GetDailyRetain(date)
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(res)
}

func TestClient_GetMonthlyRetain(t *testing.T) {
	c := NewClient(
		"wxe9284cefae4934d0",
		"0deff095cc4da391b51887cb1906e8cd",
		"",
	)

	_, err := c.GetAccessToken()
	if err != nil {
		t.Log(err)
		return
	}

	res, err := c.GetMonthlyRetain(2020, 2)
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(res)
}

func TestClient_GetWeeklyRetain(t *testing.T) {
	c := NewClient(
		"wxe9284cefae4934d0",
		"0deff095cc4da391b51887cb1906e8cd",
		"",
	)

	_, err := c.GetAccessToken()
	if err != nil {
		t.Log(err)
		return
	}

	res, err := c.GetWeeklyRetain()
	if err != nil{
		t.Log(err)
		return
	}

	t.Log(res)
}

func TestClient_GetUserPortrait(t *testing.T) {
	c := NewClient(
		"wxe9284cefae4934d0",
		"0deff095cc4da391b51887cb1906e8cd",
		"",
	)

	_, err := c.GetAccessToken()
	if err != nil {
		t.Log(err)
		return
	}

	res, err := c.GetUserPortrait("20200406")
	if err != nil{
		t.Log(err)
		return
	}

	t.Log(res)
}

func TestClient_GetDailySummary(t *testing.T) {
	c := NewClient(
		"wxe9284cefae4934d0",
		"0deff095cc4da391b51887cb1906e8cd",
		"",
	)

	_, err := c.GetAccessToken()
	if err != nil {
		t.Log(err)
		return
	}

	res, err := c.GetDailySummary("20200407")
	if err != nil{
		t.Log(err)
		return
	}

	t.Log(res)
}
