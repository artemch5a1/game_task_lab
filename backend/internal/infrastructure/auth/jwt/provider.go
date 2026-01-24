package jwt

import (
	"context"
	"errors"
	"time"

	specifictype "example/web-service-gin/internal/domain/specific_type"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Role string `json:"role"`
	jwtlib.RegisteredClaims
}

type Provider struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewProvider(secret, issuer string, ttl time.Duration) *Provider {
	return &Provider{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (p *Provider) Issue(ctx context.Context, userID uuid.UUID, role specifictype.UserRole) (string, error) {
	_ = ctx
	if userID == uuid.Nil {
		return "", errors.New("userID is empty")
	}
	now := time.Now().UTC()

	claims := Claims{
		Role: string(role),
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    p.issuer,
			Subject:   userID.String(),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(p.ttl)),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *Provider) Parse(tokenString string) (*Claims, error) {
	parser := jwtlib.NewParser(jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Name}))

	var claims Claims
	_, err := parser.ParseWithClaims(tokenString, &claims, func(token *jwtlib.Token) (any, error) {
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

