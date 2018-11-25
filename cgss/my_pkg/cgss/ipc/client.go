package ipc

import (
	"encoding/json"
)

type IpcClient struct {
	conn chan string
}

//由一个IpcServer生成一个IpcClient,且连接一个chan string
func NewIpcClient(server *IpcServer) *IpcClient {
	c := server.Connect()

	return &IpcClient{c}
}

//从IpcClient{}端传入一个Request{},并返回Response{}
func (client *IpcClient) Call(method, params string) (resp *Response, err error) {
	req := &Request{method, params} // 一个请求实例

	var b []byte
	b, err = json.Marshal(req) //Request{}转[]byte
	if err != nil {
		return
	}

	client.conn <- string(b) //[]byte转string
	str := <-client.conn

	var resp1 Response
	err = json.Unmarshal([]byte(str), &resp1) //还原成请求类实体
	resp = &resp1

	return
}

func (client *IpcClient) close() {
	client.conn <- "CLOSE"
}
