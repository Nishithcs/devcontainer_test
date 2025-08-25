package api_clients

import (
	"clusterix-code/internal/auth"
	"clusterix-code/internal/data/dto"
	"context"
	"errors"
	"fmt"
	"strconv"
)

type AuthAPIClient struct {
	client *Client
}

func NewAuthAPIClient(baseURL string, tokenProvider auth.TokenProvider) *AuthAPIClient {
	client := NewClient(baseURL, func(c *Client) {
		c.SetTokenProvider(tokenProvider)
	})

	return &AuthAPIClient{client: client}
}

type GetUserResponse struct {
	Data     []dto.AuthUserDto `json:"data"`
	LastPage int               `json:"last_page"`
	Page     int               `json:"page"`
	Code     int               `json:"code"`
}

type GetApplicationsResponse struct {
	Data []dto.AuthApplicationDto `json:"data"`
}

type GetUserDataResponse struct {
	Data *dto.AuthUserDto `json:"data"`
}

type CreateAnonymousUserReq struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	OrganizationID uint64 `json:"organization_id"`
}

type CreateAnonymousUserResponse struct {
	Data *dto.AnonymousUserDto `json:"data"`
}

func (s *AuthAPIClient) GetUsers(ctx context.Context, page int) (*GetUserResponse, error) {
	var response GetUserResponse

	path := "/users?limit=100&with[]=profile"

	if page > 0 {
		path = path + "&page=" + strconv.Itoa(page)
	}
	err := s.client.Request(ctx, "GET", path, nil, &response, nil)

	return &response, err
}

func (s *AuthAPIClient) GetUser(ctx context.Context, id uint64) (*dto.AuthUserDto, error) {
	var response GetUserDataResponse

	// API request to get the user by ID
	err := s.client.Request(ctx, "GET", "/users/"+strconv.FormatUint(id, 10), nil, &response, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user from API: %w", err)
	}

	if response.Data == nil {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}

	return response.Data, nil
}

func (s *AuthAPIClient) GetApplications(ctx context.Context) (*GetApplicationsResponse, error) {
	var response GetApplicationsResponse

	err := s.client.Request(ctx, "GET", "/applications", nil, &response, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch applications from API: %w", err)
	}

	return &response, nil
}

func (s *AuthAPIClient) CreateAnonymousUser(ctx context.Context, email, name string) (*dto.AnonymousUserDto, error) {
	// we're using orgId of 0 for anonymous users
	const orgId uint64 = 0
	req := CreateAnonymousUserReq{
		Email:          email,
		Name:           name,
		OrganizationID: orgId,
	}
	var response CreateAnonymousUserResponse

	err := s.client.Request(ctx, "POST", "/anonymous/create", req, &response, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create anonymous user from auth: %w", err)
	}
	anonymousUser := response.Data
	if anonymousUser == nil {
		return nil, errors.New("no data received from auth or data couldn't be parsed")
	}
	anonymousUser.Email = email
	anonymousUser.Name = name
	anonymousUser.OrganizationID = orgId
	return response.Data, nil
}
