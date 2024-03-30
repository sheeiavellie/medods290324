package middlewares

import (
	"context"
	"fmt"

	"github.com/sheeiavellie/medods290324/data"
)

var (
	contextKeyTokens = contextKey("accountRequest")
)

type contextKey string

func (c contextKey) String() string {
	return "middlewares context key " + string(c)
}

func GetTokensKey(ctx context.Context) (*data.TokensResponse, error) {
	if v := ctx.Value(contextKeyTokens); v != nil {
		if v, ok := v.(*data.TokensResponse); ok {
			return v, nil
		}
		err := fmt.Errorf("error getting key %s: wrong type", contextKeyTokens.String())
		return nil, err
	}
	err := fmt.Errorf("error getting key %s: no value", contextKeyTokens.String())
	return nil, err
}
