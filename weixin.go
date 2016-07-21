package main

import (
	"fmt"

	"github.com/kataras/iris"
)

func handleError(err error) {

}
func requestHandle(c *iris.Context) []byte {
	return []byte{}

	// text, err := unmarshalMsg(r)
	// if err != nil {
	// 	handleError(err)
	// 	return nil
	// }
	// fmt.Println(text)
	// resMsg := func(text Text) (s string) {

	// 	db, err := bolt.Open("data.db", 0600, nil)
	// 	if err != nil {
	// 		return ""
	// 	}
	// 	defer db.Close()
	// 	content := text.Content
	// 	if strings.HasPrefix(content, "我的学号是") {
	// 		content = content[len(content)-8:]
	// 		db.Update(func(tx *bolt.Tx) error {
	// 			b := tx.Bucket([]byte("default"))
	// 			err := b.Put([]byte(text.FromUserName), []byte(content))
	// 			s = fmt.Sprintf("你的学号是%s，你是%s", content, "我们班的同学")

	// 			return err
	// 		})
	// 	}

	// 	db.Update(func(tx *bolt.Tx) error {
	// 		b := tx.Bucket([]byte("default"))
	// 		data := b.Get([]byte(text.FromUserName))
	// 		if data == nil {
	// 			s = `请输入"我的学号是00000000"`
	// 			return nil
	// 		}

	// 		if string(data) == "09170515" {
	// 			s = "你是跳跳，一个大美女"
	// 			return nil
	// 		}
	// 		if string(data) == "09170512" {
	// 			s = "你是乐乐，一个大美女"
	// 			return nil
	// 		}
	// 		s = fmt.Sprintf("你的学号是%s，你是%s", data, "我们班的同学")

	// 		return nil
	// 	})

	// 	return s
	// }(text)

	// return []byte(marshaMsg(text.FromUserName, resMsg))

}

func main() {
	server := iris.New()

	server.HandleFunc("", "/", func(c *iris.Context) {
		//记录请求
		fmt.Println(c.MethodString(), c.URI())

		//建议域名是否真确
		if hostname := c.HostString(); hostname != "weixin.chenlixin.net" {
			c.Log("异常域名:", hostname)
			c.Write("404")
			return
		}

		//检验是否是微信服务器的请求
		if !validateURL(c.Params) {
			fmt.Println("验证错误", c.Params)
			c.Write("404")
			return
		}

		//处理请求数据
		data := c.PostBody()
		fmt.Printf("%s\n", data)
		requestHandle(c)

		var msg Text
		c.ReadXML(&msg)
		fmt.Printf("%s\n", data)

		c.Write("")

	})
	server.Listen(":80")
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.Method, r.RequestURI, r.RemoteAddr)
	// 	if !validate(r) {
	// 		w.WriteHeader(404)
	// 		w.Write([]byte("404"))
	// 		return
	// 	}

	// 	data := requestHandle(r)
	// 	w.Write(data)

	// })

	// http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println(r.Method, r.RequestURI, r.RemoteAddr)
	// 	cmd := exec.Command("go", "get", "-u", "github.com/sunclx/weixin")
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		fmt.Println("errors:")
	// 		fmt.Println(err)
	// 	}

	// })
	// http.ListenAndServe(":80", nil)
}
