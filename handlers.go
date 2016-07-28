package main

import (
	"fmt"
	"time"
)

func logHandler(c *Context) {
	r := c.Request
	fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
}

func testHandler(c *Context) {

	testID := "success-serveMux"
	c.WriteString(fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`, c.OpenID, "gh_3fb3b0b8f2fa", time.Now().Unix(), testID))

}

func messageHandler(c *Context) {
	r := c.Request
	fmt.Println(r.RemoteAddr, r.Method, r.Host, r.URL.Path, r.URL.RawQuery)
}
