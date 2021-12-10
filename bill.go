package alipay

import (
	"fmt"
	"time"
)

// BillDown 账单
type BillDown struct {
	// 账单类型 商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型
	// trade: 指商户基于支付宝交易收单的业务账单；
	// signcustomer: 是指基于商户支付宝余额收入及支出等资金变动的帐务账单
	BillType string `json:"bill_type"`
	// 账单时间：日账单格式为yyyy-MM-dd，月账单格式为yyyy-MM。
	BillDate string `json:"bill_date"`
}

// NewBillDown 创建账单
func NewBillDown(bType, date string) (*BillDown, error) {
	if bType != "trade" && bType != "signcustomer" {
		return nil, fmt.Errorf("invalid bill type %s", bType)
	}
	_, e1 := time.Parse("2006-01", date)
	_, e2 := time.Parse("2006-01-02", date)
	if e1 != nil && e2 != nil {
		return nil, fmt.Errorf("invalid data:%s", date)
	}

	return &BillDown{BillType: bType, BillDate: date}, nil
}

func (pay *Client) BillDownUrl(bill *BillDown) Request {
	return pay.newRequest(bill, "alipay.data.dataservice.bill.downloadurl.query")
}

func (pay *Client) BillDownParams(bill *BillDown) (map[string]string, error) {
	url, err := pay.BillDownUrl(bill).Build()
	if err != nil {
		return nil, err
	}
	return pay.httpDo(url, "alipay_data_dataservice_bill_downloadurl_query_response")
}
