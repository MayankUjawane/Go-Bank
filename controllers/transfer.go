package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/MayankUjawane/gobank/types"
	"github.com/MayankUjawane/gobank/util"
)

type TransferRequest struct {
	FromAccount int   `json:"fromAccount"`
	ToAccount   int   `json:"toAccount"`
	Amount      int64 `json:"amount"`
}

type TransferResponse struct {
	Message string         `json:"message"`
	Account *types.Account `json:"accountDetails"`
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	transferReq := TransferRequest{}
	if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// get the details of the account
	fromAccount, err := s.store.GetAccountByFilter("number", transferReq.FromAccount)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// check balance
	if fromAccount.Balance < int64(transferReq.Amount) {
		util.WriteJSON(w, http.StatusOK, "does not have sufficient balance")
		return
	}

	// get toAccount
	toAccount, err := s.store.GetAccountByFilter("number", transferReq.ToAccount)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// debit amount from the balance of fromAccount
	prevBalance := fromAccount.Balance
	newBalance := prevBalance - transferReq.Amount
	fromAccount, err = s.store.UpdateBalance(transferReq.FromAccount, int(newBalance))
	if err != nil {
		// if error occured than abort the transaction
		s.store.UpdateBalance(transferReq.FromAccount, int(prevBalance))
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// credit amount in the balance of toAccount
	newBalance = toAccount.Balance + transferReq.Amount
	_, err = s.store.UpdateBalance(transferReq.ToAccount, int(newBalance))
	if err != nil {
		// if error occured than abort the transaction
		s.store.UpdateBalance(transferReq.FromAccount, int(prevBalance))
		util.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	transferResponse := &TransferResponse{
		Message: "Transaction Successful",
		Account: fromAccount,
	}

	err = util.WriteJSON(w, http.StatusOK, transferResponse)
	if err != nil {
		util.WriteJSON(w, http.StatusInternalServerError, err)
	}
}
