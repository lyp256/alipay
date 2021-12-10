package alipay

import (
	"errors"
)

// CloseQuest 关闭订单请求
type CloseQuest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` // 订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	TradeNo    string `json:"trade_no,omitempty"`     // 支付宝交易号，和商户订单号不能同时为空
	OperatorId string `json:"operator_id,omitempty"`  // 自定义操作员id
}

// NewCloseRequest 创建 CloseRequest
func NewCloseRequest(OutTradeNo, TradeNo, OperatorId string) (*CloseQuest, error) {
	if OutTradeNo == "" && TradeNo == "" {
		return nil, errors.New("OutTradeNo和TradeNo不能同时为空")
	}
	return &CloseQuest{
		OutTradeNo: OutTradeNo,
		TradeNo:    TradeNo,
		OperatorId: OperatorId,
	}, nil

}

// CloseOrder 创建请求
func (pay *Client) CloseOrder(re *CloseQuest) Request {
	return pay.newRequest(re, "alipay.trade.close")
}

// CloseOrderParams 发送同步请求,获取结果
func (pay *Client) CloseOrderParams(close *CloseQuest) (map[string]string, error) {
	url, err := pay.CloseOrder(close).Build()
	if err != nil {
		return nil, err
	}
	return pay.httpDo(url, "alipay_trade_close_response")
}
