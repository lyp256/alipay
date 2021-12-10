package alipay

import (
	"errors"
)

type ExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id,omitempty"` // 否	64	系统商编号，该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID	2088511833207846
	HbFqNum              string `json:"hb_fq_num,omitempty"`               // 否	5	花呗分期数（目前仅支持3、6、12）注：使用该参数需要仔细阅读“花呗分期接入文档”	3
	HbFqSellerPercent    string `json:"hb_fq_seller_percent,omitempty"`    //	否	3	卖家承担收费比例，商家承担手续费传入100，用户承担手续费传入0，仅支持传入100、0两种，其他比例暂不支持注：使用该参数需要仔细阅读“花呗分期接入文档”	100
}

// PagePay 电脑网页支付
type PagePay struct {
	// 64商户订单号，64个字符以内、可包含字母、数字、下划线；需保证在商户端不重复	20150320010101001
	OutTradeNo string `json:"out_trade_no"`
	// 64销售产品码，与支付宝签约的产品码名称。 注：目前仅支持FAST_INSTANT_TRADE_PAY	FAST_INSTANT_TRADE_PAY
	ProductCode string `json:"product_code"`
	// 11 订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	TotalAmount float64 `json:"total_amount"`
	// 256 订单标题
	Subject string `json:"subject"`
	//	128	订单描述
	Body string `json:"body,omitempty"`
	// 订单包含的商品列表信息
	// Json格式： {&quot;show_url&quot;:&quot;https://或http://打头的商品的展示地址&quot;}
	// 在支付时，可点击商品名称跳转到该地址	{&quot;show_url&quot;:&quot;https://www.alipay.com&quot;}
	GoodsDetail string `json:"goods_detail,omitempty"`
	// 512	公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。
	// 支付宝只会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	// merchantBizType%3d3C%26merchantBizNo%3d2016010101111
	PassbackParams string `json:"passback_params,omitempty"`
	// 业务扩展参数，详见业务扩展参数说明	{&quot;sys_service_provider_id&quot;:&quot;2088511833207846&quot;}
	ExtendParams *ExtendParams `json:"extend_params,omitempty"`
	// 商品主类型：0&mdash;虚拟类商品，1&mdash;实物类商品（默认）	注：虚拟类商品不支持使用花呗渠道	0
	GoodsType string `json:"goods_type,omitempty"`
	// 该笔订单允许的最晚付款时间，逾期将关闭交易。
	// 取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。
	// 该参数数值不接受小数点， 如 1.5h，可转换为 90m。	该参数在请求到支付宝时开始计时。	90m
	TimeoutExpress string `json:"timeout_express,omitempty"`
	// 128	可用渠道，用户只能在指定渠道范围内支付
	// 当有多个渠道时用&ldquo;,&rdquo;分隔
	// 注：与disable_pay_channels互斥	pcredit,moneyFund,debitCardExpress
	EnablePayChannels string `json:"enable_pay_channels,omitempty"`
	// 禁用渠道，用户不可用指定渠道支付
	// 当有多个渠道时用&ldquo;,&rdquo;分隔
	// 注：与enable_pay_channels互斥	pcredit,moneyFund,debitCardExpress
	DisablePayChannels string `json:"disable_pay_channels,omitempty"`
	// 取用户授权信息，可实现如免登功能。
	AuthToken string `json:"auth_token,omitempty"`
	// 否 2 PC扫码支付的方式，支持前置模式和跳转模式。
	// 前置模式是将二维码前置到商户的订单确认页的模式。需要商户在自己的页面中以iframe方式请求支付宝页面。
	// 具体分为以下几种:
	// 0: 订单码-简约前置模式，对应iframe宽度不能小于600px，高度不能小于300px；
	// 1: 订单码-前置模式，对应iframe宽度不能小于300px，高度不能小于600px；
	// 2: 订单码-跳转模式
	// 3: 订单码-迷你前置模式，对应iframe宽度不能小于75px，高度不能小于75px；
	// 4: 订单码-可定义宽度的嵌入式二维码，商户可根据需要设定二维码的大小。
	// 跳转模式下，用户的扫码界面由支付宝生成的，不在商户的域名下。
	QrPayMode string `json:"qr_pay_mode,omitempty"`
	// 商户自定义二维码宽度	注：qr_pay_mode=4时该参数生效	100
	QrcodeWidth string `json:"qrcode_width,omitempty"`
}

func NewPagePay(OutTradeNo, Subject string, TotalAmount float64) (*PagePay, error) {
	if Subject == "" {
		return nil, errors.New("subject 不能为空")
	}
	if OutTradeNo == "" {
		return nil, errors.New("outTradeNo 不能为空")
	}
	return &PagePay{
		OutTradeNo:  OutTradeNo,
		TotalAmount: TotalAmount,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		Subject:     Subject,
	}, nil

}
func (pay *PagePay) SetProductCode(ProductCode string) *PagePay {
	pay.ProductCode = ProductCode
	return pay
}
func (pay *PagePay) SetBody(Body string) *PagePay {
	pay.Body = Body
	return pay
}
func (pay *PagePay) SetGoodsDetail(GoodsDetail string) *PagePay {
	pay.GoodsDetail = GoodsDetail
	return pay
}
func (pay *PagePay) SetPassbackParams(PassbackParams string) *PagePay {
	pay.PassbackParams = PassbackParams
	return pay
}
func (pay *PagePay) SetExtendParams(Ext *ExtendParams) *PagePay {
	pay.ExtendParams = Ext
	return pay
}
func (pay *PagePay) SetGoodsType(GoodsType string) *PagePay {
	pay.GoodsType = GoodsType
	return pay
}
func (pay *PagePay) SetTimeoutExpress(TimeoutExpress string) *PagePay {
	pay.GoodsDetail = TimeoutExpress
	return pay
}
func (pay *PagePay) SetEnablePayChannels(EnablePayChannels string) *PagePay {
	pay.EnablePayChannels = EnablePayChannels
	return pay
}
func (pay *PagePay) SetDisablePayChannels(DisablePayChannels string) *PagePay {
	pay.DisablePayChannels = DisablePayChannels
	return pay
}
func (pay *PagePay) SetAuthToken(AuthToken string) *PagePay {
	pay.AuthToken = AuthToken
	return pay
}
func (pay *PagePay) SetQrPayMode(QrPayMode string) *PagePay {
	pay.QrPayMode = QrPayMode
	return pay
}
func (pay *PagePay) SetQrcodeWidth(QrcodeWidth string) *PagePay {
	pay.GoodsDetail = QrcodeWidth
	return pay
}
func (pay *Client) PagePay(page *PagePay) Request {
	return pay.newRequest(page, "alipay.trade.page.pay")
}
