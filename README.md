# goEventBus

event bus for Go's
. High performance
. asynchronous
. non blocking

# Getting started

## Install

Import package:

```
go get github.com/newbietao/goEventBus/event
```

```
import (
    "github.com/newbietao/goEventBus/event"
)
```

# Usage

```
p, err := pool.New("127.0.0.1:8080", pool.DefaultOptions)
if err != nil {
    log.Fatalf("failed to new pool: %v", err)
}
defer p.Close()

conn, err := p.Get()
if err != nil {
    log.Fatalf("failed to get conn: %v", err)
}
defer conn.Close()

// cc := conn.Value()
// client := pb.NewClient(conn.Value())
```
See the complete example: [https://github.com/shimingyah/pool/tree/master/example](https://github.com/shimingyah/pool/tree/master/example)
