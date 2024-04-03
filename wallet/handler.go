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
	GetWallets(filter Wallet) ([]Wallet, error)
	CreateWallet(wallet *Wallet) error
	UpdateWallet(wallet *Wallet) error
	DeleteWallet(userID int) error
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

// GetWalletsHandler
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
func (h *Handler) GetWalletsHandler(c echo.Context) error {
	filter := Wallet{}

	// prepare filter: wallet_type
	if walletType := c.QueryParam("wallet_type"); walletType != "" {
		filter.WalletType = walletType
		if !IsWalletTypeValid(walletType) {
			return c.JSON(http.StatusOK, []Wallet{})
		}
	}

	// get wallets
	wallets, err := h.store.GetWallets(filter)
	if err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "error getting wallets"})
	}

	return c.JSON(http.StatusOK, wallets)
}

func (h *Handler) CreateWalletHandler(c echo.Context) error {
	// bind request body to wallet
	wallet := Wallet{}
	if err := c.Bind(&wallet); err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request"})
	}

	// create wallet
	if err := h.store.CreateWallet(&wallet); err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "error creating wallet"})
	}

	return c.JSON(http.StatusCreated, wallet)
}

func (h *Handler) UpdateWalletHandler(c echo.Context) error {
	// bind request body to wallet
	wallet := Wallet{}
	if err := c.Bind(&wallet); err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request"})
	}

	// update wallet
	if err := h.store.UpdateWallet(&wallet); err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "error updating wallet"})
	}

	return c.JSON(http.StatusOK, wallet)
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

	// prepare filter: user_id
	strUserID := c.Param("id")
	if strUserID == "" {
		return c.JSON(http.StatusBadRequest, Err{Message: "user_id is required"})
	}
	if userID, err := strconv.Atoi(strUserID); err == nil {
		filter.UserID = userID
	}

	// prepare filter: wallet_type
	if walletType := c.QueryParam("wallet_type"); walletType != "" {
		filter.WalletType = walletType
		if !IsWalletTypeValid(walletType) {
			return c.JSON(http.StatusOK, []Wallet{})
		}
	}

	// get wallets
	wallets, err := h.store.GetWallets(filter)
	if err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "error getting wallets"})
	}

	return c.JSON(http.StatusOK, wallets)
}

func (h *Handler) DeleteUserWalletHandler(c echo.Context) error {
	// parse user id
	userID, err := ParseUserID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	// delete wallet
	if err = h.store.DeleteWallet(userID); err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "error deleting wallet"})
	}
	return c.NoContent(http.StatusNoContent)
}
