package main

func main() {
	//启动服务
	c := &Context{}
	c.UseFunc(handlePhone)
	c.UseFunc(handleBindPhone)
	c.Run()
}
