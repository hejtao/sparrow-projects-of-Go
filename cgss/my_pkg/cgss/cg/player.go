package cg

import (
	"fmt"
)

type Player struct {
	Name  string "name"
	Level int    "level"
	Exp   int    "exp"
	Room  int    "room"

	mq chan *Message
}

//生成一个Player{},包含一个goroutine等待接收来自 chan *Message 的消息
func NewPlayer() *Player {
	m := make(chan *Message, 1024)
	p := &Player{" ", 0, 0, 0, m}

	go func(p *Player) { // 为每个player{}启动一个goroutine
		msg := <-p.mq // 等待来自通道 mq 的消息
		fmt.Println(p.Name, "received message:", msg.Content)
	}(p)

	return p
}
