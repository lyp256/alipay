package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SignType string

const (
	SignRSA  SignType = "RSA"
	SignRSA2 SignType = "RSA2"
)

type Request interface {
	SetNotifyUrl(url string)
	SetReturnUrl(url string)
	SetTimestamp(t time.Time)
	Build() (string, error)
}

// NewAlipay 创建一个alipay 应用客户端
func NewAlipay(pri, pub []byte, appid, gateway string, signType SignType) (*Client, error) {
	if appid == "" {
		return nil, errors.New("appid 不能为空")
	}
	if gateway == "" {
		return nil, errors.New("gateway 不能为空")
	}
	var c Client
	c.Gateway = gateway
	c.AppId = appid
	c.SignType = string(signType)
	if c.SignType == "RSA2" {
		c.signTypeHash = crypto.SHA256
	} else {
		c.signTypeHash = crypto.SHA1
	}
	c.PriKEY, _ = pem.Decode(pri)
	if c.PriKEY == nil {
		return nil, errors.New("私钥文件解析失败")
	}
	c.PubKEY, _ = pem.Decode(pub)
	if c.PubKEY == nil {
		return nil, errors.New("公钥钥文件解析失败")
	}

	rsapub, err := x509.ParsePKIXPublicKey(c.PubKEY.Bytes)
	if err != nil {
		return nil, err
	}
	c.pubKey = rsapub.(*rsa.PublicKey)
	c.priKey, err = x509.ParsePKCS1PrivateKey(c.PriKEY.Bytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c.basePubParam = &alPublic{
		Version: "1.0",
		Charset: "utf-8",
	}
	return &c, nil
}

// Client 应用客户端
type Client struct {
	PriKEY       *pem.Block      // 私钥pem
	PubKEY       *pem.Block      // 公钥pem
	pubKey       *rsa.PublicKey  // 公钥
	priKey       *rsa.PrivateKey // 私钥
	Gateway      string          // 支付宝网关
	SignType     string          // 签名类型
	signTypeHash crypto.Hash     // 签名使用hash类型
	AppId        string          // app_ID
	basePubParam *alPublic       // 默认通用公共参数
}

// SetNotifyUrl 设置默认异步通知url
func (pay *Client) SetNotifyUrl(url string) *Client {
	pay.basePubParam.NotifyUrl = url
	return pay
}

// SetReturnUrl 设置默认返回url
func (pay *Client) SetReturnUrl(url string) *Client {
	pay.basePubParam.ReturnUrl = url
	return pay
}

// 创建请求
func (pay *Client) newRequest(quest interface{}, method string) Request {
	if quest == nil {
		return nil
	}
	return &alRequest{
		pubParam: alPublic{
			AppId:     pay.AppId,
			Method:    method,
			ReturnUrl: pay.basePubParam.ReturnUrl,
			Charset:   pay.basePubParam.Charset,
			SignType:  pay.SignType,
			Version:   pay.basePubParam.Version,
			NotifyUrl: pay.basePubParam.NotifyUrl,
		},
		bizParam: quest,
		c:        pay,
	}
}

// 基于客户端的签名
func (pay *Client) Sign(src []byte) ([]byte, error) {
	var h = pay.signTypeHash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, pay.priKey, pay.signTypeHash, hashed)
}

// 基于客户端的验证
func (pay *Client) Verify(src, sign []byte) error {
	var h = pay.signTypeHash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(pay.pubKey, pay.signTypeHash, hashed, sign)
}

// aliResponse验证签名
func (pay *Client) ValidAliResponse(body []byte, responseName string) (map[string]string, error) {
	// 使用正则表达式寻找内容
	reg, err := regexp.Compile(`\{"` + responseName + `":(.+),"sign":"([a-zA-Z0-9/+=]+)"\}`)
	if err != nil {
		return nil, err
	}
	subs := reg.FindSubmatch(body)
	signB := make([]byte, base64.StdEncoding.DecodedLen(len(subs[2])))
	i, err := base64.StdEncoding.Decode(signB, subs[2])
	signB = signB[:i]
	err = pay.Verify(subs[1], signB)
	if err != nil {
		return nil, err
	}
	// 解析参数
	params := make(map[string]string)
	err = json.Unmarshal(subs[1], &params)
	if err != nil {
		return nil, err
	}
	return params, nil
}

// 发送同步请求,获取结果
func (pay *Client) httpDo(url, respName string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return pay.ValidAliResponse(body, respName)
}

// 组合参数成字符串
func paramsToStrings(s *[]string, i interface{}) {
	// 运用reflect 把结构体转换成参数字符窜 具体规则见蚂蚁金服签名方法 https://docs.open.alipay.com/291/106118
	tT := reflect.TypeOf(i)
	tv := reflect.ValueOf(i)
	for num, i := tT.NumField(), 0; i < num; i++ {
		var name, value string
		tags := strings.Split(tT.Field(i).Tag.Get("json"), ",")
		name = tags[0]
		if name == "" {
			continue
		}
		v := tv.Field(i)
		switch v.Type().String() {
		case "string":
			s := v.String()
			if s == "" {
				continue
			}
			value = s
		case "float64":
			f := v.Float()
			if f == 0 {
				continue
			}
			value = strconv.FormatFloat(f, 'f', 2, 64)
		}

		*s = append(*s, name+"="+value)
	}
	sort.Strings(*s) // 排序
}

// 签名方法
func Sign(src []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

// 验证签名
func Verify(src, sign []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sign)
}
