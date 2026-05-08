package terrain_service

import (
	"log"

	"github.com/TwiLightDM/diploma-gateway/package/terrainpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TerrainClient struct {
	terrain terrainpb.TerrainServiceClient
	conn    *grpc.ClientConn
}

func NewTerrainClient(address string) *TerrainClient {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to terrain service: %v", err)
	}
	return &TerrainClient{
		terrain: terrainpb.NewTerrainServiceClient(conn),
		conn:    conn,
	}
}

func (c *TerrainClient) Close() error {
	return c.conn.Close()
}
