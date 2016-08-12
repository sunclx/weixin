package main

import "github.com/boltdb/bolt"

func main() {
	c := New()

	c.Command("手机", func(ctx *Context) {
		if ctx.ArgsLen() != 2 {
			ctx.app.Printf("wrong")
			return
		}
		name := ctx.Arg(1)
		var openid string
		db.Update(func(tx *bolt.Tx) error {
			data := tx.Bucket([]byte("NameOpenID")).Get([]byte(name))
			openid = string(data)
			return nil
		})
		p := &PersonInfo{}
		err := p.Get(openid)
		if err != nil {
			ctx.app.Printf(`服务器错误`)
			ctx.app.LogWithError(err).Errorln("获取个人信息错误")
			return
		}
		if p.PhoneNumber == "" {
			ctx.app.Printf("没有%s的号码", name)
			return
		}

		ctx.app.Printf("%s %s", p.Name, p.PhoneNumber)

	})

	c.Run()
}
