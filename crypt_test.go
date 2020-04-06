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
		sEncryptBase64 = "mfBCs65c67CeJw22u4VT2TD73q5H06+ocrAIxswCaeZ/d/Lw" +
			"0msSZFHY0teqgSYiI1zR2gD2DKrB3TIrmX/liNSDrGqS8jSI/" +
			"WPeKB5VPr7Ezr7gomZAyGCwJSgT1TRFWPfONGJMxuj2nk4faTu" +
			"spAuVIFQ6SHwZuJBZC7mcJp7Cgr9cUhATQWDbOPaE7ukZBTV2Yq" +
			"yzH+UI2AK+J1S47cE79k1RX8t0hcTz/O0hlK8DGXKnvYv88qKQcI" +
			"7z4iaajqHfRVZKBNyOODabs+It+ZfM3dWTeFcPgDbGtIEnpt/EDtu" +
			"uA/zMvtkaKdHdswPnVZQ+xdwbYr3ldGvfT8HlEYEgkgKaThxTFobVl" +
			"wzu2ZkXCjicbP3xdr15Iq48ObgzPpqYuZ3IEoyggZDKClquk0u0orMck4GTF/XyE8yGzc4="

		sEncodingAesKey = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
	)

	key, err = EncAesKey2AesKey(sEncodingAesKey)
	if err != nil {
		t.Log(err)
		return
	}

	result, err = AESDecrypt(sEncryptBase64, key)
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

