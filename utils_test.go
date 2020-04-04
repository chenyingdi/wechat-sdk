package WeChatSDK

import "testing"

func TestParseMap(t *testing.T) {
	t.Log(ParseMap(map[string]string{
		"sk":       "sdfsfs",
		"appid":    "1231231",
		"zsidfsdf": "sdfsdfs",
	}))
}

func TestGeneSign(t *testing.T) {
	t.Log(GeneSign(map[string]string{
		"appid":       "wxd930ea5d5a258f4f",
		"mch_id":      "10000100",
		"device_info": "1000",
		"body":        "test",
		"nonce_str":   "ibuaiVcKdpRxkhJA",
	}, "192006250b4c09247ec02edce69f6a2d"))
}
