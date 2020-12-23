package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/nozomi-nishinohara/jwt_validation/domain/model"
	"github.com/nozomi-nishinohara/jwt_validation/domain/repository"
	"github.com/nozomi-nishinohara/jwt_validation/infrastructure"
	tk "github.com/nozomi-nishinohara/jwt_validation/token"
)

type (
	handler struct {
		repo repository.Cache
	}
	IHttpHandle func(http.ResponseWriter, *http.Request)
	IHandler    interface {
		JWTValidationMiddleware(h http.HandlerFunc) http.HandlerFunc
		JWTValidation(h http.Handler) http.Handler
	}
)

func New(repo repository.Cache) IHandler {
	if repo == nil {
		switch model.GetSetting().Cache.Name {
		case model.InMemory:
			repo = infrastructure.NewInMemory()
		case model.Redis:
			repo = infrastructure.NewRedis()
		}
	}
	if err := savePemCert(repo); err != nil {
		panic(err)
	}
	return &handler{
		repo: repo,
	}
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	buf, _ := json.Marshal(body)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func responseUnAuthorized(w http.ResponseWriter, meesage string, value ...interface{}) {
	body := map[string]interface{}{
		"code":    http.StatusUnauthorized,
		"error":   "Unauthorized",
		"message": fmt.Sprintf(meesage, value...),
	}
	writeJSON(w, http.StatusUnauthorized, body)
}

func missingAuthorization(w http.ResponseWriter) {
	responseUnAuthorized(w, "Missing authentication")
}

func (handler *handler) valid(w http.ResponseWriter, req *http.Request) (context.Context, bool) {
	ctx := req.Context()
	authorization := ""
	if authorization = req.Header.Get("Authorization"); authorization == "" {
		missingAuthorization(w)
		return req.Context(), false
	}

	splitAuthorization := strings.Split(authorization, " ")
	if len(splitAuthorization) == 2 && splitAuthorization[0] == "Bearer" && splitAuthorization[1] != "" {
		jwtStr := splitAuthorization[1] // JWTを代入
		var token *jwt.Token
		var err error
		if token, err = handler.jwtParse(ctx, jwtStr); err != nil {
			missingAuthorization(w)
			return req.Context(), false
		}
		flg := true
		for _, oauth := range model.GetSetting().Oauths {
			mapClaims := token.Claims.(jwt.MapClaims)
			if checkIss := mapClaims.VerifyIssuer(oauth.Iss, false); !checkIss {
				flg = false
			}
			audFlg := true
			if len(oauth.Aud) > 0 {
				audFlg = false
			}
			for _, aud := range oauth.Aud {
				if checkAud := mapClaims.VerifyAudience(aud, false); checkAud {
					audFlg = true
					break
				}
			}
			if !audFlg {
				flg = audFlg
			}
			if err := mapClaims.Valid(); err != nil {
				flg = false
			}
		}
		if !flg {
			missingAuthorization(w)
			return req.Context(), false
		}
		return tk.SetContext(req.Context(), token.Claims), true
	} else {
		missingAuthorization(w)
		return req.Context(), false
	}
}

func (handler *handler) JWTValidationMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if c, ok := handler.valid(w, req); ok {
			req = req.WithContext(c)
			h(w, req)
		}
	}
}

func (handler *handler) JWTValidation(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if c, ok := handler.valid(w, req); ok {
			req = req.WithContext(c)
			h.ServeHTTP(w, req)
		}
	})
}

func (handler *handler) jwtParse(c context.Context, jwtStr string) (*jwt.Token, error) {
	return jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		if kid, ok := token.Header["kid"]; ok {
			if kidStr, ok := kid.(string); ok {
				var err error
				var jsonWebKeys *model.JSONWebKeys
				// キャッシュから取得
				if jsonWebKeys, err = handler.repo.Get(c, kidStr); err != nil && err == infrastructure.ErrNotFound {
					// キャッシュから取得出来ない場合は再度JWKSの取得を行う
					if err = savePemCert(handler.repo); err != nil {
						return nil, err
					} else {
						jsonWebKeys, err = handler.repo.Get(c, kidStr)
					}
				}
				if err != nil {
					return nil, err
				} else {
					// キャッシュの更新
					handler.repo.Save(c, kidStr, jsonWebKeys)
					rsaPublicKey := jsonWebKeys.GetPublicKey()
					return rsaPublicKey, nil
				}
			}
		}
		return "", nil
	})
}
