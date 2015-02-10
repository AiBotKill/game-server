# game-server

Download, compile, and start the game-server with:
```bash
go get github.com/AiBotKill/game-server
~/go/bin/game-server
```

## NATS communication

#### "registerAi"

Request

```go
type RegisterAiMsg struct {
	BotId string `json:"botId"`
}
```

Reply

```go
type IdReplyMsg struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}
```

#### createGame

Request

```go
type CreateGameMsg struct {
	TimeLimit         time.Duration `json:"timelimit"`
	GameArea          [2]float64    `json:"gameArea"`
	Tiles             []*tile       `json:"tiles"`
	StartingPositions []*Vector     `json:"startingPositions"`
}
```

Reply

```go
type IdReplyMsg struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}
```

#### <gameId>.createPlayer

Request

```go
type CreatePlayerMsg struct {
	BotId string `json:"botId"`
	Name  string `json:"name"`
}
```

Reply

```go
type IdReplyMsg struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}
```

#### <gameId>.start

Request

```go
type StartGameMsg struct {
}
```

Reply

```go
type IdReplyMsg struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}
```

#### <playerId>.action

Request

Reply

```go
type ActionMsg struct {
	Type      string  `json:"type"`
	Direction *Vector `json:"direction"`
}
```

Reply

```go
type IdReplyMsg struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Error  string `json:"error,omitempty"`
}
```
