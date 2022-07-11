package internal

import "os"

type CSVAble interface {
	toSlice() []string
	getHeaders() []string
	WriteHeaders(file *os.File) error
	WriteToDisk(f *os.File) error
}
