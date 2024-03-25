package ws

import (
	"syscall/js"
)

type Sender struct {
	websocketClient js.Value
}

func (receiver *Sender) Send(playerID int64, gameMessage string) {
	receiver.websocketClient.Call("send", playerID, gameMessage)
}

func (receiver *Sender) Broadcast(gameMessage string) {
	receiver.websocketClient.Call("broadcast", gameMessage)
}
