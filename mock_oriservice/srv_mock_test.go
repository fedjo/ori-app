package mock_pb

import (
    "fmt"
    "testing"
    "time"
    "context"
    "gotest.tools/assert"

   	"github.com/golang/protobuf/proto"
    "github.com/golang/mock/gomock"
    srvmock "github.com/fedjo/ori-app/mock_oriservice"
    pb "github.com/fedjo/ori-app/pb"
)


var point = &pb.Point{
    X: 42,
    Y: 43,
}

// Define matcher for Point
// rpcMsg implements the gomock.Matcher interface
type rpcPoint struct {
	msg proto.Message
}

func (r *rpcPoint) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcPoint) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

// Test sum functionality for a given point
func TestSum(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockOriServiceClient := srvmock.NewMockOriServiceClient(ctrl)
    mockOriServiceClient.EXPECT().Sum(
        gomock.Any(),
        &rpcPoint{msg: point},
    ).Return(&pb.Ret{Ret: (point.X + point.Y) }, nil)
    testSum(t, mockOriServiceClient)
}

func testSum(t *testing.T, client pb.OriServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    r, err := client.Sum(ctx, point)
    if err != nil {
        t.Errorf("mocking failed")
    }
    t.Log("Reply : ", r.Ret)
    assert.Equal(t, int64(85), r.Ret)
}

// Test GCD functionality
func TestGcd(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Create mock for the client interface.
    stream := srvmock.NewMockOriService_GcdClient(ctrl)
    // set expectation on sending.
    stream.EXPECT().Send(
        gomock.Any(),
    ).Return(nil)
    // Set expectation on receiving.
    stream.EXPECT().CloseAndRecv().Return(&pb.RetSummary{Ret: []int64{1,},}, nil)
    stream.EXPECT().CloseSend().Return(nil)

    // test
    if err := stream.Send(point); err != nil {
        t.Fatalf("Test failed: %v", err)
    }
    if err := stream.CloseSend(); err != nil {
        t.Fatalf("Test failed: %v", err)
    }
    got, err := stream.CloseAndRecv()
    if err != nil {
        t.Fatalf("Test failed: %v", err)
    }
    assert.Equal(t, len(got.Ret), 1)
}
