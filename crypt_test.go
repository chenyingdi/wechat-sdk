package WeChatSDK

import (
	"encoding/xml"
	"strings"
	"testing"
)

/*
WDWVSCZABGMBRBQZ  <xml>
          <ToUserName><![CDATA[ROLL_DEE_ROCK]]></ToUserName>
          <FromUserName><![CDATA[o8Auv4qR8I8izMQC5vzsNQBSefTc]]></FromUserName>
          <CreateTime>1586160069</CreateTime>
          <MsgType><![CDATA[text]]></MsgType>
          <Content><![CDATA[success]]></Content>
          <MsgId>1</MsgId>
        </xml>��;��z�%���	�Ds�9,��Ѳ��
 */

/*
  <xml><ToUserName><![CDATA[gh_10f6c3c3ac5a]]></ToUserName>
        <FromUserName><![CDATA[oyORnuP8q7ou2gfYjqLzSIWZf0rs]]></FromUserName>
        <CreateTime>1410349438</CreateTime>
        <MsgType><![CDATA[text]]></MsgType>
        <Content><![CDATA[abcdteT]]></Content>
        <MsgId>6057404712141979648</MsgId>
        </xml>wx2c2769f8efd9abc2
 */
func TestDecrypt(t *testing.T) {
	var (
		err            error
		key            []byte
		result         []byte
		sEncryptBase64 = "Bpxiue7Ro9jQFTd+wbcRzq33CkMtFwzz77KX1PuPsPyxAJTK0vsTX9P3kKh8s8DAuiH4xRNCqd2GQNq3LB1kmqFHBz7L7AX4OOHHcCloOK4e8kZdt6uEVNwu084cX9xP7uwwzjvD/9PxUREdy5bQSIHo4d1gDxW6P/SiStcN/zK3nFTE8HvPlX2E0KmY3SZGq4junfjqqlsCa0N3x7YJKXrW5WPwg4GQoq37JYeiqAnKIEDHzlBa41s1JQ7K3lCoKSMHTqqEGpk09jx1C5cq4tqW/0gtm5AHLuwWAvWrYGwcOXVNOddArBvDvsgqTP3wxzMJaNYuE5oDfiLIGtkQiO19bYHo3AKEbmfYnVZShVtDYU2VpSGeayqOgDVE7BoxkRoNgPar6FwQDEi35mb2UJkMeHBgpslXkl4EcteIAes="

		sEncodingAesKey = "mgTTO666F3reJYYZiNsIoQnxRHPzOQ8RLMXo0bIC79V"
	)

	key = EncodingAESKey2AESKey(sEncodingAesKey)

	result, err = AesDecrypt([]byte(sEncryptBase64), key)
	if err != nil{
		t.Log(err)
		return
	}

	t.Log(string(result))

	a := strings.Split(string(result), "<xml>")

	t.Log("length: ", len(a[0]))

	var data map[string]interface{}

	err = xml.Unmarshal(result, &data)
	if err != nil{
		t.Log(err)
		return
	}

	t.Log(data)

}

