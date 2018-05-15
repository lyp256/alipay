package alipay

/*查询订单*/
import (
	"errors"
)

type QueryQuest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` //订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	TradeNo    string `json:"trade_no,omitempty"`     //支付宝交易号，和商户订单号不能同时为空
}

/*trade_no,out_trade_no如果同时存在优先取trade_no 详情见https://docs.open.alipay.com/api_1/alipay.trade.query*/
func NewQuery(OutTradeNo, TradeNo string) (*QueryQuest,error) {
	if OutTradeNo == "" && TradeNo == "" {
		return nil,errors.New("OutTradeNo和TradeNo不能同时为空")
	}
	return &QueryQuest{
		OutTradeNo: OutTradeNo,
		TradeNo:    TradeNo,
	},nil

}

//返回订单查询的url
func (this *Client) QueryOrder(query *QueryQuest) (*alquest) {
	return this.newQuest(query, "alipay.trade.query")
}

//返回查询的结果
func (this *Client) QueryOrderParams(query *QueryQuest) (map[string]string, error) {
	url, err := this.QueryOrder(query).Build()
	if err != nil {
		return nil, err
	}
	return  this.httpQuest(url,"alipay_trade_query_response")
}
