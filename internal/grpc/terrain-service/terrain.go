package terrain_service

import (
	"context"

	"github.com/TwiLightDM/diploma-gateway/package/terrainpb"
)

func (c *TerrainClient) SubmitJob(ctx context.Context, userID string, image []byte, scaleZ float32, yUp bool, textureMode string) (string, error) {
	resp, err := c.terrain.SubmitJob(ctx, &terrainpb.SubmitJobRequest{
		UserId:      userID,
		Image:       image,
		ScaleZ:      scaleZ,
		YUp:         yUp,
		TextureMode: textureMode,
	})
	if err != nil {
		return "", err
	}
	return resp.JobId, nil
}

func (c *TerrainClient) GetJobStatus(ctx context.Context, jobID string) (*terrainpb.GetJobStatusResponse, error) {
	return c.terrain.GetJobStatus(ctx, &terrainpb.GetJobStatusRequest{JobId: jobID})
}

func (c *TerrainClient) GetUserJobs(ctx context.Context, userID string) (*terrainpb.GetUserJobsResponse, error) {
	return c.terrain.GetUserJobs(ctx, &terrainpb.GetUserJobsRequest{UserId: userID})
}

func (c *TerrainClient) GetLimit(ctx context.Context, userID string) (*terrainpb.GetLimitResponse, error) {
	return c.terrain.GetLimit(ctx, &terrainpb.GetLimitRequest{UserId: userID})
}

func (c *TerrainClient) SetLimit(ctx context.Context, userID string, dailyLimit int32) (*terrainpb.SetLimitResponse, error) {
	return c.terrain.SetLimit(ctx, &terrainpb.SetLimitRequest{UserId: userID, DailyLimit: dailyLimit})
}
