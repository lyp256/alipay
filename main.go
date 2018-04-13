package main

import (
	"crypto"
	"encoding/pem"
	"errors"
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"reflect"
	"strings"
	"strconv"
	"sort"
	"encoding/json"
	"encoding/base64"
	"time"
	"net/url"
)

type SignType string

const (
	ST_RSA  SignType = "RSA"
	ST_RSA2 SignType = "RSA2"
)

type Client struct {
	PriKEY   *pem.Block //私钥
	PubKEY   *pem.Block //公钥
	Gateway  string     //支付宝网关
	SignType string     //签名类型
	AppId    string     //app_ID
}

/*创建一个alipay 应用客户端*/
func NewAlipay(pri, pub []byte, appid, gateway string, sigType SignType) (*Client, error) {
	if appid == "" {
		return nil, errors.New("appid不能为空")
	}
	if gateway == "" {
		return nil, errors.New("gateway不能为空")
	}
	var c Client
	c.Gateway = gateway
	c.AppId = appid
	c.SignType = string(sigType)
	c.PriKEY, _ = pem.Decode(pri)
	if c.PriKEY == nil {
		return nil, errors.New("私钥文件解析失败")
	}
	c.PubKEY, _ = pem.Decode(pub)
	if c.PubKEY == nil {
		return nil, errors.New("公钥钥文件解析失败")
	}

	return &c, nil
}

func (this *Client) newQuest(quest interface{}, method, returnUrl string) (string, error) {
	if quest == nil {
		return "", errors.New("请求参数不能为空")
	}
	var err error
	var PubParam alPublic
	PubParam.Charset = "utf-8" //请求字符集
	PubParam.Version = "1.0"   //接口版本
	PubParam.Method = method   //接口方法
	PubParam.AppId = this.AppId
	PubParam.ReturnUrl = returnUrl
	PubParam.SignType = this.SignType

	PubParam.Timestamp = time.Now().Format("2006-01-02 15:04:05") //接口时间
	bc, err := json.Marshal(quest)                                //json Wap请求参数
	if err != nil {
		return "", nil
	}

	PubParam.BizContent = string(bc)

	params := make([]string, 30)[0:0]

	paramsToStrings(&params, PubParam)

	src := strings.Join(params, "&") //拼接字符串
	if (this.SignType == "RSA2") {
		PubParam.Sign, err = Sign([]byte(src), this.PriKEY, crypto.SHA256)
	} else {
		PubParam.Sign, err = Sign([]byte(src), this.PriKEY, crypto.SHA1)
	}
	if err != nil {
		return "", errors.New("签名失败")
	}
	//对签名结果进行urlencode
	var v url.Values
	v = make(map[string][]string)
	v.Set("sign", PubParam.Sign)
	//组装参数
	src += "&" + v.Encode()
	//参数urlencode并返回组装完成的url返回结果
	return this.Gateway + "?" + PubParam.build().Encode(), err
}

//请求公共参数
type alPublic struct {
	AppId      string `json:"app_id"`               //支付宝分配给开发者的应用ID
	Method     string `json:"method"`               //接口名称
	Format     string `json:"format,omitempty"`     //仅支持JSON
	ReturnUrl  string `json:"return_url,omitempty"` //HTTP/HTTPS开头字符串
	Charset    string `json:"charset"`              //请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType   string `json:"sign_type"`            //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign       string `json:"sign"`                 //商户请求参数的签名串，详见签名
	Timestamp  string `json:"timestamp"`            //发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
	Version    string `json:"version"`              //调用的接口版本，固定为：1.0
	NotifyUrl  string `json:"notify_url,omitempty"` //支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BizContent string `json:"biz_content"`          //业务请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
}

/*创建为 url.values*/
func (this *alPublic) build() (url.Values) {
	var params url.Values
	params = make(map[string][]string)
	T := reflect.TypeOf(*this)
	V := reflect.ValueOf(*this)
	for i, n := 0, T.NumField(); i < n; i++ {
		var name, value string
		typeName := T.Field(i)
		tags := strings.Split(typeName.Tag.Get("json"), ",")
		name = tags[0]
		if name == "" {
			continue
		}

		switch typeName.Type.String() {
		case "string":
			s := V.Field(i).String()
			if s == "" {
				continue
			}
			value = s
		case "float64":
			f := V.Field(i).Float()
			value = strconv.FormatFloat(f, 'f', 2, 64)
		}

		params.Set(name, value)

	}
	return params
}

/*组合参数成字符串*/
func paramsToStrings(s *[]string, i interface{}) {
	//运用reflect 把结构体转换成参数字符窜 具体规则见蚂蚁金服签名方法 https://docs.open.alipay.com/291/106118
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
	sort.Strings(*s) //排序
}

//签名方法
func Sign(src []byte, key *pem.Block, hash crypto.Hash) (string, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	var err error
	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(key.Bytes)
	if err != nil {
		return "", err
	}
	b, err := rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
	en := base64.StdEncoding
	return en.EncodeToString(b), nil
}
