package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type CustomClaims struct {
	// TokenType string `json:"token_type"`
	UserId int32 `json:"user_id"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	if c.UserId <= 0 {
		return jwtmiddleware.ErrJWTInvalid
	}

	return nil
}

func getSettings() (secret []byte, cookieName string, issuer string, audience string) {
	cookieName = "access"
	issuer = "localhost"
	audience = "toiler"

	if sec := os.Getenv("JWT_SIGNING_KEY"); sec != "" {
		secret = []byte(sec)
	} else {
		panic("JWT_SIGNING_KEY variable is empty.")
	}

	if v := os.Getenv("JWT_COOKIE_NAME"); v != "" {
		cookieName = v
	}

	if v := os.Getenv("JWT_ISSUER"); v != "" {
		issuer = v
	}

	if v := os.Getenv("JWT_AUDIENCE"); v != "" {
		audience = v
	}

	return
}

func AuthMiddleware() func(http.Handler) http.Handler {
	secret, cookieName, issuer, audience := getSettings()

	keyFunc := func(ctx context.Context) (interface{}, error) {
		// Our token must be signed using this data.
		return secret, nil
	}
	// We want this struct to be filled in with
	// our custom claims from the token.
	customClaims := func() validator.CustomClaims {
		return &CustomClaims{}
	}

	tokenExtractor := jwtmiddleware.CookieTokenExtractor(cookieName)

	// Set up the validator.
	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		issuer,
		[]string{audience},
		validator.WithCustomClaims(customClaims),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	// Set up the middleware.
	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithTokenExtractor(tokenExtractor),
	)

	return middleware.CheckJWT
}

var jwtContextKey = jwtmiddleware.ContextKey{}

func GetUserId(ctx context.Context) (int32, error) {
	a := ctx.Value(jwtContextKey)
	claims, ok := a.(*validator.ValidatedClaims)
	if !ok {
		return 0, errors.New("failed to get validated claims")
	}

	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return 0, errors.New("could not cast custom claims to specific type")
	}

	return customClaims.UserId, nil
}

func GenToken(userID int32) (string, error) {
	secret, _, issuer, audience := getSettings()

	key := jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       secret,
	}

	signer, err := jose.NewSigner(key, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", fmt.Errorf("could not build signer: %w", err)
	}

	claims := jwt.Claims{
		Issuer:   issuer,
		Audience: jwt.Audience{audience},
	}
	customClaims := CustomClaims{
		UserId: userID,
	}

	token, err := jwt.Signed(signer).Claims(claims).Claims(customClaims).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("could not build token: %w", err)
	}

	return token, nil
}
