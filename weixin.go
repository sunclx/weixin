package main

import (
	"fmt"
	"net/http"
)

func handleError(err error) {

}
func requestHandle(r *http.Request) []byte {

	text, err := unmarshalMsg(r)
	if err != nil {
		handleError(err)
		return nil
	}
	msg := marshaMsg(text.FromUserName, "跳跳和乐乐都是美女")
	return []byte(msg)
}
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := requestHandle(r)
		w.Write(data)

		fmt.Println(r.URL.String())
		fmt.Println(string(data))

	})
	http.ListenAndServe(":80", nil)
}
