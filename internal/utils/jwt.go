package utils

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"strings"

// 	"github.com/golang-jwt/jwt/v4"

// 	"golang.org/x/exp/slices"
// )

// func DecodeJWT(request *http.Request) (*http.Request, error) {
// 	header := request.Header.Get("authorization")
// 	parts := strings.Split(header, " ")

// 	if len(parts) != 2 {
// 		return nil, ErrInvalidAuthorizationHeader
// 	}

// 	if !slices.Contains(AllowedTokenTypes, strings.ToLower(parts[0])) {
// 		return nil, ErrInvalidTokenType
// 	}

// 	token, err := jwt.ParseWithClaims(
// 		parts[1],
// 		&validators.MerchantToken{},
// 		keyFunc,
// 		jwt.WithValidMethods([]string{"HS256"}),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*validators.MerchantToken); ok && token.Valid {
// 		request = contextSetAuthenticatedMerchant(request, *claims)
// 		return request, nil
// 	} else {
// 		return nil, ErrInvalidToken
// 	}
// }

// var keyFunc = func(t *jwt.Token) (interface{}, error) {
// 	key := secret
// 	return key, nil
// }

// func contextSetAuthenticatedMerchant(r *http.Request, user validators.MerchantToken) *http.Request {
// 	ctx := context.WithValue(r.Context(), AuthenticatedMerchantContextKey, user)
// 	return r.WithContext(ctx)
// }

// var (
// 	ErrInvalidAuthorizationHeader = errors.New("invalid or missing authorization header")
// 	ErrInvalidTokenType           = errors.New("invalid token type")
// 	ErrInvalidToken               = errors.New("invalid token")
// )

// type contextKey string

// const (
// 	AuthenticatedMerchantContextKey = contextKey("authenticatedMerchant")
// )

// var AllowedTokenTypes = []string{"upay"}

// var secret []byte

// func SetMerchantSecret(key string) {
// 	// fmt.Println(key)
// 	secret = []byte(key)
// }
