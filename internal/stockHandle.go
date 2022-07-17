package internal

import (
	"errors"
	cli "finnhub-stock-analysis-go/cmd"
	"log"
	"os"
	"sync"
	"time"
)

type StockHandle struct {
	RollingFile        *os.File
	CandlestickFile    *os.File
	MeanFile           *os.File
	OnceFlag           *sync.Once
	StockChannel       chan time.Time
	RollingMeanChannel chan time.Time
}

func InitializeMapper(mapper map[string]*StockHandle) {
	for _, v := range cli.CLI.Stocks {
		f, s := createRollingFile(v)
		c := createCandlestickFile(v)
		m := createMeanFile(v)
		mapper[v] = &StockHandle{
			RollingFile:        f,
			CandlestickFile:    c,
			MeanFile:           m,
			OnceFlag:           s,
			StockChannel:       make(chan time.Time, 100),
			RollingMeanChannel: make(chan time.Time, 100),
		}
	}
}

func createRollingFile(stock string) (*os.File, *sync.Once) {
	var (
		file *os.File
		err  error
	)
	safeStock := SanitizeString(stock)
	file, err = os.OpenFile("data/rolling/"+safeStock+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		file, err = os.Create("data/rolling/" + safeStock + ".csv")
		if errors.Is(err, os.ErrPermission) {
			log.Fatalf("Cannot create a file due to permission reasons")
			return nil, nil
		} else {
			log.Fatalf("Couldn't create the file")
			return nil, nil
		}
	} else if err != nil {
		log.Fatalf("Couln't create file %v", err)
	}
	return file, &sync.Once{}
}

func createCandlestickFile(stock string) *os.File {
	var (
		file *os.File
		err  error
	)
	safeStock := SanitizeString(stock)
	file, err = os.OpenFile("data/candlesticks/"+safeStock+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		file, err = os.Create("data/candlesticks/" + safeStock + ".csv")
		if errors.Is(err, os.ErrPermission) {
			log.Fatalf("Cannot create a file due to permission reasons")
			return nil
		} else {
			log.Fatalf("Couldn't create the file")
			return nil
		}
	} else if err != nil {
		log.Fatalf("Couln't create file %v", err)
	}
	defer writeHeaders(file, &CandleStick{})
	return file
}

func createMeanFile(stock string) *os.File {
	var (
		file *os.File
		err  error
	)
	safeStock := SanitizeString(stock)
	file, err = os.OpenFile("data/mean/"+safeStock+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0660)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		file, err = os.Create("data/mean/" + safeStock + ".csv")
		if errors.Is(err, os.ErrPermission) {
			log.Fatalf("Cannot create a file due to permission reasons")
			return nil
		} else {
			log.Fatalf("Couldn't create the file")
			return nil
		}
	} else if err != nil {
		log.Fatalf("Couln't create file %v", err)
	}
	defer writeHeaders(file, &MeanStockData{})
	return file
}

func writeHeaders(file *os.File, c CSVAble) {
	if IsFileEmpty(file) {
		err := c.WriteHeaders(file)
		if err != nil {
			log.Fatalf("Write headers failed due to %v", err.Error())
		}
	}
}
