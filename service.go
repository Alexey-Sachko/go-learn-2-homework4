package main

import (
	"context"
	"encoding/json"
	"net"

	"google.golang.org/grpc"
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

func StartMyMicroservice(ctx context.Context, addr string, ACLData string) error {
	acl := make(map[string][]string)
	if err := json.Unmarshal([]byte(ACLData), &acl); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterAdminServer(server, NewAdminService())

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
