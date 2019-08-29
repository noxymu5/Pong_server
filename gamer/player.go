package gamer

import(
	"../github.com/gorilla/websocket"
)

type Player struct{
	Connection *websocket.Conn;
	Score int;
	Yposition float64;
}

func (p *Player)Reset(){
	p.Score = 0;
	p.Yposition = 386;
}