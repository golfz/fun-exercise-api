package wallet

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
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
//	    @Param			wallet_type     query       string false "Filter by wallet type"
//		@Success		200	            {array}	    Wallet
//		@Failure		500	            {object}	Err
//		@Router			/api/v1/wallets [get]
func (h *Handler) WalletHandler(c echo.Context) error {
	filter := Wallet{}

	// filter by wallet_type
	if walletType := c.QueryParam("wallet_type"); walletType != "" {
		filter.WalletType = walletType
		log.Printf("filter by wallet_type=%s\n", walletType)
	}

	wallets, err := h.store.Wallets(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}

// UserWalletHandler
//
//			@Summary		Get all wallets for the user
//			@Description	Get all wallets for the user
//			@Tags			wallet
//			@Accept			json
//			@Produce		json
//		    @Param			id      path        int true "User ID"
//			@Success		200     {array}	    Wallet
//	        @Failure		400	    {object}	Err
//			@Failure		500	    {object}	Err
//			@Router			/api/v1/user/{id}/wallets [get]
func (h *Handler) UserWalletHandler(c echo.Context) error {
	filter := Wallet{}

	strUserID := c.Param("id")
	if strUserID == "" {
		return c.JSON(http.StatusBadRequest, Err{Message: "user_id is required"})
	}
	if userID, err := strconv.Atoi(strUserID); err == nil {
		filter.UserID = userID
	}

	// filter by wallet_type
	if walletType := c.QueryParam("wallet_type"); walletType != "" {
		filter.WalletType = walletType
		log.Printf("filter by wallet_type=%s\n", walletType)
	}

	wallets, err := h.store.Wallets(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}
