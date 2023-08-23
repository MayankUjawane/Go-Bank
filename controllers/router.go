package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MayankUjawane/gobank/db"
	"github.com/MayankUjawane/gobank/jwt"
	"github.com/MayankUjawane/gobank/types"
	"github.com/MayankUjawane/gobank/util"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         db.Storage
}

func NewAPIServer(listenAddress string, store db.Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

// Router will create new router and will handle all the routes
func (s *APIServer) Router() {
	router := mux.NewRouter()
	router.HandleFunc("/login", s.handleLogin).Methods("POST")
	router.HandleFunc("/account", s.handleCreateAccount).Methods("POST")
	router.HandleFunc("/account", s.handleGetAccount).Methods("GET")
	router.HandleFunc("/account/{id}", jwt.WithJWTAuth(s.handleGetAccountByID, s.store)).Methods("GET")
	router.HandleFunc("/account", s.handleDeleteAccount).Methods("DELETE")
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

	if err := util.WriteJSON(w, http.StatusOK, account); err != nil {
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
