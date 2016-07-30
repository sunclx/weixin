package main

func main() {

	//启动服务
	c := New()
	c.Use(&openidHandler{})
	c.Use(&contactHandler{})
	c.Run()
}
