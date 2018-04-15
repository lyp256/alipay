package alipay
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

//订单查询
func (this *Client) QueryOrder(closeQuest *QueryQuest) (string, error) {
	return this.newQuest(closeQuest, "alipay.trade.query", "")
}
