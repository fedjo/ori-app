package main

import (
    "os"
    "io"
    "fmt"
    "log"
    "time"
    "strconv"
    "context"

    "google.golang.org/grpc"

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
    log.Printf("%v arg", args)
    // Dial gRCP
    srvAddr := os.Getenv("SERVER_ADDRESS")
    conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Cannot dial server %v", err)
    }
    defer conn.Close()

    client := pb.NewOriServiceClient(conn)

    // Call related methods
    if args[0] == "sqrt" {
        _sq, err := strconv.ParseInt(args[1], 10, 32)
        if err != nil {
            log.Fatalf("Cannot parse sqrt arg %v", err)
        }
        // Implicit casting
        sq = _sq
        sqrtVal := &pb.Value{V: sq}
        log.Printf("Hello sqrt %v\n", sq)
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        stream, err := client.Sqrt(ctx, sqrtVal)
        if err != nil {
            log.Fatalf("%v.Sqrt(_) = _, %v", client, err)
        }
        for {
            feature, err := stream.Recv()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Fatalf("%v.Sqrt(_) = _, %v", client, err)
            }
            log.Println(feature)
        }


    } else if args[0] == "gcd" {
        // Get argument values
        _a, err := strconv.ParseInt(args[1], 10, 32)
        _b, err := strconv.ParseInt(args[2], 10, 32)
        // Casting to int64
        a = _a
        b = _b
        pointCount := 2
        var points []*pb.Point
        for i := 0; i < pointCount; i++ {
            points = append(points, &pb.Point{X: a+int64(i), Y: b+int64(i)})
        }
        log.Printf("Created %v points", len(points))
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        stream, err := client.Gcd(ctx)
        if err != nil {
            log.Fatalf("%v.Gcd(_) = _, %v", client, err)
        }
        for _, point := range points {
            if err := stream.Send(point); err != nil {
                log.Fatalf("%v.Send(%v) = %v", stream, point, err)
            }
        }
        reply, err := stream.CloseAndRecv()
        if err != nil {
            log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
        }
        log.Printf("Route summary: %v", reply)

    } else if args[0] == "sum" {
        _a, err := strconv.ParseInt(args[1], 10, 32)
        if err != nil {
            log.Fatalf("Cannot parse first arg %v", err)
        }
        a = _a
        _b, err := strconv.ParseInt(args[2], 10, 32)
        if err != nil {
            log.Fatalf("Cannot parse second arg %v", err)
        }
        b = _b
        fmt.Printf("Hello sum %v %v\n", a, b)
        // Call Sum method
        req := &pb.Point{X: a, Y: b}
        log.Println(req)
        res, err := client.Sum(context.Background(), req)
        if err != nil {
            log.Fatalf("gRPC request failed: %v", err)
        }
        log.Printf("Receive: %d %s\n", res.Ret, res.Msg)
    } else {
        log.Printf("Wrong command line arguments")
        os.Exit(-1)
    }
}
