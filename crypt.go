package WeChatSDK

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

type EncryptRes struct {
	XMLName      xml.Name `xml:"xml"`
	Nonce        string   `xml:"Nonce"`
	TimeStamp    string   `xml:"TimeStamp"`
	MsgSignature CDATA    `xml:"MsgSignature"`
	Encrypt      CDATA    `xml:"Encrypt"`
	Content      string   `xml:"Content"`
}

type TextRes struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      CDATA    `xml:"MsgType"`
	Content      CDATA    `xml:"Content"`
}

type CDATA struct {
	Text string `xml:",cdata"`
}

//@brief:将EncodingAesKey转换为AesKey
func EncAesKey2AesKey(encAesKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encAesKey + "=")
}

//@brief: 填充明文
func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

//@brief: 去除填充数据
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//@brief:AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// AES分组长度为 128 位，所以 blockSize=16 字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//@brief: AES加密
func AESEncrypt(origData, key []byte) []byte {
	//获取block块
	block, _ := aes.NewCipher(key)
	//补码
	origData = PKCS5Padding(origData, block.BlockSize())
	//加密模式，
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	//创建明文长度的数组
	crypted := make([]byte, len(origData))

	//加密明文
	blockMode.CryptBlocks(crypted, origData)

	return crypted
}

func Value2CDATA(value string) CDATA {
	return CDATA{Text: value}
}

func MakeTextRes(fromUserName, toUserName, timestamp, content string) ([]byte, error) {
	var (
		textRes = new(TextRes)
	)

	textRes.FromUserName = Value2CDATA(fromUserName)
	textRes.ToUserName = Value2CDATA(toUserName)
	textRes.Content = Value2CDATA(content)
	textRes.MsgType = Value2CDATA("text")

	if timestamp == "" {
		textRes.CreateTime = strconv.Itoa(int(time.Duration(time.Now().Unix())))
	} else {
		textRes.CreateTime = timestamp
	}

	return xml.MarshalIndent(textRes, "", "  ")
}

func MakeEncryptXmlData(appId, fromUserName, toUserName, timestamp, content string, key []byte) (string, error) {
	var (
		err         error
		body        []byte
		bodyLength  []byte
		plainData   []byte
		buf         = new(bytes.Buffer)
		randomBytes = []byte("abcdefghijklmnop")
	)

	body, err = MakeTextRes(fromUserName, toUserName, timestamp, content)
	if err != nil {
		return "", err
	}

	err = binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		return "", err
	}

	bodyLength = buf.Bytes()
	plainData = bytes.Join([][]byte{randomBytes, bodyLength, body, []byte(appId)}, nil)
	return base64.StdEncoding.EncodeToString(AESEncrypt(plainData, key)), nil
}

func MakeEncryptRes(appId, token,  fromUserName, toUserName, timestamp, content string, key []byte) ([]byte, error) {
	var (
		err            error
		encryptXmlData string
		encryptRes     = new(EncryptRes)
	)

	encryptXmlData, err = MakeEncryptXmlData(appId, fromUserName, toUserName, timestamp, content, key)
	if err != nil {
		return nil, err
	}

	encryptRes.Encrypt = Value2CDATA(encryptXmlData)
	encryptRes.TimeStamp = timestamp
	encryptRes.Nonce = GeneNonceStr(32)
	encryptRes.MsgSignature = Value2CDATA(MakeMsgSignature(token, timestamp, encryptRes.Nonce, encryptXmlData))

	return xml.MarshalIndent(encryptRes, "", "  ")
}

func MakeMsgSignature(token, timestamp, nonce, msg_encrypt string) string {
	sl := []string{token, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
