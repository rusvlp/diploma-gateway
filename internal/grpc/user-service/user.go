package user_service

import (
	"context"

	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
)

func (c *UserClient) Login(ctx context.Context, email, password string) (*userservicepb.LoginResponse, error) {
	return c.user.Login(ctx, &userservicepb.LoginRequest{
		Email:    email,
		Password: password,
	})
}

func (c *UserClient) SignUp(ctx context.Context, fullName, role, email, password string) (*userservicepb.SignUpResponse, error) {
	return c.user.SignUp(ctx, &userservicepb.SignUpRequest{
		FullName: fullName,
		Role:     role,
		Email:    email,
		Password: password,
	})
}

func (c *UserClient) Refresh(ctx context.Context, id, role string) (*userservicepb.RefreshResponse, error) {
	return c.user.Refresh(ctx, &userservicepb.RefreshRequest{
		Id:   id,
		Role: role,
	})
}

func (c *UserClient) ReadUser(ctx context.Context, id string) (*userservicepb.ReadUserResponse, error) {
	return c.user.ReadUser(ctx, &userservicepb.ReadUserRequest{
		Id: id,
	})
}

func (c *UserClient) ReadUsers(ctx context.Context) (*userservicepb.ReadUsersResponse, error) {
	return c.user.ReadUsers(ctx, &userservicepb.ReadUsersRequest{})
}

func (c *UserClient) UpdateUser(ctx context.Context, id, fullName, email string) (*userservicepb.UpdateUserResponse, error) {
	return c.user.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		Id:       id,
		FullName: fullName,
		Email:    email,
	})
}

func (c *UserClient) UpdateUserRole(ctx context.Context, id, role string) (*userservicepb.UpdateUserRoleResponse, error) {
	return c.user.UpdateUserRole(ctx, &userservicepb.UpdateUserRoleRequest{
		Id:   id,
		Role: role,
	})
}

func (c *UserClient) ChangePassword(ctx context.Context, id, password string) (*userservicepb.ChangePasswordResponse, error) {
	return c.user.ChangePassword(ctx, &userservicepb.ChangePasswordRequest{
		Id:       id,
		Password: password,
	})
}
