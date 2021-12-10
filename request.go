package alipay

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type alRequest struct {
	pubParam alPublic    // 公共参数
	bizParam interface{} // 业务参数
	c        *Client     // 客户端
}

// SetNotifyUrl 设置异步通知地址
func (r *alRequest) SetNotifyUrl(url string) {
	r.pubParam.NotifyUrl = url
}

// SetReturnUrl 设置返回地址
func (r *alRequest) SetReturnUrl(url string) {
	r.pubParam.ReturnUrl = url
}

// SetTimestamp 设置Timestamp 此字段为请求发送时间
// 如非必要请不要调用此方法.此方法未调用
// 时Timestamp会在在Build()方法调用
// 时自动生成,调用此方法后Build()方法
// 会保留此方法设置的Timestamp
func (r *alRequest) SetTimestamp(t time.Time) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	r.pubParam.Timestamp = t.In(location).Format("2006-01-02 15:04:05")
}

// Build 创建请求
func (r *alRequest) Build() (string, error) {
	if r.bizParam == nil {
		return "", errors.New("请求业务参数不能为空")
	}
	if r.pubParam.Timestamp == "" {
		r.SetTimestamp(time.Now())
	}

	bc, err := json.Marshal(r.bizParam) // json Wap请求参数
	if err != nil {
		return "", nil
	}

	r.pubParam.BizContent = string(bc)

	params := make([]string, 0, 11)

	paramsToStrings(&params, r.pubParam)

	src := strings.Join(params, "&") // 拼接字符串

	signData, err := r.c.Sign([]byte(src))
	if err != nil {
		return "", errors.New("签名失败")
	}
	r.pubParam.Sign = base64.StdEncoding.EncodeToString(signData)

	// 参数 url encode 并返回组装完成的url返回结果
	return r.c.Gateway + "?" + r.pubParam.build().Encode(), err
}

// 其他参数目前均未固定值  暂未提供修改方法

// 请求公共参数
type alPublic struct {
	// 支付宝分配给开发者的应用ID
	AppId string `json:"app_id"`
	// 接口名称
	Method string `json:"method"`
	// 仅支持JSON
	Format string `json:"format,omitempty"`
	// 付款成功返回地址 HTTP/HTTPS开头字符串
	ReturnUrl string `json:"return_url,omitempty"`
	// 请求使用的编码格式，如utf-8,gbk,gb2312等
	Charset string `json:"charset"`
	// 商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	SignType string `json:"sign_type"`
	// 商户请求参数的签名串，详见签名
	Sign string `json:"sign"`
	// 发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
	Timestamp string `json:"timestamp"`
	// 调用的接口版本，固定为：1.0
	Version string `json:"version"`
	// 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	NotifyUrl string `json:"notify_url,omitempty"`
	// 第三方应用授权
	AppAuthToken string `json:"app_auth_token,omitempty"`
	// 业务请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
	BizContent string `json:"biz_content"`
}

// 创建为 url.values
func (p *alPublic) build() url.Values {
	var params url.Values
	params = make(map[string][]string)
	T := reflect.TypeOf(*p)
	V := reflect.ValueOf(*p)
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
