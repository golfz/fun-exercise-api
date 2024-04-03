package wallet

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

func IsWalletTypeValid(walletType string) bool {
	for _, t := range AvailableWalletTypes {
		if t == walletType {
			return true
		}
	}
	return false
}

func ParseUserID(c echo.Context) (int, error) {
	id := c.Param("id")
	if id == "" {
		return 0, errors.New("id is required")
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("invalid user id")
	}

	return userID, nil
}
