package cg

import (
	"encoding/json"
	"errors"
	"sync"

	"my_pkg/cgss/ipc"
)

//CenterServer实现了接口Server，可以接收request并返回response

type Message struct {
	From    string "from"
	To      string "to"
	Content string "content"
}

type CenterServer struct {
	servers map[string]ipc.Server
	players []*Player //玩家数据库
	//rooms   []*Room
	mutex sync.RWMutex
}

func NewCenterServer() *CenterServer {
	servers := make(map[string]ipc.Server)
	players := make([]*Player, 0)

	return &CenterServer{servers: servers, players: players}
}

func (server *CenterServer) addPlayer(params string) error { //添加一个player{}

	player := NewPlayer()

	err := json.Unmarshal([]byte(params), &player)
	if err != nil {
		return err
	}

	for _, v := range server.players { // 若玩家已经登陆，直接返回
		if v == player {
			return nil
		}
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.players = append(server.players, player)

	return nil
}

func (server *CenterServer) removePlayer(params string) error {
	server.mutex.Lock()
	defer server.mutex.Unlock() // 使用defer解锁更优雅！！！

	for i, v := range server.players {
		if v.Name == params {
			if len(server.players) == 1 { // 删除唯一玩家
				server.players = make([]*Player, 0)
			} else if i == len(server.players)-1 { // 删除最后一玩家
				server.players = server.players[:i-1]
			} else if i == 0 { // 删除第一个玩家
				server.players = server.players[1:]
			} else {
				server.players = append(server.players[:i-1], server.players[:i+1]...)
			}
			return nil
		}
	}

	return errors.New("Player not found.")
}

func (server *CenterServer) listPlayer(params string) (players string, err error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if len(server.players) > 0 {
		b, _ := json.Marshal(server.players) //
		players = string(b)
	} else {
		err = errors.New("No players online.")
	}

	return
}

func (server *CenterServer) broadcast(params string) error {
	var message Message
	err := json.Unmarshal([]byte(params), &message) //string转[]byte转Message{}

	if err != nil {
		return err
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if len(server.players) > 0 {
		for _, player := range server.players { // 给每位玩家发布消息
			player.mq <- &message // 创建每位玩家时都有一个goroutine
		} // 在等待接收消息
	} else {
		err = errors.New(("No players online."))
	}

	return err
}

func (server *CenterServer) Handle(method, params string) *ipc.Response {
	switch method {
	case "addplayer":
		err := server.addPlayer(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
	case "removeplayer":
		err := server.removePlayer(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
	case "listplayer":
		players, err := server.listPlayer(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{"200", players}
	case "broadcast":
		err := server.broadcast(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{Code: "200"}
	default:
		return &ipc.Response{Code: "404", Body: method + ":" + params}
	}
	return &ipc.Response{Code: "200"}
}

func (server *CenterServer) Name() string {
	return "CenterServer"
}
