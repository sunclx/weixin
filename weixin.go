package main

import (
	"fmt"
	"os/exec"

	"github.com/kataras/iris"
)

func main() {
	server := iris.New()

	//监听微信服务器的信息
	server.HandleFunc("", "/", mainHandle)

	//监听github.com的自动更新
	server.Post("/update", updateHandle)

	//启动服务
	server.Listen(":80")
}

func mainHandle(c *iris.Context) {
	logConect(c)

	//排除非POST请求
	if c.MethodString() != iris.MethodPost {
		c.WriteString("404")
		return
	}

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
	if !validateURL(signature, timestamp, nonce) {
		c.Log("参数错误%s,%s,%s.\n", signature, timestamp, nonce)
		c.Write("404")
		return
	}

	//过滤openid
	openid := c.URLParam("openid")
	filteOpenid(openid)

	//处理请求数据
	handlerMux(c)
}

func updateHandle(c *iris.Context) {
	logConect(c)

	//执行命令
	cmd := exec.Command("go", "get", "-u", "github.com/sunclx/weixin")
	err := cmd.Run()
	if err != nil {
		fmt.Println("errors:", err)

	}
	c.WriteString("success")
}

func logConect(c *iris.Context) {
	//记录请求
	fmt.Println(c.MethodString(), c.URI(), c.RemoteAddr())
}

func filteOpenid(openid string) {

}
