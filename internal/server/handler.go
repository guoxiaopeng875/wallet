package server

import (
	"github.com/guoxiaopeng875/wallet/internal/wallet"
	"net/http"
)

// Handler handles HTTP requests for wallet operations
type Handler struct {
	uc wallet.UseCase
}

func NewHandler(uc wallet.UseCase) *Handler {
	return &Handler{uc: uc}
}

// Deposit handles wallet deposit requests
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	id, req := parseWalletID(w, r), &DepositRequest{}
	if id == 0 || !parseReqBody(w, r, req) {
		return
	}

	if err := h.uc.Deposit(r.Context(), id, req.Amount); err != nil {
		handleError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, nil)
}

// Withdraw handles wallet withdrawal requests
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	id, req := parseWalletID(w, r), &WithdrawRequest{}
	if id == 0 || !parseReqBody(w, r, req) {
		return
	}

	if err := h.uc.Withdraw(r.Context(), id, req.Amount); err != nil {
		handleError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, nil)
}

// Transfer handles wallet transfer requests
func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	id, req := parseWalletID(w, r), &TransferRequest{}
	if id == 0 || !parseReqBody(w, r, req) {
		return
	}

	if err := h.uc.Transfer(r.Context(), id, req.TargetWalletID, req.Amount); err != nil {
		handleError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, nil)
}

// Balance retrieves wallet balance
func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	id := parseWalletID(w, r)
	if id == 0 {
		return
	}

	wallet, err := h.uc.Wallet(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, &BalanceResponse{Balance: wallet.Balance.String()})
}

// Transactions retrieves wallet transaction history
func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	id := parseWalletID(w, r)
	if id == 0 {
		return
	}

	txs, err := h.uc.WalletTransactions(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, txs)
}
