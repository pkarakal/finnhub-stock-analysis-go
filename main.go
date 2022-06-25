package main

import (
	"encoding/json"
	"errors"
	cli "finnhub-stock-analysis-go/cmd"
	"finnhub-stock-analysis-go/trades"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
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

	mapper := make(map[string]*os.File, len(cli.CLI.Stocks))
	syncMap := make(map[string]*sync.Once, len(cli.CLI.Stocks))

	closeFile := func(f *os.File) {
		f.Close()
		f.Sync()
	}
	closeFiles := func(f []*os.File) {
		for _, file := range f {
			closeFile(file)
		}
	}
	files := make([]*os.File, len(cli.CLI.Stocks))
	// create data directory
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if err := os.Mkdir("data", 0600); err != nil {
			log.Fatalf("Couldn't create data directory because %v. Exiting... ", err)
		}
	}
	for _, stock := range cli.CLI.Stocks {
		var (
			file *os.File
			err  error
		)
		file, err = os.OpenFile("data/"+stock+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			file, err = os.Create("data" + stock + ".csv")
			if errors.Is(err, os.ErrPermission) {
				log.Fatalf("Cannot create a file due to permission reasons")
			} else {
				log.Fatalf("Couldn't create the file")
			}
		} else if err != nil {
			log.Fatalf("Couln't create file %v", err)
		}
		mapper[stock] = file
		syncMap[stock] = &sync.Once{}
	}
	defer closeFiles(files)

	symbols := cli.CLI.Stocks
	for _, s := range symbols {
		msg, _ := json.Marshal(trades.Subscribe{Type: "subscribe", Symbol: s})
		err := w.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Fatalf("error writing message to websocket %e", err)
		}
	}

	var msg trades.Response
	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			log.Fatalf("Failed to decode json %e", err)
		}
		parseMessage(&msg, mapper, syncMap)
	}
}

func parseMessage(msg *trades.Response, mapper map[string]*os.File, syncMap map[string]*sync.Once) {
	for _, stock := range msg.Data {
		syncMap[stock.Symbol].Do(func() {
			stock.WriteHeaders(mapper[stock.Symbol])
		})
		err := stock.WriteToDisk(mapper[stock.Symbol])
		if errors.Is(err, os.ErrPermission) {
			log.Fatalf("Permission denied while reading or writing to file")
		}
	}
}
