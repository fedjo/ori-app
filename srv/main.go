package main

import (
    "os"
    "log"
    "net"
    "io"
    "math"
    // Change this for your own project
    "github.com/fedjo/ori-app/pb"
    context "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)


type server struct{}



func (s *server) Sum(ctx context.Context, p *pb.Point) (*pb.Ret, error) {
    log.Println(p)
    x, y := p.X, p.Y
    sum := x + y
    log.Printf("Sum calculated: %d\n", sum)
    return &pb.Ret{Ret: sum, Msg: "Success"}, nil
}

func (s *server) Gcd(stream pb.OriService_GcdServer) error {

    point, err := stream.Recv()
    if err == io.EOF {
        x, y := point.X, point.Y
        for y != 0 {
            x, y = y, x%y
        }
        return stream.SendAndClose(&pb.Ret{
            Ret: x,
            Msg: "Success",
        })
    }
    if err != nil {
        return err
    }
    return nil
}

func (s *server) Sqrt(v *pb.Value, stream pb.OriService_SqrtServer) error {

    sqrt := math.Sqrt(float64(v.V))
    res := &pb.Ret{Ret: int64(sqrt), Msg: "Success"}
    if err := stream.Send(res); err != nil {
        return err
    }
    return nil
}


func main() {

    srvPort := os.Getenv("BIND_PORT")
    lis, err := net.Listen("tcp", ":" + srvPort)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // Create server instance
    s := grpc.NewServer()
    pb.RegisterOriServiceServer(s, &server{})

    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
