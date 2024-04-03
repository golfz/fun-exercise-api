package wallet

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Handler struct {
	store Storer
}

type Filter struct {
	WalletType string
}

type Storer interface {
	Wallets(filter Wallet) ([]Wallet, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

// WalletHandler
//
//		@Summary		Get all wallets
//		@Description	Get all wallets
//		@Tags			wallet
//		@Accept			json
//		@Produce		json
//	    @Param			wallet_type query string false "Filter by wallet type"
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/wallets [get]
//		@Failure		500	{object}	Err
func (h *Handler) WalletHandler(c echo.Context) error {
	filter := Wallet{}
	if walletType := c.QueryParams().Get("wallet_type"); walletType != "" {
		filter.WalletType = walletType
		log.Printf("filter by wallet_type: %s", walletType)
	}
	wallets, err := h.store.Wallets(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}
