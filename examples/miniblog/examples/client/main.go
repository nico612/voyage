package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
	pb "github.com/nico612/go-project/examples/miniblog/pkg/proto/miniblog/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	addr  = flag.String("addr", "localhost:9090", "The address to connect to.")
	limit = flag.Int64("limit", 10, "Limit to list users.")
)

func main() {
	flag.Parse()
	// 建立与服务器的连接
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalw("Did not connect", "err", err)
	}
	defer conn.Close()
	c := pb.NewMiniBlogClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 请求 ListUser 接口
	offset := int64(0)
	r, err := c.ListUser(ctx, &pb.ListUserRequest{Offset: &offset, Limit: limit})
	if err != nil {
		log.Fatalw("could not greet: %v", err)
	}

	// 打印请求结果
	fmt.Println("TotalCount:", r.TotalCount)
	for _, u := range r.Users {
		d, _ := json.Marshal(u)
		fmt.Println(string(d))
	}
}