# grest

Functional HTTP/REST library in Go

## usage

```go
import (
    "context"
    "os"
    . "github.com/ryupold/grest"
)

func main() {
    ctx, stopServer := context.WithCancel(context.Background())
    go StartListening(ctx, "0.0.0.0", 8080, Choose(
        Path("/").OKS("hello world"),
        NotFoundS("404"),
    ))

    os.Stdin.Read([]byte{0})
    stopServer()
}
```