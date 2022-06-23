package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func main() {
	w, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token=capm2q2ad3i1rqbdbqk0", nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	symbols := []string{"AAPL", "AMZN", "MSFT", "BINANCE:BTCUSDT", "GTLB"}
	for _, s := range symbols {
		msg, _ := json.Marshal(Subscribe{"subscribe", s})
		err := w.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing message to websocket %e", err)
		}
	}

	var msg Response
	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error", msg)
			panic(err)
		}
		fmt.Println(msg, msg.Type)
	}
}
