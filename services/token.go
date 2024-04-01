package services

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sheeiavellie/medods290324/data"
)

const (
	jwtTTL = time.Minute * 15
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
	userID string,
	sessionID string,
) (*data.TokensResponse, error) {
	issuingTime := time.Now()

	tokenJWT := ts.generateJWT(userID, issuingTime)
	tokenJWTSigned, err := ts.signJWT(tokenJWT)
	if err != nil {
		return nil, err
	}

	tokenRefresh, err := ts.generateRefreshToken(sessionID, userID)
	if err != nil {
		return nil, err
	}

	return &data.TokensResponse{
		AccessToken:  tokenJWTSigned,
		RefreshToken: tokenRefresh,
	}, nil
}

func (ts *TokenService) DecodeRefreshToken(
	refreshToken string,
) (*data.RefreshTokenClaims, error) {
	refreshStr, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return nil, err
	}

	var claims data.RefreshTokenClaims
	err = json.Unmarshal(refreshStr, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func (ts *TokenService) generateJWT(
	userID string,
	issuingTime time.Time,
) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data.JWTTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: issuingTime.Add(jwtTTL)},
			IssuedAt:  &jwt.NumericDate{Time: issuingTime},
		},
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

func (ts *TokenService) generateRefreshToken(
	sessionID string,
	userID string,
) (string, error) {
	claims := data.RefreshTokenClaims{
		SessionID: sessionID,
		UserID:    userID,
	}

	refreshStr, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(refreshStr)

	return refreshToken, nil
}
