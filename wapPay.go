package alipay

import (
	"time"
	"errors"
)

//手机订单
type WapQuest struct {
	Body               string  `json:"body,omitempty"`                 //对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body。
	Subject            string  `json:"subject"`                        //	商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo         string  `json:"out_trade_no"`                   //商户网站唯一订单号
	TimeoutExpress     string  `json:"timeout_express,omitempty"`      //该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。注：若为空，则默认为15d。
	TimeExpire         string  `json:"time_expire,omitempty"`          //绝对超时时间，格式为yyyy-MM-dd HH:mm。 注：1）以支付宝系统时间为准；2）如果和timeout_express参数同时传入，以time_expire为准。
	TotalAmount        float64 `json:"total_amount"`                   //订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	AuthToken          string  `json:"auth_token,omitempty"`           //针对用户授权接口，获取用户相关数据时，用于标识用户授权关系注：若不属于支付宝业务经理提供签约服务的商户，暂不对外提供该功能，该参数使用无效。
	ProductCode        string  `json:"product_code"`                   //销售产品码，商家和支付宝签约的产品码。该产品请填写固定值：QUICK_WAP_WAY
	GoodsType          string  `json:"goods_type,omitempty"`           //商品主类型：0—虚拟类商品，1—实物类商品注：虚拟类商品不支持使用花呗渠道
	PassbackParams     string  `json:"passback_params,omitempty"`      //公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	PromoParams        string  `json:"promo_params,omitempty"`         //优惠参数注：仅与支付宝协商后可用
	ExtendParams       string  `json:"extend_params,omitempty"`        //业务扩展参数，详见下面的“业务扩展参数说明”
	EnablePayChannels  string  `json:"enable_pay_channel,omitempty"`   //可用渠道，用户只能在指定渠道范围内支付当有多个渠道时用“,”分隔注：与disable_pay_channels互斥
	DisablePayChannels string  `json:"disable_pay_channels,omitempty"` //禁用渠道，用户不可用指定渠道支付当有多个渠道时用“,”分隔注：与enable_pay_channels互斥
	StoreId            string  `json:"store_id,omitempty"`             //商户门店编号。该参数用于请求参数中以区分各门店，非必传项。
	QuitUrl            string  `json:"quit_url,omitempty"`             //添加该参数后在h5支付收银台会出现返回按钮，可用于用户付款中途退出并返回到该参数指定的商户网站地址。注：该参数对支付宝钱包标准收银台下的跳转不生效。
	ExtUserInfo        string  `json:"ext_user_info,omitempty"`        //外部指定买家，详见外部用户ExtUserInfo参数说明
}

/*创建一个wap订单请求*/
func NewWapPay(OutTradeNo,Subject  string, TotalAmount float64) (*WapQuest, error) {
	if Subject == "" {
		return nil, errors.New("Subject不能为空")
	}
	if OutTradeNo == "" {
		return nil, errors.New("OutTradeNo不能为空")
	}
	return &WapQuest{
		Subject:     Subject,
		OutTradeNo:  OutTradeNo,
		TotalAmount: TotalAmount,
	}, nil
}

/*
参数修改方法
返回值提供链式写法*/
func (this *WapQuest) SetTimeExpire(unix int64) (*WapQuest) {
	this.TimeExpire = time.Unix(unix, 0).Format("2006-01-02 15:04:05")
	return this

}
func (this *WapQuest) SetBody(body string) (*WapQuest) {
	this.Body = body
	return this

}

//手机支付url
func (this *Client) WapPay(wapQuest *WapQuest) (*alquest) {
	return this.newQuest(wapQuest, "alipay.trade.wap.pay")
}
