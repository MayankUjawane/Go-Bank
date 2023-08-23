package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/MayankUjawane/gobank/jwt"
	"github.com/MayankUjawane/gobank/util"
)

type loginRequest struct {
	Number   string `json:"number"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	loginRequest := loginRequest{}

	// Get the id and password from request body
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Look up requested user in database
	number, err := util.ConvertIntoInt(loginRequest.Number)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := s.store.GetAccountByFilter("number", number)
	if err != nil {
		util.WriteJSON(w, http.StatusNotFound, err.Error())
		return
	}

	// Compare sent in password with saved user hashed password
	savedPassword := account.HashedPassword
	sentPassword := loginRequest.Password
	err = util.CheckPassword(sentPassword, savedPassword)
	if err != nil {
		util.WriteJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate a JWT Token
	tokenString, err := jwt.CreateJWT(account)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	loginResponse := loginResponse{Token: tokenString}

	// send token back in header
	err = util.WriteJSON(w, http.StatusOK, loginResponse)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
}
