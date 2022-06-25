package trades

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Symbol     string   `json:"s"`
	Price      float64  `json:"p"`
	Timestamp  uint64   `json:"t"`
	Conditions []string `json:"c" omitempty:"true"`
}

type Response struct {
	Type string `json:"type"`
	Data []Data `json:"data"`
}

type Subscribe struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

func (stock *Data) WriteToDisk(file *os.File) error {
	w := csv.NewWriter(file)
	err := w.Write(stock.toSlice())
	if err != nil {
		return err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (stock *Data) WriteHeaders(file *os.File) error {
	if file == nil {
		return errors.New("file handler cannot be nil")
	}
	if isFileEmpty(file) {
		w := csv.NewWriter(file)
		if bytes, _ := file.Read([]byte{}); bytes == 0 {
			err := w.Write(stock.getHeaders())
			if err != nil {
				return err
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			fmt.Printf("before log fatal %e\n", err)
			log.Fatal(err)
		}
	}
	return nil
}

func isFileEmpty(f *os.File) bool {
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 //Number of records per record. Set to Negative value for variable
	r.TrimLeadingSpace = true
	r.LazyQuotes = true
	if records, err := r.Read(); records == nil && err == io.EOF {
		return true
	}
	return false
}

func (stock *Data) getHeaders() []string {
	return []string{"Symbol", "Price", "Timestamp", "Write Timestamp"}
}

func (stock *Data) toSlice() []string {
	return []string{
		stock.Symbol,
		strconv.FormatFloat(stock.Price, 'E', -1, 64),
		strconv.FormatUint(stock.Timestamp, 10),
		strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
	}
}
