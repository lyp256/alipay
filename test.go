package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	o := NewWapOrder("酒店预定", "6", 199.09)
	query := NewQueryQuest("6", "")
	close := NewCloseQuest("6", "", "")
	s, _ := alipay.WapPay(o, "http://d.lyp256.cn")
	q, _ := alipay.QueryOrder(query)
	c, _ := alipay.CloseOrser(close)
	fmt.Println("create ", s)
	fmt.Println("query ", q)
	fmt.Println("close ", c)

}

var alipay *Client

func init() {
	var err error
	prib, err := ioutil.ReadFile("pri.txt")
	if err != nil {
		log.Fatalln("私钥读取失败")
	}
	pubb, err := ioutil.ReadFile("pub.txt")
	if err != nil {
		log.Fatalln("公钥读取失败")
	}
	alipay, _ = NewAlipay(prib, pubb, "2016091200494196", "https://openapi.alipaydev.com/gateway.do", ST_RSA2)
}
