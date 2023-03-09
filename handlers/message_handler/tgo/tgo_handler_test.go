// https: //github.com/p4gefau1t/trojan-go/blob/master/api/service/client.go

package tgo

import (
	"testing"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"run1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, _ := grpc.Dial("127.0.0.1:10000", grpc.WithInsecure())
			service := NewTrojanServerServiceClient(conn)
			c, err := service.ListUsers(context.TODO(), &ListUsersRequest{})
			if err != nil {
				println(err)
				return
			}

			var status *UserStatus
			for response, err := c.Recv(); err == nil; response, err = c.Recv() {
				print(response.Status)
				status = response.GetStatus()
			}

			if status != nil {
				client, _ := service.SetUsers(context.TODO())
				status.TrafficTotal.DownloadTraffic += 1 * 1024 * 1024
				err := client.Send(&SetUsersRequest{Status: status, Operation: SetUsersRequest_Modify})
				if err != nil {
					println(err)
				}
			}
		})
	}
}
