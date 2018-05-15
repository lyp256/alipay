package alipay

import (
	"net/http"
	"io/ioutil"
	"errors"
)

//
func (this *Client) CloseOrder(closeQuest *CloseQuest) (*alquest) {
	return this.newQuest(closeQuest, "alipay.trade.close")
}
/*关闭订单请求*/
type CloseQuest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` //订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	TradeNo    string `json:"trade_no,omitempty"`     //支付宝交易号，和商户订单号不能同时为空
	OperatorId string `json:"operator_id,omitempty"`  //自定义操作员id
}

func NewClose(OutTradeNo, TradeNo, OperatorId string) (*CloseQuest,error) {
	if OutTradeNo == "" && TradeNo == "" {
		return nil,errors.New("OutTradeNo和TradeNo不能同时为空")
	}
	return &CloseQuest{
		OutTradeNo: OutTradeNo,
		TradeNo:    TradeNo,
		OperatorId: OperatorId,
	},nil

}
func (this *Client) CloseOrser(re *CloseQuest) (*alquest) {
	return this.newQuest(re, "alipay.trade.close")
}
func (this *Client) CloseOrderParams(close *CloseQuest) (map[string]string, error) {
	url, err := this.CloseOrder(close).Build()
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
	return this.ValidAliResponse(body, "alipay_trade_close_response")
}