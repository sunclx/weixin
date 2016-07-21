package main

import "github.com/kataras/iris"

func handlerMux(c *iris.Context) {
	var t msgType
	c.ReadXML(&t)
	switch t.MsgType {
	case MsgTypeImage:
		handleImage(c)
	case MsgTypeVoice:
		handleVoice(c)
	case MsgTypeVideo:
		handleVideo(c)
	case MsgTypeMusic:
		handleMusic(c)
	case MsgTypeNews:
		handleNews(c)
	default:
		c.Log("不支持该类型，%s.\n", t.MsgType)
		c.WriteString("")
	}
}

func handleText(c *iris.Context) {

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

func handleImage(c *iris.Context) {}
func handleVoice(c *iris.Context) {}
func handleVideo(c *iris.Context) {}
func handleMusic(c *iris.Context) {}
func handleNews(c *iris.Context)  {}
