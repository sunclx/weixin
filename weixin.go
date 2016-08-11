package main

func main() {

	//启动服务
	c := New()
	c.Handler(&contactHandler{})
	c.Run()
}
