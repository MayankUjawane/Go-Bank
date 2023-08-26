package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MayankUjawane/gobank/token"
	"github.com/MayankUjawane/gobank/util"
	"github.com/gorilla/context"
)

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := util.GetId(r)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// check if user is asking for his data only or not
	account, err := s.store.GetAccountByFilter("id", id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	authPayload := context.Get(r, authorizationPayloadKey).(*token.Payload)
	authNumber, err := strconv.Atoi(authPayload.Number)
	if err != nil {
		error := fmt.Errorf("while conversion in handleAccountGetByID: %s", err.Error())
		util.WriteJSON(w, http.StatusInternalServerError, error)
		return
	}
	if account.Number != int64(authNumber) {
		util.WriteJSON(w, http.StatusUnauthorized, "account doesn't belong to user")
		return
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = util.WriteJSON(w, http.StatusOK, map[string]int{"account deleted": id})
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}
