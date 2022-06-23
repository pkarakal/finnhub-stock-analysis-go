package main

import (
	"encoding/json"
	cli "finnhub-stock-analysis-go/cmd"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
)

func main() {
	err := cli.Execute()
	if err != nil {
		return
	}
	fmt.Println("flags: ", cli.CLI)
	fmt.Println(fmt.Sprintf("wss://ws.finnhub.io?token=%s", cli.CLI.Token))
	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", cli.CLI.Token), nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	symbols := cli.CLI.Stocks
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
