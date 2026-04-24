package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Address  string
	Insecure bool
	Timeout  time.Duration
}

func Dial(ctx context.Context, cfg Config, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			loggingInterceptor,
		),
	}

	if cfg.Insecure {
		defaultOpts = append(defaultOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	dialCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, cfg.Address, append(defaultOpts, opts...)...)
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}

	return conn, nil
}

func loggingInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("[grpc] method=%s duration=%s error=%v", method, time.Since(start), err)
	return err
}
