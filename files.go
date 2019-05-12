package grest

import (
	"fmt"
	"io"
	"os"
)

//ServeFile tries to serve the file
func ServeFile(file string) WebPart {
	return func(u WebUnit) *WebUnit {
		return ServeReadCloser(func(WebUnit) (io.ReadCloser, error) {
			return os.Open(file)
		})(u)
	}
}

//ServeFile tries to serve the file
func (w WebPart) ServeFile(file string) WebPart {
	return Compose(w, ServeFile(file))
}

//ServeFolder serves a folder overview similar to Pythons `python -m http.server`, where you would see an overview of the files in the folder and can navigate in it opening files that are served with ServeFile(...)
func ServeFolder(path string) WebPart {
	return func(u WebUnit) *WebUnit {
		return Panic(fmt.Errorf("ServeFolder is not implemented yet"))(u)
	}
}

//ServeFolder serves a folder overview similar to Pythons `python -m http.server`, where you would see an overview of the files in the folder and can navigate in it opening files that are served with ServeFile(...)
func (w WebPart) ServeFolder(path string) WebPart {
	return Compose(w, ServeFolder(path))
}
