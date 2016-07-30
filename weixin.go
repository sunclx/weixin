package main

func main() {
	//启动服务
	c := New()
	c.Use(&contactHandler{})
	c.UseFunc(handleBindPhone)
	c.Run()
}
