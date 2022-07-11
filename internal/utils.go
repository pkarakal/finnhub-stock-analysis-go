package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsFileEmpty(f *os.File) bool {
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 //Number of records per record. Set to Negative value for variable
	r.TrimLeadingSpace = true
	r.LazyQuotes = true
	if records, err := r.Read(); records == nil && err == io.EOF {
		return true
	}
	return false
}

func FindItems(f string, t *time.Time) []Data {
	file, _ := os.OpenFile(f, os.O_RDWR, 0600)
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = 4
	offset := 0
	var i int64
	data := make([]Data, 100)
	d, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for idx, v := range d {
		if idx == 0 {
			continue
		}
		i, err = strconv.ParseInt(v[3], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		ts := time.UnixMilli(i).UTC()
		if ts.After(*t) && ts.After(t.Add(1*time.Minute)) {
			price, _ := strconv.ParseFloat(v[1], 64)
			if offset >= cap(data) {
				sl := make([]Data, 2*len(data), 2*len(data))
				copy(sl, data)
				data = sl
			}
			data[offset] = Data{
				Symbol:     v[0],
				Price:      price,
				Timestamp:  uint64(i),
				Conditions: nil,
			}
			offset++
		}
	}
	return data[:offset]
}

func WaitForCandlestick(f *os.File, in chan time.Time, cstick *os.File) {
	for {
		tm := <-in
		items := FindItems(f.Name(), &tm)
		if len(items) == 0 || items == nil {
			log.Println("Cannot calculate candlestick when there is no data")
		} else {
			cs := CalculateCandlestick(items)
			_ = cs.WriteToDisk(cstick)
		}
	}
}

func SanitizeString(s string) string {
	re, err := regexp.Compile(`\W`)
	if err != nil {
		log.Fatal(err)
	}
	return re.ReplaceAllString(s, "_")
}

func CreateDirs(p string) error {
	if err := os.MkdirAll(p, 0750); err != nil && !os.IsExist(err) {
		log.Fatalf("Couldn't create data directory because %v. Exiting... ", err)
	}
	return nil
}
