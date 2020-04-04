package WeChatSDK

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net"
	"sort"
	"strings"
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
	var (
		stringA        string
		stringSignTemp string
		sign           string
		m              = md5.New()
	)
	// 1. 字典序排序
	stringA = ParseMap(args)

	// 2. 与key拼接得到stringSignTemp
	stringSignTemp = stringA + "&key=" + key

	m.Write([]byte(stringSignTemp))

	sign = strings.ToUpper(hex.EncodeToString(m.Sum(nil)))

	return sign
}

func ParseMap(args map[string]string) string {
	var (
		keys    []string
		stringA string
	)
	for k, _ := range args {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		if stringA == "" {
			stringA += v + "=" + args[v]
		} else {
			stringA += "&" + v + "=" + args[v]
		}
	}

	return stringA
}

// 获取本机IP
func GetIp() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addr {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}
