package main

import (
	"context"
	"encoding/json"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type adminService struct {
}

func (as *adminService) Logging(n *Nothing, logServ Admin_LoggingServer) error {
	return nil
}

func (as *adminService) Statistics(interval *StatInterval, statServ Admin_StatisticsServer) error {
	return nil
}

func NewAdminService() AdminServer {
	service := &adminService{}
	return service
}

type bizService struct {
}

func (bz *bizService) Check(ctx context.Context, n *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (bz *bizService) Add(ctx context.Context, n *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (bz *bizService) Test(ctx context.Context, n *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}
func NewBizService() BizServer {
	service := &bizService{}
	return service
}

func authInterceptor(acl map[string][]string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// start := time.Now()

		md, _ := metadata.FromIncomingContext(ctx)

		consumers := md.Get("consumer")
		if len(consumers) == 0 {
			return nil, grpc.Errorf(codes.Unauthenticated, "Unauthenticated")
		}

		cons := consumers[0]
		consMethods, ok := acl[cons]
		if !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "Unauthenticated")
		}

		allowed := false
		for _, m := range consMethods {
			// fmt.Println("allowed calc: ", info.FullMethod, m)

			if strings.Contains(m, "/*") {
				if strings.Contains(info.FullMethod, strings.Replace(m, "/*", "", -1)) {
					allowed = true
					break
				}
			} else if info.FullMethod == m {
				allowed = true
				break
			}
		}

		if !allowed {
			return nil, grpc.Errorf(codes.Unauthenticated, "Unauthenticated")
		}

		reply, err := handler(ctx, req)

	// 	fmt.Printf(`
	// info.FullMethod = %v
	// req = %#v
	// reply = %#v
	// time = %v
	// md = %v
	// err = %v
	// `, info.FullMethod, req, reply, time.Since(start), md, err)

		return reply, err
	}
}

func StartMyMicroservice(ctx context.Context, addr string, ACLData string) error {
	acl := make(map[string][]string)
	if err := json.Unmarshal([]byte(ACLData), &acl); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor(acl)),
	)
	RegisterAdminServer(server, NewAdminService())
	RegisterBizServer(server, NewBizService())

	go server.Serve(lis)

	go func() {
		<-ctx.Done()
		err := lis.Close()
		if err != nil {
			panic(err.Error())
		}
	}()

	return nil
}
