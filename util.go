package weapp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

// tokenAPI 获取带 token 的 API 地址
func tokenAPI(api, token string) (string, error) {
	queries := requestQueries{
		"access_token": token,
	}

	return encodeURL(api, queries)
}

// encodeURL add and encode parameters.
func encodeURL(api string, params requestQueries) (string, error) {
	url, err := url.Parse(api)
	if err != nil {
		return "", err
	}

	query := url.Query()

	for k, v := range params {
		query.Set(k, v)
	}

	url.RawQuery = query.Encode()

	return url.String(), nil
}

// randomString random string generator
//
// ln length of return string
func randomString(ln int) string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, ln)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}

	return string(b)
}

// postJSON perform a HTTP/POST request with json body
func postJSON(url string, params interface{}, response interface{}) error {
	resp, err := postJSONWithBody(url, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}

func getJSON(url string, response interface{}) error {

	resp, err := httpClient().Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}

// postJSONWithBody return with http body.
func postJSONWithBody(url string, params interface{}) (*http.Response, error) {
	b := &bytes.Buffer{}
	if params != nil {
		enc := json.NewEncoder(b)
		enc.SetEscapeHTML(false)
		err := enc.Encode(params)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("=====",b.String())

	return httpClient().Post(url, "application/json; charset=utf-8", b)
}

func postFormByFile(url, field, filename string, response interface{}) error {
	// Add your media file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return postForm(url, field, filename, file, response)
}

func postForm(url, field, filename string, reader io.Reader, response interface{}) error {
	// Prepare a form that you will submit to that URL.
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, err := w.CreateFormFile(field, filename)
	if err != nil {
		return err
	}

	if _, err = io.Copy(fw, reader); err != nil {
		return err
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := httpClient()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// convert bool to int
func bool2int(ok bool) uint8 {

	if ok {
		return 1
	}

	return 0
}

//--------- aes256加解密

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func PKCS7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	padtext = append(text, padtext...)
	return padtext
}

func PKCS7UnPadding(text []byte) []byte {
	length := len(text)
	unpadding := int(text[length-1])
	return text[:(length - unpadding)]
}

// 补0
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 获取随机字符byte
func GetRandomByte(length int) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	strLen := len(str)
	strBytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, strBytes[r.Intn(strLen)])
	}
	return result
}

func getIv(length int) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	strBytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, strBytes[i])
	}
	return result
}

//aes256cbc加密
func Encrypt(key string, data string) (string, error) {
	// 获取key
	keyByte := []byte(key)

	// 选取加密算法
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	// pkcs7补位
	text := PKCS7Padding([]byte(data), block.BlockSize())

	//iv
	//iv := GetRandomByte(block.BlockSize())
	iv := getIv(block.BlockSize())

	//加密
	cbc := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(text))
	cbc.CryptBlocks(ciphertext, text)

	//base64加密
	ciphertextStr := EncodeBase64(ciphertext)
	return ciphertextStr, nil
}

func Decrypt(key string, b64 string) (string, error) {
	// 获取key
	keyByte := []byte(key)

	//base64解密
	ciphertext, err := DecodeBase64(b64)
	if err != nil {
		return "", err
	}

	// 选取加密算法
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	//iv
	iv := getIv(block.BlockSize())

	//解密
	text := make([]byte, len(ciphertext))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(text, ciphertext)

	// 反解pkcs7补位
	text = PKCS7UnPadding(text)
	return string(text), nil
}