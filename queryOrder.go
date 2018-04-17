package alipay

import (
	"net/http"
	"io/ioutil"
)

type QueryQuest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` //订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	TradeNo    string `json:"trade_no,omitempty"`     //支付宝交易号，和商户订单号不能同时为空
}

/*trade_no,out_trade_no如果同时存在优先取trade_no 详情见https://docs.open.alipay.com/api_1/alipay.trade.query*/
func NewQueryQuest(OutTradeNo, TradeNo string) (*QueryQuest) {
	if OutTradeNo == "" && TradeNo == "" {
		return nil
	}
	return &QueryQuest{
		OutTradeNo: OutTradeNo,
		TradeNo:    TradeNo,
	}

}

//返回订单查询的url
func (this *Client) QueryOrder(query *QueryQuest) (string, error) {
	return this.newQuest(query, "alipay.trade.query", "")
}

//返回查询的结果
func (this *Client) QueryOrderParams(query *QueryQuest) (map[string]string, error) {
	url, err := this.QueryOrder(query)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return this.ValidAliResponse(body, "alipay_trade_query_response")
}
