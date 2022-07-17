package internal

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
	"time"
)

type MeanStockData struct {
	TotalTransactions int
	MeanPrice         float64
	StockSymbol       string
	StartTime         time.Time
	EndTime           time.Time
}

func CalculateMeanStockData(t []Data, s, e time.Time) *MeanStockData {
	return &MeanStockData{
		TotalTransactions: len(t),
		StockSymbol:       t[0].Symbol,
		MeanPrice:         getMeanPrice(t),
		StartTime:         s.Truncate(1 * time.Minute),
		EndTime:           e.Truncate(1 * time.Minute),
	}
}

func getMeanPrice(p []Data) (mean float64) {
	sum := 0.0
	for _, v := range p {
		sum += v.Price
	}
	return sum / float64(len(p))
}

func (m *MeanStockData) toSlice() []string {
	s := SanitizeString(m.StockSymbol)
	return []string{
		s,
		m.StartTime.String(),
		m.EndTime.String(),
		strconv.FormatFloat(m.MeanPrice, 'f', -1, 64),
		strconv.FormatInt(int64(m.TotalTransactions), 10),
	}
}
func (m *MeanStockData) getHeaders() []string {
	return []string{"Symbol", "StartTime", "EndTime", "MeanPrice", "Transactions"}
}

func (m *MeanStockData) WriteToDisk(f *os.File) error {
	w := csv.NewWriter(f)
	s := m.toSlice()
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

func (m *MeanStockData) WriteHeaders(file *os.File) error {
	if file == nil {
		return errors.New("file handler cannot be nil")
	}
	if IsFileEmpty(file) {
		w := csv.NewWriter(file)
		if bytes, _ := file.Read([]byte{}); bytes == 0 {
			err := w.Write(m.getHeaders())
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
