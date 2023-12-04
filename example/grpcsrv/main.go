package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/nico612/voyage/api/gen/voyage/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"strings"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: in.Name + "world",
	}, nil

}

func main() {

	listenGrpcAndHttpSamePort()
}

func listenGrpcAndHttpNotSamePort() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln(err)
	}

	// 创建一个 grpc server 对象
	s := grpc.NewServer()
	// 注册 Greeter 服务 到 server
	pb.RegisterGreeterServer(s, &server{})

	// 启动 grpc server
	log.Println("Serving gRPC on 0.0.0.0:9090")

	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// 创建一个连接到我们刚刚启动的 gRPC 服务器的客户端连接
	// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:9090",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// 注册Greeter
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	// 8090端口提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func listenGrpcAndHttpSamePort() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// 创建一个 grpc 服务
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	// gRPC-Gateway mux
	gwmux := runtime.NewServeMux()
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterGreeterHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:9090", dops)
	if err != nil {
		log.Fatalln("failed to register gwmux:", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// 定义HTTP server配置
	gwServer := &http.Server{
		Addr:    "127.0.0.1:9090",
		Handler: grpcHandlerFunc(s, mux), // 请求的统一入口
	}
	log.Println("Serving on http://127.0.0.1:9090")
	log.Fatalln(gwServer.Serve(lis)) // 启动HTTP服务
}

// grpcHandlerFunc 将gRPC请求和HTTP请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
