package alipay

import (
	"errors"
)

/*退款接口*/

//商品形象
type Goods struct {
	GoodsId       string  `json:"goods_id"`                  //商品的编号
	AlipayGoodsId string  `json:"alipay_goods_id,omitempty"` //支付宝定义的统一商品编号	20010001
	GoodsName     string  `json:"goods_name"`                //256	商品名称
	Quantity      int     `json:"quantity"`                  //10	商品数量	1
	Price         float64 `json:"price"`                     //9 商品单价，单位为元	2000
	GoodsCategory string  `json:"goods_category,omitempty"`  //24	商品类目	34543238
	Body          string  `json:"body,omitempty"`            //1000 商品描述信息	特价手机
	ShowUrl       string  `json:"show_url,omitempty"`        //400	商品的展示地址	http://www.alipay.com/xxx.jpg
}

//退款
type Refunnd struct {
	OutTradeNo     string   `json:"out_trade_no,omitempty"`    //订单支付时传入的商户订单号,不能和 trade_no同时为空。
	TradeNo        string   `json:"trade_no,omitempty"`        //支付宝交易号，和商户订单号不能同时为空
	RefundAmount   float64  `json:"refund_amount"`             //需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundCurrency string   `json:"refund_currency,omitempty"` //订单退款币种信息，非外币交易，不能传入退款币种信息
	RefundReason   string   `json:"refund_reason,omitempty"`   //退款的原因说明
	OutRequestNo   string   `json:"out_request_no,omitempty"`  //标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorId     string   `json:"operator_id,omitempty"`     //商户的操作员编号
	StoreId        string   `json:"store_id,omitempty"`        //商户的门店编号
	TerminalId     string   `json:"terminal_id,omitempty"`     //商户的终端编号
	GoodsDetail    []*Goods `json:"goods_detail,omitempty"`    //退款包含的商品列表信息，Json格式。
}

//创建一个退款
func NewRefunnd(outNo, NO string, RefundAmount float64) (*Refunnd,error) {
	if outNo==""&&NO=="" {
		return nil,errors.New("OutNo和No不能同时为空")
	}


	return &Refunnd{
		OutTradeNo:   outNo,
		TradeNo:      NO,
		RefundAmount: RefundAmount,
	},nil
}

//设置退款币种
func (this *Refunnd) SetRefundCurrency(Currency string) (*Refunnd) {
	this.RefundCurrency = Currency
	return this
}

//设置退款说明
func (this *Refunnd) SetRefundReason(Reason string) (*Refunnd) {
	this.RefundReason = Reason
	return this
}

//设置退款编号
func (this *Refunnd) SetRefundOutRequestNo(No string) (*Refunnd) {
	this.OutRequestNo = No
	return this
}

//设置操作员id
func (this *Refunnd) SetOperatorId(id string) (*Refunnd) {
	this.OperatorId = id
	return this
}

//设置店铺id
func (this *Refunnd) SetStoreId(id string) (*Refunnd) {
	this.StoreId = id
	return this
}

//设置终端id
func (this *Refunnd) SetTerminalId(id string) (*Refunnd) {
	this.TerminalId = id
	return this
}

//设置商品列表
func (this *Refunnd) SetGoodsDetail(gs []*Goods) (*Refunnd) {
	this.GoodsDetail = gs
	return this
}

//添加商品
func (this *Refunnd) AddGoodsDetail(g *Goods) (*Refunnd) {
	if this.GoodsDetail == nil {
		this.GoodsDetail = make([]*Goods, 8)[0:0]
	}
	this.GoodsDetail = append(this.GoodsDetail, g)
	return this
}
func (this *Client) Refund(re *Refunnd) (*alquest) {
	return this.newQuest(re, "alipay.trade.refund")
}
func (this *Client) RefundParams(re *Refunnd) (map[string]string, error) {
	url, err := this.Refund(re).Build()
	if err != nil {
		return nil, err
	}
	return  this.httpQuest(url,"alipay_trade_refund_response")
}
