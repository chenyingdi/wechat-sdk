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
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func EncodingAESKey2AESKey(encodingKey string) []byte {
	data, _ := base64.StdEncoding.DecodeString(encodingKey + "=")
	return data
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Url          string
	PicUrl       string
	MediaId      string
	ThumbMediaId string
	Content      string
	MsgId        int
	Location_X   string
	Location_Y   string
	Label        string
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   string
	MsgType      CDATAText
	Content      CDATAText
}

type EncryptRequestBody struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
}

type EncryptResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATAText
	MsgSignature CDATAText
	TimeStamp    string
	Nonce        CDATAText
}

type EncryptResponseBody1 struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      string
	MsgSignature string
	TimeStamp    string
	Nonce        string
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

func MakeSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func MakeMsgSignature(token, timestamp, nonce, msg_encrypt string) string {
	sl := []string{token, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func ValidateUrl(token, timestamp, nonce, signatureIn string) bool {
	signatureGen := MakeSignature(token, timestamp, nonce)
	if signatureGen != signatureIn {
		return false
	}
	return true
}

func ValidateMsg(token, timestamp, nonce, msgEncrypt, msgSignatureIn string) bool {
	msgSignatureGen := MakeMsgSignature(token, timestamp, nonce, msgEncrypt)
	if msgSignatureGen != msgSignatureIn {
		return false
	}
	return true
}

func Value2CDATA(v string) CDATAText {
	//return CDATAText{[]byte("<![CDATA[" + v + "]]>")}
	return CDATAText{"<![CDATA[" + v + "]]>"}
}

func MakeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = Value2CDATA(fromUserName)
	textResponseBody.ToUserName = Value2CDATA(toUserName)
	textResponseBody.MsgType = Value2CDATA("text")
	textResponseBody.Content = Value2CDATA(content)
	textResponseBody.CreateTime = strconv.Itoa(int(time.Duration(time.Now().Unix())))
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}
func MakeEncryptResponseBody(appId, token, fromUserName, toUserName, content, nonce, timestamp string, key []byte) ([]byte, error) {
	encryptBody := &EncryptResponseBody{}

	encryptXmlData, _ := MakeEncryptXmlData(appId, fromUserName, toUserName, timestamp, content, key)
	encryptBody.Encrypt = Value2CDATA(encryptXmlData)
	encryptBody.MsgSignature = Value2CDATA(MakeMsgSignature(token, timestamp, nonce, encryptXmlData))
	encryptBody.TimeStamp = timestamp
	encryptBody.Nonce = Value2CDATA(nonce)

	return xml.MarshalIndent(encryptBody, " ", "  ")
}

func MakeEncryptXmlData(appId, fromUserName, toUserName, timestamp, content string, key []byte) (string, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = Value2CDATA(fromUserName)
	textResponseBody.ToUserName = Value2CDATA(toUserName)
	textResponseBody.MsgType = Value2CDATA("text")
	textResponseBody.Content = Value2CDATA(content)
	textResponseBody.CreateTime = timestamp
	body, err := xml.MarshalIndent(textResponseBody, " ", "  ")
	if err != nil {
		return "", errors.New("xml marshal error")
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		return "", err
	}
	bodyLength := buf.Bytes()

	randomBytes := []byte("abcdefghijklmnop")

	plainData := bytes.Join([][]byte{randomBytes, bodyLength, body, []byte(appId)}, nil)
	cipherData, err := AesEncrypt(plainData, key)
	if err != nil {
		return "", errors.New("AesEncrypt error")
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// PadLength calculates padding length, from github.com/vgorin/cryptogo
func PadLength(slice_length, blocksize int) (padlen int) {
	padlen = blocksize - slice_length%blocksize
	if padlen == 0 {
		padlen = blocksize
	}
	return padlen
}

//from github.com/vgorin/cryptogo
func PKCS7Pad(message []byte, blocksize int) (padded []byte) {
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		panic("block size is too small (minimum is 2 bytes)")
	}
	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		// calculate padding length
		padlen := PadLength(len(message), blocksize)

		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)

		// apply padding
		padded = append(message, padding...)
		return padded
	}
	// block size bigger or equal 256 is not currently supported
	panic("unsupported block size")
}

func AesEncrypt(plainData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey)
	if len(plainData)%k != 0 {
		plainData = PKCS7Pad(plainData, k)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)

	return cipherData, nil
}

func AesDecrypt(cipherData []byte, aesKey []byte) ([]byte, error) {
	k := len(aesKey) //PKCS#7
	if len(cipherData)%k != 0 {
		return nil, errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
}

func ValidateAppId(id []byte, appId string) bool {
	if string(id) == appId {
		return true
	}
	return false
}

func ParseEncryptTextRequestBody(appId string, plainText []byte) (*TextRequestBody, error) {

	// Read length
	buf := bytes.NewBuffer(plainText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)

	// appID validation
	appIDstart := 20 + length
	id := plainText[appIDstart : int(appIDstart)+len(appId)]
	if !ValidateAppId(id, appId) {
		return nil, errors.New("Appid is invalid")
	}

	textRequestBody := &TextRequestBody{}
	xml.Unmarshal(plainText[20:20+length], textRequestBody)
	return textRequestBody, nil
}

func ParseEncryptResponse(token string, responseEncryptTextBody, key []byte) {
	textResponseBody := &EncryptResponseBody1{}
	xml.Unmarshal(responseEncryptTextBody, textResponseBody)

	if !ValidateMsg(token, textResponseBody.TimeStamp, textResponseBody.Nonce, textResponseBody.Encrypt, textResponseBody.MsgSignature) {
		log.Println("msg signature is invalid")
		return
	}

	cipherData, err := base64.StdEncoding.DecodeString(textResponseBody.Encrypt)
	if err != nil {
		log.Println(err, "Wechat Message Service: Decode base64 error")
		return
	}

	plainText, err := AesDecrypt(cipherData, key)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(plainText))
}

func DecryptWechatAppletUser(encryptedData string, session_key string, iv string) ([]byte, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedData)
	key, _ := base64.StdEncoding.DecodeString(session_key)
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	iv_b, _ := base64.StdEncoding.DecodeString(iv)
	blockModel := cipher.NewCBCDecrypter(block, iv_b)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
