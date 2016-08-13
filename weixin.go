package main

import "github.com/boltdb/bolt"

func main() {
	c := New()
	c.Command("我的姓名", func(ctx *Context) {
		if ctx.NArg() != 1 {
			ctx.Print("我的姓名 XXX")
			return
		}
		if ctx.User.Name != "" {
			ctx.Print("你的姓名是" + ctx.User.Name + "，如错误请联系管理员")
			return
		}
		ctx.User.OpenID = ctx.Message.FromUserName
		ctx.User.Name = ctx.Arg(1)
		if err := ctx.User.Put(); err == nil {
			ctx.Print("姓名设置成功")
		}
	})

	c.Command("我的学号", func(ctx *Context) {
		if ctx.NArg() != 1 {
			return
		}
		if ctx.User.StudentID != "" {
			ctx.Printf("你的学号是%s,错误请联系管理员", ctx.User.StudentID)
			return
		}
		ctx.User.OpenID = ctx.Message.FromUserName
		ctx.User.StudentID = ctx.Arg(1)
		ctx.User.Put()
		if err := ctx.User.Put(); err == nil {
			ctx.Print("学号设置成功")
		}
	})
	c.Command("我的手机", func(ctx *Context) {
		if ctx.NArg() != 1 {
			return
		}
		ctx.User.OpenID = ctx.Message.FromUserName
		ctx.User.PhoneNumber = ctx.Arg(1)
		ctx.User.Put()
		if err := ctx.User.Put(); err == nil {
			ctx.Print("手机设置成功")
		}
	})

	c.Command("手机", func(ctx *Context) {
		if ctx.NArg() != 1 {
			return
		}
		name := ctx.Arg(0)
		var openid string
		db.Update(func(tx *bolt.Tx) error {
			data := tx.Bucket([]byte("NameOpenID")).Get([]byte(name))
			openid = string(data)
			return nil
		})
		p := &User{}
		err := p.Get(openid)
		if err != nil {
			ctx.Printf(`服务器错误`)
			return
		}
		if p.PhoneNumber == "" {
			ctx.Printf("没有%s的号码", name)
			return
		}

		ctx.Printf("%s %s", p.Name, p.PhoneNumber)
	})

	c.Run()
}
