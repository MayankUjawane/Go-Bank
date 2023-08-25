package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MayankUjawane/gobank/token"
	"github.com/MayankUjawane/gobank/util"
	"github.com/gorilla/context"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// checking if token is passed in the header or not
		authorizationHeader := r.Header.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		// this will split our authorizationHeader around spce, so we should get 2 parts
		// one is Bearer and second is token
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			util.WriteJSON(w, http.StatusUnauthorized, err)
			return
		}

		context.Set(r, authorizationPayloadKey, payload)

		// after authentication callind the function
		next.ServeHTTP(w, r)
	}
}
