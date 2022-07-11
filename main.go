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

	mapper := make(map[string]*os.File, len(cli.CLI.Stocks))
	syncMap := make(map[string]*sync.Once, len(cli.CLI.Stocks))
	candlestickMap := make(map[string]*os.File, len(cli.CLI.Stocks))
	stockChans := make(map[string]chan time.Time, len(cli.CLI.Stocks))

	for _, v := range cli.CLI.Stocks {
		stockChans[v] = make(chan time.Time, 100)
	}

	closeFile := func(f *os.File) {
		f.Sync()
		f.Close()
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
	createFiles(mapper, syncMap)
	createCstickFiles(candlestickMap)
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
						stockChans[v] <- time.Now().Truncate(time.Minute).Add(-1 * time.Minute)
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
		go internal.WaitForCandlestick(mapper[v], stockChans[v], candlestickMap[v])
	}

	var msg internal.Response
	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			log.Fatalf("Failed to decode json %e", err)
		}
		parseMessage(&msg, mapper, syncMap)
	}
}

func createFiles(mapper map[string]*os.File, syncMap map[string]*sync.Once) {
	for _, stock := range cli.CLI.Stocks {
		var (
			file *os.File
			err  error
		)
		safeStock := internal.SanitizeString(stock)
		file, err = os.OpenFile("data/rolling/"+safeStock+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			file, err = os.Create("data/rolling/" + safeStock + ".csv")
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

}

func createCstickFiles(candlestickMap map[string]*os.File) {
	for _, v := range cli.CLI.Stocks {
		var (
			file *os.File
			err  error
		)
		safeStock := internal.SanitizeString(v)
		file, err = os.OpenFile("data/candlesticks/"+safeStock+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			file, err = os.Create("data/candlesticks/" + safeStock + ".csv")
			if errors.Is(err, os.ErrPermission) {
				log.Fatalf("Cannot create a file due to permission reasons")
			} else {
				log.Fatalf("Couldn't create the file")
			}
		} else if err != nil {
			log.Fatalf("Couln't create file %v", err)
		}
		candlestickMap[v] = file
		if internal.IsFileEmpty(file) {
			cs := &internal.CandleStick{}
			err := cs.WriteHeaders(file)
			if err != nil {
				log.Fatalf("Write headers failed due to %v", err.Error())
			}
		}
	}
}

func parseMessage(msg *internal.Response, mapper map[string]*os.File, syncMap map[string]*sync.Once) {
	wg := new(sync.WaitGroup)
	wg.Add(len(msg.Data))
	for _, stock := range msg.Data {
		go func(stock internal.Data, syncMap map[string]*sync.Once, w *sync.WaitGroup) {
			syncMap[stock.Symbol].Do(func() {
				err := stock.WriteHeaders(mapper[stock.Symbol])
				if err != nil {
					log.Fatalf("Write headers once failed due to %v", err.Error())
				}
			})
			err := stock.WriteToDisk(mapper[stock.Symbol])
			if errors.Is(err, os.ErrPermission) {
				log.Fatalf("Permission denied while reading or writing to file")
			}
			w.Done()
		}(stock, syncMap, wg)
	}
	wg.Wait()
}
