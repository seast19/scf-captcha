package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/liujiawm/gocaptcha"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

const salt = "@!asdc3453"

// DefineEvent 入参参数
type DefineEvent struct {
	// test event define
	Body string `json:"body"`
}

// Input 用户输入
type Input struct {
	UserCode       string `json:"usercode"`       //用户输入验证码
	UserCipherText string `json:"userciphertext"` //用户输入密文

	Action string `json:"action"` //目标动作，new 或者 check
}

// Data 数据结构
type Data struct {
	Img        string `json:"img"`        //图片base64
	CipherText string `json:"ciphertext"` // sha256(code+salt+time)

	CheckStatus string `json:"checkstatus"` //验证状态 ok 或者 fail
}

// New 生成图片验证码
func New() *Data {
	// 生成图片验证码
	data, _ := gocaptcha.New(&gocaptcha.Options{
		CharPreset: "0123456789", // 数字作基数
		Curve:      2,            // 两条弧线
		Length:     4,            // 长度为4的验证码
		Width:      80,           // 图片宽
		Height:     33,           // 图片高
	})

	//生成hash     sha256(code + salt+time)
	timeText := strconv.Itoa(int(time.Now().Unix()))
	cipherText := GetSHA256HashCode([]byte(data.Text + salt + timeText))

	res := &Data{
		Img:        data.EncodeB64string(),
		CipherText: cipherText + "#" + timeText,
	}

	// fmt.Println(res)
	fmt.Println(res.CipherText, data.Text)

	return res

}

// Check 检验用户输入验证码
func Check(data *Input) *Data {
	//获取参数
	res := &Data{
		CheckStatus: "-1",
	}

	cipherText := strings.Split(data.UserCipherText, "#")[0]
	timeText := strings.Split(data.UserCipherText, "#")[1]

	//校验时间戳
	now := time.Now().Unix()
	timeNum, err := strconv.Atoi(timeText)
	if err != nil {
		fmt.Println(err)
		return res
	}

	//过期时间10分钟
	if now-int64(timeNum) > 600 {
		res.CheckStatus = "-2"
		fmt.Println(res)
		return res
	}

	//校验验证码

	if GetSHA256HashCode([]byte(data.UserCode+salt+timeText)) == cipherText {
		res.CheckStatus = "1"
	}

	fmt.Println(res)

	return res
}

// GetSHA256HashCode SHA256生成哈希值
func GetSHA256HashCode(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

// scf入口函数
func scf(event DefineEvent) (interface{}, error) {
	params := Input{}

	err := json.Unmarshal([]byte(event.Body), &params)
	if err != nil {
		fmt.Println("json反序列失败->", err)
		return "", err
	}

	if params.Action == "new" {
		return New(), nil
	} else if params.Action == "check" {
		return Check(&params), nil
	}

	return "", errors.New("无action参数")
}

func main() {
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(scf)
	// New()

}
