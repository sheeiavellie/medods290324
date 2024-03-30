package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sheeiavellie/medods290324/data"
)

const (
	jwtTTL     = time.Minute * 15
	refreshTTL = time.Hour * 24 * 30
)

type TokenService struct {
	secret []byte
}

func NewTokenService(signSecret []byte) *TokenService {
	return &TokenService{
		secret: signSecret,
	}
}

func (ts *TokenService) IssueTokens(
	user *data.User,
) (*data.TokensResponse, error) {
	tokenJWT := ts.generateJWT(user.ID.String())
	tokenJWTSigned, err := ts.signJWT(tokenJWT)
	if err != nil {
		return nil, err
	}

	tokenRefresh, err := ts.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &data.TokensResponse{
		AccessToken:  tokenJWTSigned,
		RefreshToken: tokenRefresh,
	}, nil
}

func (ts *TokenService) generateJWT(subject string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(jwtTTL)},
		Subject:   subject,
	})

	return token
}

func (ts *TokenService) signJWT(token *jwt.Token) (string, error) {
	tokenStr, err := token.SignedString(ts.secret)
	if err != nil {
		return "", nil
	}

	return tokenStr, nil
}

func (ts *TokenService) generateRefreshToken() (string, error) {
	refreshToken := uuid.New().String()

	return refreshToken, nil
}
