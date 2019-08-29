package party

import(
	"../github.com/gorilla/websocket"
	"../gamer"
	"../helpers"
	"time"
)

const(
	deltaTime = 1.0/60.0
	paddleHeight = 90
)

type Party struct {
	Name string;
	LeftPlayer *gamer.Player;
	RightPlayer *gamer.Player;
	BallInfo helpers.Ball;
	Restart int;
	GameLaunch bool;
	
	Rest chan string;
	Players chan string;
}

var AllParties = make(map[string]*Party)
var FreeParties = make(map[string]*Party)

func (p *Party)Reset(){
	p.BallInfo = helpers.GetBall()
	p.Restart = 0
	p.GameLaunch = false
	if p.LeftPlayer != nil {
		p.LeftPlayer.Reset()
	}
	if p.RightPlayer != nil{
		p.RightPlayer.Reset()
	}
}

func (p *Party)Supervisor() {
	for range time.Tick(time.Second / 60){
		if p.GameLaunch{
			p.MoveBall()
			p.LeftPlayer.Connection.WriteJSON(helpers.Game{"game", p.RightPlayer.Yposition, p.BallInfo.X, p.BallInfo.Y})
			p.RightPlayer.Connection.WriteJSON(helpers.Game{"game", p.LeftPlayer.Yposition, p.BallInfo.X, p.BallInfo.Y})
		}

		select{
		case c := <-p.Rest:
			if c == "restart"{
				if p.Restart != 1{
					p.Restart++
				}else{
					p.Reset()
					p.SendToAll(helpers.Ready{"ready"});
				}
			} else if c == "close"{
				return
			} else if c == "start" {
				p.GameLaunch = true;
			}
		break;
		default:
		break;
		}
	}
}

func (p *Party)AddPlayer(conn *websocket.Conn){
	if p.LeftPlayer == nil {
		p.LeftPlayer = &gamer.Player{conn, 0, 386}
		conn.WriteJSON(helpers.Init{"init", true})
		go p.LPlayerListener()
	} else {
		p.RightPlayer = &gamer.Player{conn, 0, 386}
		conn.WriteJSON(helpers.Init{"init", false})


		go p.RPlayerListener()
		p.SendToAll(helpers.Ready{"ready"});
	}
}

func NewConnection(conn *websocket.Conn) {
	if len(FreeParties) != 0 {
		for _, p := range FreeParties {
			p.AddPlayer(conn)
			AllParties[p.Name] = p
			delete(FreeParties, p.Name)
			break
		}
	} else {
		newParty := &Party{
			Name:        helpers.RandString(6),
			LeftPlayer:  nil,
			RightPlayer: nil,
			BallInfo:    helpers.GetBall(),
			Restart:     0,
			Rest:        make(chan string),
			Players:     make(chan string),
			GameLaunch:  false,
		}
		go newParty.Supervisor()
		newParty.AddPlayer(conn)
		FreeParties[newParty.Name] = newParty
	}
}

func (p* Party)MoveBall(){
	if p.BallInfo.X < 10{
		p.SendScore(false)
	}

	if p.BallInfo.X > 790{
		p.SendScore(true)
	}

	if p.BallInfo.Y < 135 || p.BallInfo.Y > 717{
		p.BallInfo.Dy *= -1
	}

	if (p.BallInfo.X < 27 && p.BallInfo.X > 22) && (p.BallInfo.Y > p.LeftPlayer.Yposition && p.BallInfo.Y < p.LeftPlayer.Yposition + paddleHeight){
		p.BallInfo.Dx *= -1
	}
	if (p.BallInfo.X > 765 && p.BallInfo.X < 770) && (p.BallInfo.Y > p.RightPlayer.Yposition && p.BallInfo.Y < p.RightPlayer.Yposition + paddleHeight){
		p.BallInfo.Dx *= -1
	}

	p.BallInfo.Move(deltaTime)
}

func (p *Party)SendScore(pl bool){
	if pl{
		p.LeftPlayer.Score++;
	} else {
		p.RightPlayer.Score++;
	}

	p.BallInfo.X = 396;
	p.BallInfo.Y = 426;

	if p.LeftPlayer.Score == 5 || p.RightPlayer.Score == 5{
		p.GameLaunch = false
	}
	p.SendToAll(helpers.Score{"score",pl,p.LeftPlayer.Score, p.RightPlayer.Score})
}

func (p Party)SendToAll(mess interface{}){
	p.LeftPlayer.Connection.WriteJSON(mess)
	p.RightPlayer.Connection.WriteJSON(mess)
}

func (p *Party)LPlayerListener() {
	defer func(){
		if p.LeftPlayer == nil && p.RightPlayer == nil {
			p.Rest <- "close"
			if _, ok := AllParties[p.Name]; ok{
				delete(AllParties, p.Name)
			} else {
				delete(FreeParties, p.Name)
			}
		}
	}()

	for{
		if p.LeftPlayer == nil {
			return
		}
		var mess interface{}

		err := p.LeftPlayer.Connection.ReadJSON(&mess)

		if err != nil{
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				p.GameLaunch = false

				if p.RightPlayer != nil {
					p.RightPlayer.Connection.WriteJSON(struct{Command string}{"opponentLeft"})
				}
			}
			p.LeftPlayer = nil
			return
		}

		if mess == nil{
			continue
		}
		info := mess.(map[string]interface{})

		switch info["Command"].(string) {
		case "coords":
			p.LeftPlayer.Yposition = info["Player"].(float64)
		break;
		case "readyState":
			p.Rest <- "restart"
		break;
		case "start":
			p.Rest <- "start"
		break;
		case "findNew":
			NewConnection(p.LeftPlayer.Connection)
			p.LeftPlayer = nil
			return
		break;
		}
	}
}

func (p *Party)RPlayerListener() {
	defer func(){
		if p.LeftPlayer == nil && p.RightPlayer == nil {
			p.Rest <- "close"
			if _, ok := AllParties[p.Name]; ok{
				delete(AllParties, p.Name)
			} else {
				delete(FreeParties, p.Name)
			}
		}
	}()

	for{
		if p.RightPlayer == nil {
			return
		}
		var mess interface{}

		err := p.RightPlayer.Connection.ReadJSON(&mess)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				p.GameLaunch = false

				if p.LeftPlayer != nil{
					p.LeftPlayer.Connection.WriteJSON(struct{Command string}{"opponentLeft"})
				}
			}
			p.RightPlayer = nil
			return
		}

		if mess == nil{
			continue
		}
		info := mess.(map[string]interface{})

		switch info["Command"].(string) {
		case "coords":
			p.RightPlayer.Yposition = info["Player"].(float64)
		break;
		case "readyState":
			p.Rest <- "restart"
		break;
		case "start":
			p.Rest <- "start"
		break;
		case "findNew":
			NewConnection(p.RightPlayer.Connection)
			p.RightPlayer = nil
			return
		break;
		}
	}
}