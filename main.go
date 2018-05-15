package alipay

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
	"regexp"
	"fmt"
	"net/http"
	"io/ioutil"
)

type SignType string

const (
	ST_RSA  SignType = "RSA"
	ST_RSA2 SignType = "RSA2"
)

// 应用客户端
type Client struct {
	PriKEY       *pem.Block      //私钥pem
	PubKEY       *pem.Block      //公钥pem
	pubKey       *rsa.PublicKey  //公钥
	priKey       *rsa.PrivateKey //私钥
	Gateway      string          //支付宝网关
	SignType     string          //签名类型
	signTypeHash crypto.Hash     //签名使用hash类型
	AppId        string          //app_ID
	basePubParam *alPublic       //默认通用公共参数
}

/*创建一个alipay 应用客户端*/
func NewAlipay(pri, pub []byte, appid, gateway string, signType SignType) (*Client, error) {
	var err error
	if appid == "" {
		return nil, errors.New("appid不能为空")
	}
	if gateway == "" {
		return nil, errors.New("gateway不能为空")
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

//设置默认异步通知url
func (this *Client) SetNotifyUrl(url string) (*Client) {
	this.basePubParam.NotifyUrl = url
	return this
}

//设置默认返回url
func (this *Client) SetReturnUrl(url string) (*Client) {
	this.basePubParam.ReturnUrl = url
	return this
}

type alquest struct {
	pubParam alPublic    //公共参数
	bizParam interface{} //业务参数
	c        *Client     //客户端
}

//创建请求 内部通用方法
func (this *Client) newQuest(quest interface{}, method string) (*alquest) {
	if quest == nil {
		return nil
	}
	var q alquest
	q.c = this
	q.bizParam = quest
	q.pubParam.Charset = this.basePubParam.Charset     //请求字符集
	q.pubParam.Version = this.basePubParam.Version     //接口版本
	q.pubParam.Method = method                         //接口方法
	q.pubParam.AppId = this.AppId                      //appid
	q.pubParam.NotifyUrl = this.basePubParam.NotifyUrl //异步通知地址
	q.pubParam.ReturnUrl = this.basePubParam.ReturnUrl //返回地址
	q.pubParam.SignType = this.SignType
	return &q
}

//设置异步通知地址
func (this *alquest) SetNotifyUrl(url string) (*alquest) {
	this.pubParam.NotifyUrl = url
	return this
}

//设置返回地址
func (this *alquest) SetReturnUrl(url string) (*alquest) {
	this.pubParam.ReturnUrl = url
	return this
}

//设置Timestamp 此字段为请求发送时间
//如非必要请不要调用此方法.此方法未调用
// 时Timestamp会在在Build()方法调用
// 时自动生成,调用此方法后Build()方法
// 会保留此方法设置的Timestamp

func (this *alquest) SetTimestamp(t time.Time) (*alquest) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	this.pubParam.Timestamp = t.In(location).Format("2006-01-02 15:04:05")
	return this
}

/*其他字段为固定值 暂时未提供修改方法*/

//创建请求 内部通用方法

func (this *alquest) Build() (string, error) {
	if this.bizParam == nil {
		return "", errors.New("请求业务参数不能为空")
	}
	if this.bizParam == nil {
		return "", errors.New("请求业务参数不能为空")
	}
	if this.pubParam.Timestamp=="" {
		this.SetTimestamp(time.Now())
	}

	bc, err := json.Marshal(this.bizParam)                                      //json Wap请求参数
	if err != nil {
		return "", nil
	}

	this.pubParam.BizContent = string(bc)

	params := make([]string, 11)[0:0]

	paramsToStrings(&params, this.pubParam)

	src := strings.Join(params, "&") //拼接字符串

	this.pubParam.Sign, err = this.c.Sign([]byte(src))
	if err != nil {
		return "", errors.New("签名失败")
	}
	/*//对签名结果进行urlencode
	var v url.Values
	v = make(map[string][]string)
	v.Set("sign", this.pubParam.Sign)
	//组装参数
	src += "&" + v.Encode()*/
	//参数urlencode并返回组装完成的url返回结果
	return this.c.Gateway + "?" + this.pubParam.build().Encode(), err
}

//其他参数目前均未固定值  暂未提供修改方法

//请求公共参数
type alPublic struct {
	AppId      string `json:"app_id"`               //支付宝分配给开发者的应用ID
	Method     string `json:"method"`               //接口名称
	Format     string `json:"format,omitempty"`     //仅支持JSON
	ReturnUrl  string `json:"return_url,omitempty"` //付款成功返回地址 HTTP/HTTPS开头字符串
	Charset    string `json:"charset"`              //请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType   string `json:"sign_type"`            //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign       string `json:"sign"`                 //商户请求参数的签名串，详见签名
	Timestamp  string `json:"timestamp"`            //发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
	Version    string `json:"version"`              //调用的接口版本，固定为：1.0
	NotifyUrl  string `json:"notify_url,omitempty"` //支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	AppAuthToken string `json:"app_auth_token,omitempty"` //第三方应用授权
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

//基于客户端的签名
func (this *Client) Sign(src []byte) (string, error) {
	var h = this.signTypeHash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	b, err := rsa.SignPKCS1v15(rand.Reader, this.priKey, this.signTypeHash, hashed)
	if err != nil {
		return "", err
	}
	en := base64.StdEncoding
	return en.EncodeToString(b), nil
}

//验证签名 供外部使用
func Verify(src, sign []byte, key *pem.Block, hash crypto.Hash) (error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	var err error
	var pub interface{}
	pub, err = x509.ParsePKIXPublicKey(key.Bytes)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), hash, hashed, sign)
}

//基于客户端的验证
func (this *Client) Verify(src, sign []byte) (error) {

	var h = this.signTypeHash.New()
	h.Write(src)
	var hashed = h.Sum(nil)
	var err error
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(this.pubKey, this.signTypeHash, hashed, sign)
}

//aliResponse验证签名
func (this *Client) ValidAliResponse(body []byte, responseName string) (map[string]string, error) {
	//使用正则表达式寻找内容
	reg, err := regexp.Compile(`\{"` + responseName + `":(.+),"sign":"([a-zA-Z0-9/+=]+)"\}`)
	if err != nil {
		return nil, err
	}
	subs := reg.FindSubmatch(body)
	signB := make([]byte, base64.StdEncoding.DecodedLen(len(subs[2])))
	i, err := base64.StdEncoding.Decode(signB, subs[2])
	signB = signB[:i]
	err = this.Verify(subs[1], signB)
	if err != nil {
		return nil, err
	}
	//解析参数
	params := make(map[string]string)
	err = json.Unmarshal(subs[1], &params)
	if err != nil {
		return nil, err
	}
	return params, nil
}

//发送同步请求,获取结果
func (this *Client) httpQuest(url,respName string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return this.ValidAliResponse(body, respName)
}