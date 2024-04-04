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

// CreateWalletHandler
//
//	@Summary		Create wallet
//	@Description	Create wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			wallet	body	WalletForCreate	true	"Wallet object"
//	@Success		201	{object}	Wallet
//	@Failure		400	{object}	Err
//	@Failure		500	{object}	Err
//	@Router			/api/v1/wallets [post]
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

// UpdateWalletHandler
//
//	@Summary		Update wallet
//	@Description	Update wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			wallet	body	WalletForUpdate	true	"Wallet object"
//	@Success		200	{object}	Wallet
//	@Failure		400	{object}	Err
//	@Failure		500	{object}	Err
//	@Router			/api/v1/wallets [put]
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

// GetUserWalletHandler
//
// @Summary		Get all wallets for the user
// @Description	Get all wallets for the user
// @Tags		user wallet
// @Accept		json
// @Produce		json
// @Param		id      path        int true "User ID"
// @Success		200     {array}	    Wallet
// @Failure		400	    {object}	Err
// @Failure		500	    {object}	Err
// @Router		/api/v1/users/{id}/wallets [get]
func (h *Handler) GetUserWalletHandler(c echo.Context) error {
	filter := Wallet{}

	// prepare filter: user_id
	var userID int
	var err error
	strUserID := c.Param("id")
	if strUserID == "" {
		log.Printf("error: user_id is ''\n")
		return c.JSON(http.StatusBadRequest, Err{Message: "user_id is required"})
	}
	if userID, err = strconv.Atoi(strUserID); err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid user_id"})
	}
	filter.UserID = userID

	// prepare filter: wallet_type
	//if walletType := c.QueryParam("wallet_type"); walletType != "" {
	//	filter.WalletType = walletType
	//	if !IsWalletTypeValid(walletType) {
	//		return c.JSON(http.StatusOK, []Wallet{})
	//	}
	//}

	// get wallets
	wallets, err := h.store.GetWallets(filter)
	if err != nil {
		log.Printf("error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "error getting wallets"})
	}

	return c.JSON(http.StatusOK, wallets)
}

//	DeleteUserWalletHandler
//
// @Summary		Delete wallet for the user
// @Description	Delete wallet for the user
// @Tags		user wallet
// @Produce		json
// @Param		id      path        int true "User ID"
// @Success		204
// @Failure		400	    {object}	Err
// @Failure		500	    {object}	Err
// @Router		/api/v1/user/{id}/wallets [delete]
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
