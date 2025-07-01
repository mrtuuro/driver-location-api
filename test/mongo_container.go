package test

import (
	"context"
	"fmt"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartMongo(ctx context.Context) (uri string, terminate func()) {
	req := tc.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}
	c, _ := tc.GenericContainer(ctx, tc.GenericContainerRequest{Started: true, ContainerRequest: req})
	host, _ := c.Host(ctx)
	port, _ := c.MappedPort(ctx, "27017")
	return fmt.Sprintf("mongodb://%s:%s", host, port.Port()), func() { _ = c.Terminate(ctx) }
}
