package reflection

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"tests-coverage-tool/tool/config"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type GRPCReflectionClient struct {
	conn             *grpc.ClientConn
	reflectionClient *grpcreflect.Client
}

func NewGRPCReflectionClient(ctx context.Context, config config.Config) (*GRPCReflectionClient, error) {
	conn, err := grpc.NewClient(
		config.Service.GetURL(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	reflectionClient := grpcreflect.NewClientV1Alpha(ctx, grpc_reflection_v1alpha.NewServerReflectionClient(conn))

	return &GRPCReflectionClient{conn: conn, reflectionClient: reflectionClient}, nil
}

func (c *GRPCReflectionClient) GetServices() ([]string, error) {
	services, err := c.reflectionClient.ListServices()
	if err != nil {
		return nil, err
	}

	return lo.Filter(services, func(item string, _ int) bool {
		return !strings.HasPrefix(item, "grpc.reflection")
	}), nil
}

func (c *GRPCReflectionClient) GetServiceDescriptor(service string) (*desc.ServiceDescriptor, error) {
	return c.reflectionClient.ResolveService(service)
}

func (c *GRPCReflectionClient) GetServiceMethodsDescriptors(service string) ([]*desc.MethodDescriptor, error) {
	serviceDescriptor, err := c.GetServiceDescriptor(service)
	if err != nil {
		return nil, err
	}

	return serviceDescriptor.GetMethods(), nil
}
