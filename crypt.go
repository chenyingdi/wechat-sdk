package WeChatSDK

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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
}

type TextRes struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      CDATA    `xml:"MsgType"`
	Content      CDATA    `xml:"Content"`
	MsgId        string   `xml:"MsgId"`
}

type CDATA struct {
	Text string `xml:",cdata"`
}

//@brief:将EncodingAesKey转换为AesKey
func EncAesKey2AesKey(encAesKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encAesKey + "=")
}

//@brief: 填充明文
//补码
func PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}

//@brief: 去除填充数据
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

//@brief:AES解密
func AesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

//@brief: AES加密
//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func AESEncrypt(rawData, key []byte) (string, error) {
	data, err := AesCBCEncrypt(rawData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func AESDecrypt(rawData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return nil, err
	}
	dnData, err := AesCBCDecrypt(data, key)
	if err != nil {
		return nil, err
	}
	return dnData, nil
}

func Value2CDATA(value string) CDATA {
	return CDATA{Text: value}
}

func MakeTextRes(fromUserName, toUserName, timestamp, content, msgId string) ([]byte, error) {
	var (
		textRes = new(TextRes)
	)

	textRes.FromUserName = Value2CDATA(fromUserName)
	textRes.ToUserName = Value2CDATA(toUserName)
	textRes.Content = Value2CDATA(content)
	textRes.MsgType = Value2CDATA("text")
	textRes.MsgId = msgId

	if timestamp == "" {
		textRes.CreateTime = strconv.Itoa(int(time.Duration(time.Now().Unix())))
	} else {
		textRes.CreateTime = timestamp
	}

	return xml.MarshalIndent(textRes, "", "  ")
}

func MakeEncryptXmlData(appId, fromUserName, toUserName, timestamp, content, msgId string, key []byte) (string, error) {
	var (
		err        error
		body       []byte
		bodyLength []byte
		random     []byte
		plainData  []byte
		buf        = new(bytes.Buffer)
	)
	// random(16B) + msg_len(4B) + msg + appId
	body, err = MakeTextRes(fromUserName, toUserName, timestamp, content, msgId)
	if err != nil {
		return "", err
	}

	err = binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		return "", err
	}
	bodyLength = buf.Bytes()

	random = []byte(GeneNonceStr(16))
	plainData = bytes.Join([][]byte{random, bodyLength, body, []byte(appId)}, nil)

	return AESEncrypt(plainData, key)
}

func MakeEncryptRes(appId, token, fromUserName, toUserName, timestamp, content, msgId string, key []byte) ([]byte, error) {
	var (
		err            error
		encryptXmlData string
		encryptRes     = new(EncryptRes)
	)

	encryptXmlData, err = MakeEncryptXmlData(appId, fromUserName, toUserName, timestamp, content, msgId, key)
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
