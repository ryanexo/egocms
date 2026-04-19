package provider

import (
    `context`
    
    `cms/internal/domain/token/service`
    `cms/internal/middleware/authz`
)

type parser struct {
    srv *service.TokenService
}

var _ authz.TokenParser = (*parser)(nil)

func (s parser) Parse(ctx context.Context, token string) (authz.User, error) {
    u, err := s.srv.GetUserFromToken(ctx, token)
    if err != nil {
        return nil, err
    }
    return user{User: u}, nil
}

func NewTokenParser(srv *service.TokenService) authz.TokenParser {
    return parser{srv}
}
