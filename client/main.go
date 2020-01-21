package main

import (
    "os"
    "fmt"
    "log"
    "strconv"
    "context"
    // "flag"
    // "io"
    // "math/rand"
    // "time"

    "google.golang.org/grpc"
    // "google.golang.org/grpc/credentials"

    "github.com/fedjo/ori-app/pb"
)


var (
    sq  int64
    a   int64
    b   int64
)


func main() {

    // Parse arguments
    args := os.Args[1:]
    if args[0] == "sqrt" {
        _sq, err := strconv.ParseInt(args[1], 10, 32)
        if err != nil {
            log.Fatalf("Cannot parse sqrt arg %v", err)
        }
        sq = _sq
        fmt.Printf("Hello sqrt %v\n", sq)
    } else {
        _a, err := strconv.ParseInt(args[0], 10, 32)
        if err != nil {
            log.Fatalf("Cannot parse first arg %v", err)
        }
        a = _a
        _b, err := strconv.ParseInt(args[1], 10, 32)
        if err != nil {
            log.Fatalf("Cannot parse second arg %v", err)
        }
        b = _b
        fmt.Printf("Hello sum %v %v\n", a, b)
    }


    // Dial gRCP
    srvAddr := os.Getenv("SERVER_ADDRESS")
    conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Cannot dial server %v", err)
    }
    defer conn.Close()

    client := pb.NewOriServiceClient(conn)

    // Call Sum method
    req := &pb.Point{X: a, Y: b}
    log.Println(req)
    res, err := client.Sum(context.Background(), req)
    if err != nil {
        log.Fatalf("gRPC request failed: %v", err)
    }
    log.Printf("Receive: %d %s\n", res.Ret, res.Msg)
}
