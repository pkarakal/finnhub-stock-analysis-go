package main

import (
	"encoding/json"
	"errors"
	cli "finnhub-stock-analysis-go/cmd"
	"finnhub-stock-analysis-go/internal"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	err := cli.Execute()
	if err != nil {
		return
	}
	w, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://ws.finnhub.io?token=%s", cli.CLI.Token), nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	dirs := []string{"data/rolling", "data/candlesticks", "data/mean"}
	for _, dir := range dirs {
		err := internal.CreateDirs(dir)
		if err != nil {
			log.Fatalf("Couldn't create directories. %v", err)
		}
	}

	mapper := make(map[string]*internal.StockHandle, len(cli.CLI.Stocks))
	internal.InitializeMapper(mapper)

	closeFile := func(f *os.File) {
		f.Sync()
		f.Close()
	}
	closeFiles := func(f []*os.File) {
		for _, file := range f {
			closeFile(file)
		}
	}
	files := make([]*os.File, 3*len(cli.CLI.Stocks))
	for i, v := range cli.CLI.Stocks {
		files[i] = mapper[v].CandlestickFile
		files[i+len(cli.CLI.Stocks)] = mapper[v].RollingFile
		files[2*i+len(cli.CLI.Stocks)] = mapper[v].MeanFile
	}
	defer closeFiles(files)

	symbols := cli.CLI.Stocks
	for _, s := range symbols {
		msg, _ := json.Marshal(internal.Subscribe{Type: "subscribe", Symbol: s})
		err := w.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Fatalf("error writing message to websocket %e", err)
		}
	}

	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				{
					for _, v := range cli.CLI.Stocks {
						mapper[v].StockChannel <- time.Now().Truncate(time.Minute).Add(-1 * time.Minute)
						mapper[v].RollingMeanChannel <- time.Now().Truncate(time.Minute).Add(-15 * time.Minute)
					}
					break
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	for _, v := range cli.CLI.Stocks {
		go internal.WaitForCandlestick(mapper[v])
		go internal.WaitForMeanData(mapper[v])
	}

	var msg internal.Response
	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			log.Fatalf("Failed to decode json %e", err)
		}
		parseMessage(&msg, mapper)
	}
}

func parseMessage(msg *internal.Response, mapper map[string]*internal.StockHandle) {
	wg := new(sync.WaitGroup)
	wg.Add(len(msg.Data))
	for _, stock := range msg.Data {
		go func(stock internal.Data, synchro *sync.Once, w *sync.WaitGroup) {
			synchro.Do(func() {
				err := stock.WriteHeaders(mapper[stock.Symbol].RollingFile)
				if err != nil {
					log.Fatalf("Write headers once failed due to %v", err.Error())
				}
			})
			err := stock.WriteToDisk(mapper[stock.Symbol].RollingFile)
			if errors.Is(err, os.ErrPermission) {
				log.Fatalf("Permission denied while reading or writing to file")
			}
			w.Done()
		}(stock, mapper[stock.Symbol].OnceFlag, wg)
	}
	wg.Wait()
}
