package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RText struct {
	ToUserName   string
	FromUserName string
	//CreateTime   time.Time
	MsgType string
	Content string
	MsgId   int64
}

var msgtmp = `
<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[gh_3fb3b0b8f2fa]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var m RText
		data, _ := ioutil.ReadAll(r.Body)
		xml.Unmarshal(data, &m)
		usrid := m.FromUserName
		fmt.Fprintf(w, msgtmp, usrid, time.Now().Unix(), "你是一个美女")

		fmt.Println(r.URL.String())
		fmt.Println(string(data))

	})
	http.ListenAndServe(":80", nil)
}
