package alipay

import (
	"crypto"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const success = "success"

func (pay *Client) ValidateNotify(from url.Values) error {
	sign := from.Get("sign")
	signType := from.Get("sign_type")
	if sign == "" || signType == "" {
		return errors.New("没有sign或sign_type")
	}
	//
	params := make([]string, 0, len(from)-2)
	// 获取所有除sign和sign_type外的所有参数
	for k := range map[string][]string(from) {
		if k != "sign" && k != "sign_type" {
			params = append(params, k+"="+from[k][0])
		}
	}
	sort.Strings(params)
	src := strings.Join(params, "&")
	signBuf, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return Verify([]byte(src), signBuf, pay.pubKey, crypto.SHA256)
}

func NotIfySuccess(w http.ResponseWriter) (err error) {
	_, err = w.Write([]byte(success))
	return err
}
