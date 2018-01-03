package grest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type router struct {
	routes WebPart
}

func (r router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(req.Context())

	defer func() {
		err := recover()
		if err != nil {
			cancel()
			p, isErrorStruct := err.(Error)
			if isErrorStruct {
				ContentType(ContentTypeJSON).
					ResponseJ(p.Status, p)(WebUnit{w, req, ctx})
			} else {
				ContentType(ContentTypeJSON).
					ResponseJ(http.StatusInternalServerError, Error{Status: http.StatusInternalServerError, Message: fmt.Sprintf("%v", err)})(WebUnit{w, req, ctx})
			}

			panic(err)
		}
	}()
	r.routes(WebUnit{w, req, ctx})
}

//StartListening starts a HTTP listener on given port
func StartListening(ctx context.Context, ip string, port uint16, routes WebPart) <-chan error {
	errChan := make(chan error)
	go func() {
		if routes == nil {
			errChan <- fmt.Errorf("no routes defined")
		} else {

			serv := http.Server{
				Addr:    ip + ":" + strconv.Itoa(int(port)),
				Handler: router{routes},
			}

			go func() {
				errChan <- serv.ListenAndServe()
			}()

			go func() {
				<-ctx.Done()
				close(errChan)
				try(serv.Close())
			}()
		}
	}()
	return errChan
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}
