package grest

import (
	"context"
	"fmt"
	"os"
)

func ExampleRouteSetup() {
	ctx, stopServer := context.WithCancel(context.Background())
	go func() {
		serverLife := StartListening(ctx, "", 44444, Choose(
			// /
			Path("/").OK().ServeString("hello world"),

			// /hello/STRING
			TypedPath("/hello/%s", func(u WebUnit, params []interface{}) *WebUnit {
				return OK().ServeString(fmt.Sprintf("hello %s", params[0]))(u)
			}),

			// /add/NUMBER/NUMBER
			TypedPath("/add/%d/%d", func(u WebUnit, params []interface{}) *WebUnit {
				n1 := params[0].(int)
				n2 := params[1].(int)
				return OK().ServeString(fmt.Sprintf("%d + %d = %d", n1, n2, n1+n2))(u)
			}),

			//404 - no route matched
			NotFound().ServeString("404"),
		))

		select {
		case err, alive := <-serverLife:
			if err != nil {
				panic(err)
			} else if !alive {
				fmt.Println("server stopped")
			}
		}

		fmt.Println("server listening...")
	}()

	fmt.Println("press ENTER to stop server")
	os.Stdin.Read([]byte{0})
	stopServer()
}
