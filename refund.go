package alipay

import (
	"errors"
)

/*退款接口*/

// Goods 商品形象
type Goods struct {
	// 商品的编号
	GoodsId string `json:"goods_id"`
	// 支付宝定义的统一商品编号	20010001
	AlipayGoodsId string `json:"alipay_goods_id,omitempty"`
	// 256	商品名称
	GoodsName string `json:"goods_name"`
	// 10	商品数量	1
	Quantity int `json:"quantity"`
	// 9 商品单价，单位为元	2000
	Price float64 `json:"price"`
	// 24	商品类目	34543238
	GoodsCategory string `json:"goods_category,omitempty"`
	// 1000 商品描述信息	特价手机
	Body string `json:"body,omitempty"`
	// 400	商品的展示地址	http://www.alipay.com/xxx.jpg
	ShowUrl string `json:"show_url,omitempty"`
}

// Refund 退款
type Refund struct {
	// 订单支付时传入的商户订单号,不能和 trade_no同时为空。
	OutTradeNo string `json:"out_trade_no,omitempty"`
	// 支付宝交易号，和商户订单号不能同时为空
	TradeNo string `json:"trade_no,omitempty"`
	// 需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundAmount float64 `json:"refund_amount"`
	// 订单退款币种信息，非外币交易，不能传入退款币种信息
	RefundCurrency string `json:"refund_currency,omitempty"`
	// 退款的原因说明
	RefundReason string `json:"refund_reason,omitempty"`
	// 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OutRequestNo string `json:"out_request_no,omitempty"`
	// 商户的操作员编号
	OperatorId string `json:"operator_id,omitempty"`
	// 商户的门店编号
	StoreId string `json:"store_id,omitempty"`
	// 商户的终端编号
	TerminalId string `json:"terminal_id,omitempty"`
	// 退款包含的商品列表信息，Json格式。
	GoodsDetail []*Goods `json:"goods_detail,omitempty"`
}

// NewRefund 创建一个退款
func NewRefund(outNo, NO string, RefundAmount float64) (*Refund, error) {
	if outNo == "" && NO == "" {
		return nil, errors.New("OutNo和No不能同时为空")
	}

	return &Refund{
		OutTradeNo:   outNo,
		TradeNo:      NO,
		RefundAmount: RefundAmount,
	}, nil
}

// SetRefundCurrency 设置退款币种
func (r *Refund) SetRefundCurrency(Currency string) *Refund {
	r.RefundCurrency = Currency
	return r
}

// SetRefundReason 设置退款说明
func (r *Refund) SetRefundReason(Reason string) *Refund {
	r.RefundReason = Reason
	return r
}

// SetRefundOutRequestNo 设置退款编号
func (r *Refund) SetRefundOutRequestNo(No string) *Refund {
	r.OutRequestNo = No
	return r
}

// SetOperatorId 设置操作员id
func (r *Refund) SetOperatorId(id string) *Refund {
	r.OperatorId = id
	return r
}

// SetStoreId 设置店铺id
func (r *Refund) SetStoreId(id string) *Refund {
	r.StoreId = id
	return r
}

// SetTerminalId 设置终端id
func (r *Refund) SetTerminalId(id string) *Refund {
	r.TerminalId = id
	return r
}

// SetGoodsDetail 设置商品列表
func (r *Refund) SetGoodsDetail(gs []*Goods) *Refund {
	r.GoodsDetail = gs
	return r
}

// AddGoodsDetail 添加商品
func (r *Refund) AddGoodsDetail(g *Goods) *Refund {
	if r.GoodsDetail == nil {
		r.GoodsDetail = make([]*Goods, 0, 8)
	}
	r.GoodsDetail = append(r.GoodsDetail, g)
	return r
}
func (pay *Client) Refund(re *Refund) Request {
	return pay.newRequest(re, "alipay.trade.refund")
}
func (pay *Client) RefundParams(re *Refund) (map[string]string, error) {
	url, err := pay.Refund(re).Build()
	if err != nil {
		return nil, err
	}
	return pay.httpDo(url, "alipay_trade_refund_response")
}

// QueryRefund 查询退款
type QueryRefund struct {
	TradeNo      string `json:"trade_no,omitempty"`     // 支付宝交易号，和商户订单号不能同时为空
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no
	OutRequestNo string `json:"out_request_no"`         // 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
}

// NewQueryRefund 交易退款查询
func NewQueryRefund(OutTradeNo, TradeNo, OutRequestNo string) (*QueryRefund, error) {
	if TradeNo == "" && OutTradeNo == "" {
		return nil, errors.New("TradeNo和OutTradeNo不能同时为空")
	}
	if OutRequestNo == "" {
		return nil, errors.New("OutRequestNo不能为空")

	}
	return &QueryRefund{
		TradeNo:      TradeNo,
		OutTradeNo:   OutTradeNo,
		OutRequestNo: OutRequestNo,
	}, nil
}

// QueryRefund 交易退款查询
func (pay *Client) QueryRefund(query *QueryRefund) Request {
	return pay.newRequest(query, "alipay.trade.fastpay.refund.query")
}
