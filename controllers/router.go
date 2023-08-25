package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MayankUjawane/gobank/db"
	"github.com/MayankUjawane/gobank/token"
	"github.com/MayankUjawane/gobank/types"
	"github.com/MayankUjawane/gobank/util"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         db.Storage
	tokenMaker    token.Maker
}

func NewAPIServer(listenAddress string, store db.Storage, tokenMaker token.Maker) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
		tokenMaker:    tokenMaker,
	}
}

// Router will create new router and will handle all the routes
func (s *APIServer) SetupRouter(tokenMaker token.Maker) {
	router := mux.NewRouter()

	router.HandleFunc("/signup", s.handleCreateAccount).Methods("POST")
	router.HandleFunc("/login", s.handleLogin).Methods("POST")
	router.HandleFunc("/account", s.handleGetAccount).Methods("GET")

	router.HandleFunc("/account/{id}", authMiddleware(tokenMaker, s.handleGetAccountByID)).Methods("GET")
	router.HandleFunc("/account/{id}", authMiddleware(tokenMaker, s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/transfer", s.handleTransfer).Methods("POST")

	log.Println("JSON API server running on port: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.GetAllAccounts()
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, accounts); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := util.GetId(r)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := s.store.GetAccountByFilter("id", id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// check if user is asking for his data only or not
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

	err = util.WriteJSON(w, http.StatusOK, account)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}

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

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	transferReq := types.TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err := util.WriteJSON(w, http.StatusOK, transferReq)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}
}
