package WeChatSDK

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

// 生成随机字符串
func GeneNonceStr(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}

	return string(bytes)
}

// 签名
func GeneSign(args map[string]string, key string) string {
	stringA := ""

	if args != nil {
		for k, v := range args {
			if stringA == "" {
				stringA += k + "=" + v
			}else {
				stringA += "&" + k + "=" + v
			}
		}
	}

	stringSignTemp := stringA + "&key=" + key

	m := md5.New()

	m.Write([]byte(stringSignTemp))

	sign := hex.EncodeToString(m.Sum(nil))

	return sign
}


