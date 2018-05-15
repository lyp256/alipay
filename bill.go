package alipay

import (
	"errors"
	"time"
	"net/http"
	"io/ioutil"
)

/*账单*/

type BillDown struct {
	BillType string `json:"bill_type"` //账单类型 商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单
	BillDate string `json:"bill_date"` //账单时间：日账单格式为yyyy-MM-dd，月账单格式为yyyy-MM。
}

func NewBillDown(btype, date string) (*BillDown, error) {
	if btype != "trade" && btype != "signcustomer" {
		return nil, errors.New("错误的账单类型")
	}
	_, e1 := time.Parse("2006-01", date)
	_, e2 := time.Parse("2006-01-02", date)
	if e1 != nil && e2 != nil {
		return nil, errors.New("错误的账单日期")
	}

	return &BillDown{BillType: btype, BillDate: date,}, nil
}
func (this *Client)BillDownUrl(bill *BillDown) *alquest {
 return this.newQuest(bill,"alipay.data.dataservice.bill.downloadurl.query")
}
func (this *Client) BillDownParams(bill *BillDown) (map[string]string, error) {
	url, err := this.BillDownUrl(bill).Build()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return this.ValidAliResponse(body, "alipay_data_dataservice_bill_downloadurl_query_response")
}
