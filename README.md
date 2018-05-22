# golang alipay 支付宝支付接口  
## 已完成支付接口:
* `alipay.trade.wap.pay`[手机网站支付接口](https://docs.open.alipay.com/203)
* `alipay.trade.page.pay`[电脑网站支付接口](https://docs.open.alipay.com/270/alipay.trade.page.pay)
* `alipay.trade.close`[统一收单交易关闭接口](https://docs.open.alipay.com/api_1/alipay.trade.close/)
* `alipay.trade.query`[统一收单线下交易查询](https://docs.open.alipay.com/api_1/koubei.trade.itemorder.query/)
* `alipay.trade.refund`[统一收单交易退款接口](https://docs.open.alipay.com/api_1/alipay.trade.refund/)
* `alipay.trade.fastpay.refund.query`[统一收单交易退款查询接口](https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query/)
* `alipay.data.dataservice.bill.downloadurl.query`[查询对账单下载地址](https://docs.open.alipay.com/api_15/alipay.data.dataservice.bill.downloadurl.query)

[支持异步通知验签](https://docs.open.alipay.com/203/105286)

支付宝应用创建及配置参见支付宝[快速接入](https://docs.open.alipay.com/203/105285/)介绍,也可使用支付宝[沙箱环境](https://docs.open.alipay.com/200/105311/)进行测试开发
## 安装
    get github.com/lyp256/alipay
## 使用
### 导入`github.com/lyp256/alipay`包
    import "github.com/lyp256/alipay"
### 创建一个alipay客户端
    priByte, err := ioutil.ReadFile("pri.txt")
	if err != nil {
		log.Fatalln("私钥文件读取失败")
	}
	pubByte, err := ioutil.ReadFile("pub.txt")
	if err != nil {
		log.Fatalln("公钥文件读取失败")
	}
	payClick, _ = alipay.NewAlipay(priByte, pubByte, "APPID", "支付宝网关地址", alipay.ST_RSA2)
	payClick.SetNotifyUrl("默认通知地址")
### 说明
所有请求分三步完成
1. 创建业务:调用alipay.NewWapPay,alipay.NewPagePay,alipay.NewQuery等方法传入必要参数会返回业务struct,调用业务struct的SetXXXX方法(如:SetTimeExpire,SetBody等)可改变业务参数.
  
2. 创建请求:调用alipay.WapPay,alipay.PagePay等方法传入业务struct返回请求struct,
调用该struct的SetXXXX方法(如:SetReturnUrl,SetNotifyUrl等)可以修改公共参数
3. 调用请求struct的Build()方法返回请求链接

创建业务请求时传入必要参数
### 手机网站支付
    //创建一个手机支付订单
    wapOrder,err := alipay.NewWapPay("唯一订单号", "标题", 9.9 /*价格*/ )
    if err!=nil{
        fmt.Println(err)
    }
    //设置支付超时时间
	wapOrder.SetTimeExpire(time.Now().Unix() + 180)
    //设置回跳地址并生成支付链接
	url, _ := pay.WapPay(o).SetReturnUrl("http://lyp256.cn").Build()
	fmt.Println("手机网站支付链接:", url)
### 电脑网站支付
	pageOrder, _ := alipay.NewPagePay("唯一订单号", "标题", 9.9)
	url, _ := pay.PagePay(pageOrder).Build()
	fmt.Println("电脑网站支付链接:", url)
### 查询订单
    //创建查询,两个订单号选择一个传入即可
	query, _ := alipay.NewQuery("你自己创建订单时传入的订单号", "支付宝返回的订单号")
    //生成查询地址
	url, _ := pay.QueryOrder(query).Build()
    //查询订单返回查询结果map,出错或者验签失败会返回空
    resultMap,err:=pay.QueryOrderParams(query)
    if err!=nil{
	fmt.Println(err)
    }
	fmt.Println("查询地址:", url)
	fmt.Println("查询结果:",resultMap)
### 异步通知处理
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		err:= pay.ValidateNotify(request.Form)
        if err!=nil{
            /*验证支付宝签名失败*/
        }else{
            /*验证支付宝签成功*/

            //业务处理代码
            some code......
            //业务处理完成,返回success标识符.注意返回成功标志前后不可返回其他任意数据.
           alipay.NotIfySuccess(writer)
        }
	})
	http.ListenAndServe(":80", nil)
## 贡献
欢迎[提问](https://github.com/lyp256/alipay/issues/new)或Pull Request 联系我[lyp256@gmail.com](mailto:lyp256@gmail.com)