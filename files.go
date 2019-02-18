package grest

import (
	"io"
	"os"
)

//ServeFile tries to serve the file
func ServeFile(file string) WebPart {
	return func(u WebUnit) *WebUnit {
		return ServeReadCloser(func() (io.ReadCloser, error) {
			return os.Open(file)
		})(u)
	}
}
