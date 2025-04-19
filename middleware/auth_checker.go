package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Suryarpan/user-auth-jwt/database"
	"github.com/Suryarpan/user-auth-jwt/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ctxKeyUserData string

const (
	ctxUserDataKey  ctxKeyUserData = "USER_AUTH_USER_DATA"
	userAuthHeader  string         = "Authorization"
	userTokenPrefix string         = "Bearer "
	TokenIssuer     string         = "chat-api-auth"
	lenPrefix       int            = len(userTokenPrefix)
	// messages
	notAuthed     string = "token not available"
	wrongHeader   string = "incorrect header structure"
	invalidToken  string = "token is invalid"
	invalidUser   string = "user is unavailable or invalid"
	serverProblem string = "unable to login"
)

var (
	RegularAudience []string = []string{"user"}
	AdminAudience   []string = []string{"user", "admin"}
)

type tokenData struct {
	jwt.RegisteredClaims
	UserId pgtype.UUID `json:"uid"`
}

func tokenToUser(s string, secret []byte) (*tokenData, error) {
	token, err := jwt.ParseWithClaims(
		s,
		&tokenData{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		},
		jwt.WithAudience(RegularAudience[0]),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithLeeway(time.Second*10),
		jwt.WithJSONNumber(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS384.Name}),
	)
	if err != nil {
		slog.Warn("invalid token encountered", "error", err)
		return nil, err
	}
	if claims, ok := token.Claims.(*tokenData); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("unknown claims type")
}

func UserToToken(u database.User) (string, error) {
	config := utils.NewConf()
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, tokenData{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  RegularAudience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    TokenIssuer,
			Subject:   u.Username,
		},
		UserId: u.UserId,
	})
	return token.SignedString(config.Secret)
}

func Authentication(next http.Handler) http.Handler {
	config := utils.NewConf()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(userAuthHeader)
		if token == "" {
			utils.EncodeError(w, http.StatusUnauthorized, notAuthed)
			return
		} else if !strings.HasPrefix(token, userTokenPrefix) {
			utils.EncodeError(w, http.StatusBadRequest, wrongHeader)
			return
		}
		// validate token
		data, err := tokenToUser(token[lenPrefix:], config.Secret)
		if err != nil {
			utils.EncodeError(w, http.StatusUnauthorized, invalidToken)
			return
		}
		// query database
		llo := GetLLObject(r)
		user, err := database.GetUserByUUID(r, llo.Conn, data.UserId)
		if errors.Is(err, pgx.ErrNoRows) {
			utils.EncodeError(w, http.StatusUnauthorized, invalidUser)
			return
		} else if err != nil {
			utils.EncodeError(w, http.StatusInternalServerError, serverProblem)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserDataKey, user)
		rr := r.WithContext(ctx)
		next.ServeHTTP(w, rr)
	})
}

func GetUserData(r *http.Request) database.User {
	user, ok := r.Context().Value(ctxUserDataKey).(database.User)
	if !ok {
		slog.Error("user data is corrupted or overwritten")
		os.Exit(1)
	}
	return user
}
