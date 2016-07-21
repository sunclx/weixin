package main

import (
	"fmt"
	"os/exec"

	"github.com/kataras/iris"
)

func handleError(err error) {

}

func main() {
	server := iris.New()

	server.Post("/", func(c *iris.Context) {
		//记录请求
		fmt.Println(c.MethodString(), c.URI(), c.RemoteAddr())

		//建议域名是否真确
		if hostname := c.HostString(); hostname != "weixin.chenlixin.net" {
			c.Log("异常域名:", hostname)
			c.Write("404")
			return
		}

		//检验是否是微信服务器的请求
		signature := c.URLParam("signature")
		timestamp := c.URLParam("timestamp")
		nonce := c.URLParam("nonce")
		//openid:=c.URLParam("openid")
		if !validateURL(signature, timestamp, nonce) {
			c.Log("参数错误%s,%s,%s.\n", signature, timestamp, nonce)
			c.Write("404")
			return
		}

		//处理请求数据
		data := c.PostBody()
		c.Log("%s\n", data)
		requestHandle(c)

		var msg Text
		c.ReadXML(&msg)
		c.Log("%v\n", data)

		c.Write("")

	})

	//自动更新
	server.Post("/update", func(c *iris.Context) {
		//记录请求
		fmt.Println(c.MethodString(), c.URI(), c.RemoteAddr())

		//执行命令
		cmd := exec.Command("go", "get", "-u", "github.com/sunclx/weixin")
		err := cmd.Run()
		if err != nil {
			fmt.Println("errors:", err)

		}
	})

	server.Listen(":80")

	// http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.Method, r.RequestURI, r.RemoteAddr)

}
