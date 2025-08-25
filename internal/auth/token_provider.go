package auth

import "context"

type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

type StaticTokenProvider struct {
	token string
}

func NewStaticTokenProvider(token string) *StaticTokenProvider {
	return &StaticTokenProvider{token: token}
}

func (p *StaticTokenProvider) GetToken(ctx context.Context) (string, error) {
	return p.token, nil
}

type RequestTokenProvider struct {
	tokenFunc func(ctx context.Context) (string, error)
}

func NewRequestTokenProvider(tokenFunc func(ctx context.Context) (string, error)) *RequestTokenProvider {
	return &RequestTokenProvider{tokenFunc: tokenFunc}
}

func (p *RequestTokenProvider) GetToken(ctx context.Context) (string, error) {
	return p.tokenFunc(ctx)
}
