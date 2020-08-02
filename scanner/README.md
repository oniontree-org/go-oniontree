# Scanner

Scanner is a concurrent, configurable TCP scanner for OnionTree content.

## Example

```go
package main

import (
    "fmt"
    "context"
    "github.com/onionltd/go-oniontree"
    "github.com/onionltd/go-oniontree/scanner"
)

func main() {
    s := scanner.NewScanner()
    s.Start()
    s.Events()

    eventCh := make(chan events.Event)

    go func(){
        if err := w.Watch(context.TODO(), eventCh); err != nil {
            panic(err)
        }
    }()

    for {
        select {
        case e := <-eventCh:
            switch e.(type) {
            case events.ServiceAdded:
                fmt.Println("service added!")
            case events.ServiceRemoved:
                fmt.Println("service removed!")
            }
        }
    }
}
```
