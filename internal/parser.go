package internal

type CSVAble interface {
	toSlice() []string
	getHeaders() []string
}
