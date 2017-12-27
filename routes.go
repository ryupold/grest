package grest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type test struct {
	A string
	B string
}

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
func StartListening(ctx context.Context, ip string, port uint16, routes WebPart) error {
	if routes == nil {
		return fmt.Errorf("no routes defined")
	}

	serv := http.Server{
		Addr:    ip + ":" + strconv.Itoa(int(port)),
		Handler: router{routes},
	}

	go func() {
		<-ctx.Done()
		try(serv.Close())
	}()

	return nil
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}
