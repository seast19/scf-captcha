package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	Code       int    `json:"code"`       //状态码2000成功，2901失败
	Msg        string `json:"msg"`        //提示信息
	Img        string `json:"img"`        //图片base64
	CipherText string `json:"ciphertext"` // sha256(code+salt+time)

	CheckStatus string `json:"checkstatus"` //验证状态 1 或者 -1 -2
}

// New 生成图片验证码
func New() *Data {
	// 生成图片验证码
	data, _ := gocaptcha.New(&gocaptcha.Options{
		CharPreset: "0123456789", // 数字作基数
		Curve:      2,            // 两条弧线
		Length:     4,            // 长度为4的验证码
		Width:      80,           // 图片宽
		Height:     40,           // 图片高
	})

	//生成hash     sha256(code + salt+time)
	timeText := strconv.Itoa(int(time.Now().Unix()))
	cipherText := GetSHA256HashCode([]byte(data.Text + salt + timeText))

	res := &Data{
		Code:       2000,
		Msg:        "创建验证码成功",
		Img:        data.EncodeB64string(),
		CipherText: cipherText + "#" + timeText,
	}

	fmt.Println("创建验证码成功")
	fmt.Printf("%s -> %s", data.Text, res.CipherText)

	return res
}

// Check 检验用户输入验证码
func Check(data *Input) *Data {
	//返回参数
	res := &Data{
		Code:        2901,
		Msg:         "验证码错误",
		CheckStatus: "-1",
	}

	fmt.Println("开始校验验证码")

	// 验证参数放置报错
	if len(strings.Split(data.UserCipherText, "#")) != 2 {
		fmt.Println("ctoken参数错误")
		return res
	}

	cipherText := strings.Split(data.UserCipherText, "#")[0]
	timeText := strings.Split(data.UserCipherText, "#")[1]

	//校验时间戳
	now := time.Now().Unix()
	timeNum, err := strconv.Atoi(timeText)
	if err != nil {
		fmt.Println("时间戳错误")
		return res
	}

	//过期时间10分钟
	if now-int64(timeNum) > 600 {
		fmt.Println("ctoken已过期")
		res.Msg = "验证码已过期"
		res.CheckStatus = "-2"
		return res
	}

	//校验验证码

	if GetSHA256HashCode([]byte(data.UserCode+salt+timeText)) == cipherText {
		fmt.Println("校验成功")
		res.Code = 2000
		res.Msg = "验证码校验成功"
		res.CheckStatus = "1"
	}

	fmt.Println("验证码校验失败")
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
	res := &Data{
		Code: 2901,
		Msg:  "请求参数错误",
	}

	params := Input{}

	err := json.Unmarshal([]byte(event.Body), &params)
	if err != nil {
		fmt.Println("json反序列入参失败->", err)
		return res, nil
	}

	fmt.Println("输入内容")
	fmt.Println(params)

	if params.Action == "new" {
		return New(), nil
	} else if params.Action == "check" {
		return Check(&params), nil
	}

	return res, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(scf)
	// New()

}
