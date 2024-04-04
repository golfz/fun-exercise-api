package wallet

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockWalletStorer struct {
	wallets      []Wallet
	err          error
	methodToCall map[string]bool
	whatIsFilter Wallet
}

func NewMockWalletStorer() *mockWalletStorer {
	return &mockWalletStorer{
		methodToCall: make(map[string]bool),
	}
}

func (m *mockWalletStorer) GetWallets(filter Wallet) ([]Wallet, error) {
	m.methodToCall["GetWallets"] = true
	m.whatIsFilter = filter
	return m.wallets, m.err
}

func (m *mockWalletStorer) CreateWallet(w *Wallet) error {
	m.methodToCall["CreateWallet"] = true
	return m.err
}

func (m *mockWalletStorer) UpdateWallet(w *Wallet) error {
	m.methodToCall["UpdateWallet"] = true
	return m.err
}

func (m *mockWalletStorer) DeleteWallet(userID int) error {
	m.methodToCall["DeleteWallet"] = true
	return m.err
}

func (m *mockWalletStorer) ExpectToCall(methodName string) {
	if m.methodToCall == nil {
		m.methodToCall = make(map[string]bool)
	}
	m.methodToCall[methodName] = false
}

func (m *mockWalletStorer) Verify(t *testing.T) {
	for methodName, called := range m.methodToCall {
		if !called {
			t.Errorf("expected %s to be called", methodName)
		}
	}
}

func testSetup(method, url string, body io.Reader) (*httptest.ResponseRecorder, echo.Context, *Handler, *mockWalletStorer) {
	req := httptest.NewRequest(method, url, body)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	mock := NewMockWalletStorer()
	h := New(mock)

	return rec, c, h, mock
}

func TestGetWallets(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets", nil)
		mock.err = errors.New("unable to get wallets")
		mock.ExpectToCall("GetWallets")

		// Act
		err := h.GetWalletsHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets", nil)
		log.Printf("c: %#v\n", c)
		want := []Wallet{
			{
				ID:       1,
				UserName: "user1",
				Balance:  1000,
			},
			{
				ID:       2,
				UserName: "user2",
				Balance:  2000,
			},
		}
		mock.wallets = want
		mock.ExpectToCall("GetWallets")

		// Act
		err := h.GetWalletsHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		var got []Wallet
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, want, got)
	})

	t.Run("given user filter by available wallet_types should return 200 and list of wallets", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets?wallet_type=Savings", nil)
		expectedFilter := Wallet{
			WalletType: WalletTypeSavings,
		}
		want := []Wallet{
			{
				ID:       1,
				UserName: "user1",
				Balance:  1000,
			},
			{
				ID:       2,
				UserName: "user2",
				Balance:  2000,
			},
		}
		mock.wallets = want
		mock.ExpectToCall("GetWallets")

		// Act
		err := h.GetWalletsHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, expectedFilter, mock.whatIsFilter)
		assert.Equal(t, http.StatusOK, resp.Code)
		var got []Wallet
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, want, got)
	})

	t.Run("given user filter by unavailable wallet_types should return 200 and empty list of wallet", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets?wallet_type=Unknown", nil)
		want := []Wallet{}
		mock.wallets = want

		// Act
		err := h.GetWalletsHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		var got []Wallet
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, want, got)
	})

}

func TestGetUserWallet(t *testing.T) {
	t.Run("given no user_id in path param should return 400 and error message", func(t *testing.T) {
		// Arrange
		resp, c, h, _ := testSetup(http.MethodGet, "/", nil)
		c.SetPath("/api/v1/users/:id/wallets")
		c.SetParamNames("id")
		c.SetParamValues("")
		// see: https://echo.labstack.com/docs/testing#getuser

		// Act
		err := h.GetUserWalletHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
	})

	t.Run("given user_id is not number should return 400 and error message", func(t *testing.T) {
		// Arrange
		resp, c, h, _ := testSetup(http.MethodGet, "/", nil)
		c.SetPath("/api/v1/users/:id/wallets")
		c.SetParamNames("id")
		c.SetParamValues("abc")

		// Act
		err := h.GetUserWalletHandler(c)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
	})

	t.Run("given error should return 500 and error message", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/", nil)
		c.SetPath("/api/v1/users/:id/wallets")
		c.SetParamNames("id")
		c.SetParamValues("999")
		mock.err = errors.New("unable to get wallets")
		mock.ExpectToCall("GetWallets")
		expectedFilter := Wallet{
			UserID: 999,
		}

		// Act
		err := h.GetUserWalletHandler(c)

		// Assert
		mock.Verify(t)
		assert.Equal(t, expectedFilter, mock.whatIsFilter)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.NotEmpty(t, got.Message)
	})

	t.Run("given no error should return 200 and []wallets", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/", nil)
		c.SetPath("/api/v1/users/:id/wallets")
		c.SetParamNames("id")
		c.SetParamValues("1")
		mock.ExpectToCall("GetWallets")
		expectedFilter := Wallet{
			UserID: 1,
		}
		expectedWallets := []Wallet{
			{
				ID:       1,
				UserName: "user1",
				Balance:  1000,
			},
			{
				ID:       2,
				UserName: "user2",
				Balance:  2000,
			},
		}
		mock.wallets = expectedWallets

		// Act
		err := h.GetUserWalletHandler(c)

		// Assert
		mock.Verify(t)
		assert.Equal(t, expectedFilter, mock.whatIsFilter)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		var got []Wallet
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, expectedWallets, got)
	})
}
