package authentication

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nkien0204/projectTemplate/configs"
	"github.com/nkien0204/projectTemplate/internal/log"
	"github.com/nkien0204/projectTemplate/internal/network/http_handler/responses"
	"go.uber.org/zap"
)

const RefreshAllowTime int = 30
const TokenExpireTime int = 45
const AccessTokenKey string = "AccessToken"
const RefreshTokenKey string = "RefreshToken"

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"user_name"`
	jwt.StandardClaims
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	jwtKey := configs.GetConfigs().SecretKey.Key
	var creds Credentials
	logger := log.Logger()
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		logger.Error("decode request failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// just for testing
	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		logger.Error("username or password is not correct")
		responses.CustomResponse(w, responses.ResAuthFailed, "wrong username or password", nil)
		return
	}

	accessToken, err := generateToken(jwtKey, time.Now().Add(time.Duration(TokenExpireTime)*time.Second), creds.Username)
	if err != nil {
		logger.Error("generate accessToken failed")
		responses.CustomResponse(w, responses.ResGenTokenFailed, "generate accessToken failed", nil)
		return
	}
	refreshToken, err := generateToken(jwtKey, time.Now().Add(time.Duration(TokenExpireTime+RefreshAllowTime)*time.Second), creds.Username)
	if err != nil {
		logger.Error("generate refreshToken failed")
		responses.CustomResponse(w, responses.ResGenTokenFailed, "generate refreshToken failed", nil)
		return
	}

	logger.Info("sign-in successfully", zap.String("accessToken", accessToken), zap.String("refreshToken", refreshToken))
	data := map[string]string{
		AccessTokenKey:  accessToken,
		RefreshTokenKey: refreshToken,
	}
	responses.CustomResponse(w, responses.ResOk, "ok", data)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger()
	jwtKey := configs.GetConfigs().SecretKey.Key

	tknStr := r.Header.Get(AccessTokenKey)

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error("invalid signature", zap.Error(err))
			responses.CustomResponse(w, responses.ResInvalidSignature, "invalid signature", nil)
			return
		}
		logger.Error("parse token failed", zap.Error(err))
		responses.CustomResponse(w, responses.ResParseTokenFailed, "parse token failed", nil)
		return
	}
	if !tkn.Valid {
		logger.Error("invalid token")
		responses.CustomResponse(w, responses.ResInvalidToken, "invalid token", nil)
		return
	}

	data := map[string]interface{}{
		"username":         claims.Username,
		"expire time left": claims.ExpiresAt - time.Now().Unix(),
	}
	responses.CustomResponse(w, responses.ResOk, "ok", data)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	logger := log.Logger()
	jwtKey := configs.GetConfigs().SecretKey.Key
	tknStr := r.Header.Get(RefreshTokenKey)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error("invalid signature", zap.Error(err))
			responses.CustomResponse(w, responses.ResInvalidSignature, "invalid signature", nil)
			return
		}
		logger.Error("parse token failed", zap.Error(err))
		responses.CustomResponse(w, responses.ResParseTokenFailed, "parse token failed", nil)
		return
	}
	if !tkn.Valid {
		logger.Error("invalid token", zap.Error(err))
		responses.CustomResponse(w, responses.ResInvalidToken, "invalid token", nil)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).After(time.Now()) {
		logger.Error("no need to refresh token", zap.Error(err))
		responses.CustomResponse(w, responses.ResInvalidToken, "no need to refresh token", nil)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Add(time.Duration(RefreshAllowTime) * time.Second).Before(time.Now()) {
		logger.Error("refresh token was expired", zap.Error(err))
		responses.CustomResponse(w, responses.ResTokenExpired, "refresh token was expired", nil)
		return
	}

	accessToken, err := generateToken(jwtKey, time.Now().Add(time.Duration(TokenExpireTime)*time.Second), claims.Username)
	if err != nil {
		logger.Error("generate accessToken failed")
		responses.CustomResponse(w, responses.ResGenTokenFailed, "generate accessToken failed", nil)
		return
	}
	refreshToken, err := generateToken(jwtKey, time.Now().Add(time.Duration(TokenExpireTime+RefreshAllowTime)*time.Second), claims.Username)
	if err != nil {
		logger.Error("generate refreshToken failed")
		responses.CustomResponse(w, responses.ResGenTokenFailed, "generate refreshToken failed", nil)
		return
	}

	data := map[string]string{
		AccessTokenKey:  accessToken,
		RefreshTokenKey: refreshToken,
	}
	responses.CustomResponse(w, responses.ResOk, "ok", data)
}

func generateToken(jwtKey []byte, expirationTime time.Time, username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
