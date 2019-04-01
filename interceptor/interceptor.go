package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

func WithServerUnaryInterceptor() grpc.ServerOption {
	interceptor := &LogInterceptor{metrics: make(map[string]int)}
	go interceptor.metricReporter()

	return grpc.UnaryInterceptor(interceptor.logInterceptor)
}


type LogInterceptor struct {
	metrics 		map[string]int
}


func (l *LogInterceptor) metricReporter() {
	for true {
		time.Sleep(1 * time.Minute)
		for meth, count := range l.metrics {
			fmt.Printf(fmt.Sprintf("Method: %s invocations: %d ", meth, count))
		}
	}
}

func (l *LogInterceptor) logInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	fmt.Println("Executing method " + info.FullMethod)
	count, found := l.metrics[info.FullMethod]
	if !found {
		l.metrics[info.FullMethod] = 1
	} else {
		l.metrics[info.FullMethod] = count + 1
	}

	// Calls the handler
	h, err := handler(ctx, req)

	// Logging with grpclog (grpclog.LoggerV2)
	if err != nil {
		fmt.Println(fmt.Sprintf("Request - Method: %s\nDuration:%s\nError:%v\n",
			info.FullMethod,
			time.Since(start),
			err))
	} else {
		fmt.Println(fmt.Sprintf("Request - Method: %s\nDuration:%s\n",
			info.FullMethod,
			time.Since(start)))
	}
	return h, err
}