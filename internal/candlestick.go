package internal

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
	"time"
)

type CandleStick struct {
	OpenPrice         float64
	ClosePrice        float64
	HighestPrice      float64
	LowestPrice       float64
	TotalTransactions int
	StockSymbol       string
	MinuteOfDay       time.Time
}

func CalculateCandlestick(t []Data) *CandleStick {
	min, max := getMinMax(t)
	return &CandleStick{
		OpenPrice:         t[0].Price,
		ClosePrice:        t[len(t)-1].Price,
		HighestPrice:      max,
		LowestPrice:       min,
		TotalTransactions: len(t),
		StockSymbol:       t[0].Symbol,
		MinuteOfDay:       time.UnixMilli(int64(t[0].Timestamp)).Truncate(1 * time.Minute),
	}
}

func getMinMax(p []Data) (min, max float64) {
	min = p[0].Price
	max = p[0].Price

	for _, v := range p {
		if v.Price < min {
			min = v.Price
		}
		if v.Price > max {
			max = v.Price
		}
	}
	return min, max
}

func (c *CandleStick) toSlice() []string {
	s := SanitizeString(c.StockSymbol)
	return []string{
		s,
		c.MinuteOfDay.String(),
		strconv.FormatFloat(c.OpenPrice, 'f', -1, 64),
		strconv.FormatFloat(c.ClosePrice, 'f', -1, 64),
		strconv.FormatFloat(c.HighestPrice, 'f', -1, 64),
		strconv.FormatFloat(c.LowestPrice, 'f', -1, 64),
		strconv.FormatInt(int64(c.TotalTransactions), 10),
	}
}
func (c *CandleStick) getHeaders() []string {
	return []string{"Symbol", "MinuteOfDay", "OpenPrice", "ClosePrice", "HighestPrice", "LowestPrice", "Transactions"}
}

func (c *CandleStick) WriteToDisk(f *os.File) error {
	w := csv.NewWriter(f)
	s := c.toSlice()
	err := w.Write(s)
	if err != nil {
		return err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (c *CandleStick) WriteHeaders(file *os.File) error {
	if file == nil {
		return errors.New("file handler cannot be nil")
	}
	if IsFileEmpty(file) {
		w := csv.NewWriter(file)
		if bytes, _ := file.Read([]byte{}); bytes == 0 {
			err := w.Write(c.getHeaders())
			if err != nil {
				return err
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
