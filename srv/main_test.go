package main_test

import (
    "os"
    "context"
    "testing"

    . "github.com/onsi/gomega"
	kitlog "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/log/level"

	srv "github.com/fedjo/ori-app/srv"
	pb "github.com/fedjo/ori-app/pb"
)

var (
    appName = "Server unit test"

    logger            kitlog.Logger
)


func TestSum(t *testing.T) {


    // Set up logger
	logger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "app", appName)
    logger = level.NewFilter(logger, level.AllowAll())

    level.Info(logger).Log("msg", "Starting server testcase")

    testCases := []struct {
        name    string
        req     *pb.Point
        msg     string
        expectedErr bool
    }{
        {
            name: "OK",
            req: &pb.Point{
                X: 42,
                Y: 42,
            },
            msg: "This is the sum",
            expectedErr: false,
        },
        {
            name: "Nil request",
            req: nil,
            expectedErr: true,
        },
    }


    for _, tc := range(testCases) {
        testCase := tc
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()
            g := NewGomegaWithT(t)

            ctx := context.Background()

            // call srv
            grpcSrv := srv.NewServer("Test server", logger)
            level.Info(logger).Log("msg", "testcase", "Request", testCase.req)
            resp, err := grpcSrv.Sum(ctx, testCase.req)

            level.Info(logger).Log("msg", "grcp repsonse", "Response", resp)
            t.Log("Got : ", resp)

            if testCase.expectedErr {
                g.Expect(resp).ToNot(BeNil(), "Result should be nil")
				g.Expect(err).ToNot(BeNil(), "Result should be nil")
            } else {
                g.Expect(resp.Ret).To(Equal(testCase.req.X + testCase.req.Y))
            }
        })
    }
}
