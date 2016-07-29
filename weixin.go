package main

func main() {
	//启动服务
	c := New()
	//c.Use(&Contact{})
	c.UseFunc(handleBindPhone)
	c.Run()
}
