package alipay

import (
	"net/url"
	"sort"
	"strings"
	"errors"
	"crypto"
	"encoding/base64"
	"fmt"
	"net/http"
)

/*通知相关*/
const success = "success"

func (this *Client) ValidateNotify(from url.Values) (error) {

	sign := from.Get("sign")
	signType := from.Get("sign_type")
	if sign == "" || signType == "" {
		return errors.New("没有sign或sign_type")
	}
	//
	params := make([]string, len(from)-2)[0:0]
	//获取所有除sign和sign_type外的所有参数
	for k := range map[string][]string(from) {
		if k != "sign" && k != "sign_type" {
			params = append(params, k+"="+from[k][0])
		}
	}
	sort.Strings(params)
	src := strings.Join(params, "&")
	signbyte, err := base64.StdEncoding.DecodeString(sign)
	fmt.Println(signbyte)
	fmt.Println(sign)
	fmt.Println(signType)
	fmt.Println(src)
	if err != nil {
		return err
	}
	return Verify([]byte(src), signbyte, this.PubKEY, crypto.SHA256)
}

func NotIfySuccess(w http.ResponseWriter) (err error) {
	_, err = w.Write([]byte(success))
	return err

}
