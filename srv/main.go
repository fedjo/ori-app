package main

import (
	"io"
	"math"
	"net"
	"os"
    "os/signal"
    "syscall"
    "time"
    "errors"

	kitlog "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/log/level"
	context "golang.org/x/net/context"
    "golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
    "google.golang.org/grpc/keepalive"

	pb "github.com/fedjo/ori-app/pb"
)

const appName = "ori-grpc"

// Define global variables like logger, grpcServer
var (
    version = os.Getenv("VERSION")
	srvPort = os.Getenv("BIND_PORT")

    logger            kitlog.Logger
    grpcServer        *grpc.Server
)

// Define the server struct which methods are to be exposed if gRPC
type Server struct {
	appName string
	logger  kitlog.Logger
}

// Constructor method
func NewServer(appName string, logger kitlog.Logger) *Server {
    level.Info(logger).Log("msg", "Creating server")
	return &Server{
		appName: appName,
		logger:  logger,
	}
}

// Calculating sum
func (s *Server) Sum(ctx context.Context, p *pb.Point) (*pb.Ret, error) {
    level.Info(s.logger).Log("msg", "Received", "point", p)

    if p == nil {
        level.Error(s.logger).Log("msg", "Nil Point provided")
        return &pb.Ret{}, errors.New("Nil Point provided")
    }
	x, y := p.X, p.Y
	sum := x + y
    level.Info(s.logger).Log("msg", "Sum calculated: ", "value", sum)

	return &pb.Ret{Ret: sum, Msg: "Success"}, nil
}

// Calculating gcd of points
func (s *Server) Gcd(stream pb.OriService_GcdServer) error {

    var pointCount    int64
    var gcdList     []int64
    for {
        point, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&pb.RetSummary{
                Ret: gcdList,
                TotalPoints: pointCount,
                Msg: "Success",
            })
        }
        if err != nil {
            return err
        }
        if point == nil {
            level.Error(s.logger).Log("msg", "Nil Point provided")
            return errors.New("Nil Point provided")
        }
        pointCount++
        // Calculate gcd
        x, y := point.X, point.Y
        for y != 0 {
            x, y = y, x%y
        }
        gcdList = append(gcdList, x)
    }
}

// Calculating the sqrt
func (s *Server) Sqrt(v *pb.Value, stream pb.OriService_SqrtServer) error {

	sqrt := math.Sqrt(math.Abs(float64(v.V)))
	res := &pb.Ret{Ret: int64(sqrt), Msg: "Success"}
    for i := 1;  i<=3; i++ {
        level.Debug(s.logger).Log("msg", "Sending to stream")
        if err := stream.Send(res); err != nil {
            return err
        }
	}
	return nil
}

// Main func
func main() {

    // Set up logger
	logger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "app", appName)
    logger = level.NewFilter(logger, level.AllowAll())

    level.Info(logger).Log("msg", "Starting app", "version", version)

    ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

    // catch termination
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

    // Create errgroup with context
    g, ctx := errgroup.WithContext(ctx)

    // Initialize server struct
	s := NewServer(appName, logger)

    // Create go routine for grpc server
	g.Go(func () error {
        level.Info(logger).Log("msg", "Starting gRPC server on", "port", srvPort)
	    lis, err := net.Listen("tcp", ":"+srvPort)
	    if err != nil {
            level.Error(logger).Log("msg", "gRPC server: failed to listen", "error", err)
		    os.Exit(-1)
	    }

	    // Create server instance
	    grpcServer = grpc.NewServer(
		    grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
	    )
	    pb.RegisterOriServiceServer(grpcServer, s)

        level.Info(logger).Log("msg", "gRPC server serving", "localhost:", srvPort)
	    // Register reflection service for client info
	    reflection.Register(grpcServer)
	    if err := grpcServer.Serve(lis); err != nil {
            level.Info(logger).Log("msg", "Failed to server gRPC server")
	    }

        return grpcServer.Serve(lis)
    })


    // Catch SIGTERM
    select {
	case <-interrupt:
        level.Info(logger).Log("msg", "Interrupt case")
		break
	case <-ctx.Done():
        level.Info(logger).Log("msg", "Context Done")
		break
	}

	level.Warn(logger).Log("msg", "received shutdown signal")
    // Shutdown gRPC server
    if grpcServer != nil {
		level.Error(logger).Log("msg", "gRPC Server bye bye...")
		grpcServer.GracefulStop()
	}
	cancel()

    // Shutdown context
    _, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()


    // Wait and close
    err := g.Wait()
	if err != nil {
		level.Error(logger).Log("msg", "server returning an error", "error", err)
		os.Exit(-1)
	}
    level.Info(logger).Log("msg", "Waiting...")
}
