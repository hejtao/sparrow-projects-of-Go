package ipc

import (
	"encoding/json"
	"fmt"
)

//定义request,response和Server接口,IpcServer封装了Server
//IpcServer拥有一个Connec方法，通过chan string接收request，交给server.Handle(),并返回responce

type Request struct { //服务器请求
	Method string "method"
	Params string "params"
}

type Response struct { //服务器响应
	Code string "code"
	Body string "body"
}

type Server interface {
	Name() string                           //命名
	Handle(method, params string) *Response //请求处理函数
}

type IpcServer struct {
	Server //匿名字段
}

func NewIpcServer(server Server) *IpcServer { //将Server封装成IpcServer
	return &IpcServer{server}
}

func (server *IpcServer) Connect() chan string { //通过chan string连接到一个goroutine，传入请求返回响应
	session := make(chan string)

	go func(c chan string) {
		for {
			request := <-c

			if request == "CLOSE" {
				break
			}

			var req Request
			err := json.Unmarshal([]byte(request), &req) //string转[]byte转Request{}
			if err != nil {
				fmt.Println("Invalid request format:", request)
			}

			resp := server.Handle(req.Method, req.Params) //处理request
			b, err := json.Marshal(resp)                  //将Response{}转[]byte
			c <- string(b)
		}

		fmt.Println("Session closed.")

	}(session)

	fmt.Println("A new session has been created successfully.")

	return session
}
