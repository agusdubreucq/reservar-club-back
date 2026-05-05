package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwtlib.RegisteredClaims
}

type TokenManager struct {
	secret string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: secret}
}

func (tm *TokenManager) GenerateToken(userID, email string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.secret))
}

func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (interface{}, error) {
		return []byte(tm.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwtlib.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwtlib.ErrSignatureInvalid
	}

	return claims, nil
}
