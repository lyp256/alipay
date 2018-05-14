package alipay
//
func (this *Client) CloseOrser(closeQuest *CloseQuest) (*alquest) {
	return this.newQuest(closeQuest, "alipay.trade.close")
}
/*关闭订单请求*/
type CloseQuest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` //订单支付时传入的商户订单号,和支付宝交易号不能同时为空。
	TradeNo    string `json:"trade_no,omitempty"`     //支付宝交易号，和商户订单号不能同时为空
	OperatorId string `json:"operator_id,omitempty"`  //自定义操作员id
}

func NewCloseQuest(OutTradeNo, TradeNo, OperatorId string) (quest *CloseQuest) {
	if OutTradeNo == "" && TradeNo == "" {
		return nil
	}
	return &CloseQuest{
		OutTradeNo: OutTradeNo,
		TradeNo:    TradeNo,
		OperatorId: OperatorId,
	}

}
