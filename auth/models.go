package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtPayload struct {
	IssuedAt  time.Time `json:"iat,omitempty"`
	ExpiredAt time.Time `json:"exp,omitempty"`
	NotBefore time.Time `json:"nbf,omitempty"`
}

////
//func NewJwtPayload() jwt.Claims {
//	return &JwtPayload{
//		UserPayload: UserPayload{
//			UserID:    "",
//			Role:      "GUEST",
//			UserName:  "",
//			SessionId: "",
//			Rank:      11,
//		},
//		ExpiredAt: time.Now(),
//		NotBefore: time.Now(),
//		IssuedAt:  time.Now(),
//	}
//}
//func NewGuestPayload() *jwt.Token {
//	return &jwt.Token{
//		Claims: NewJwtPayload(),
//	}
//}

func (t JwtPayload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(t.ExpiredAt), nil
}

func (t JwtPayload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(t.IssuedAt), nil
}

func (t JwtPayload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(t.NotBefore), nil
}

func (t JwtPayload) GetIssuer() (string, error) {
	return "", nil
}

func (t JwtPayload) GetSubject() (string, error) {
	return "", nil
}

func (t JwtPayload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}
