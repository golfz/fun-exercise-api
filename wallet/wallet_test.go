package wallet

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
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

func TestWallet(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets", nil)
		mock.err = errors.New("unable to get wallets")
		mock.ExpectToCall("GetWallets")

		// Act
		err := h.WalletHandler(c)

		// Assert
		mock.Verify(t)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		var got Err
		if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", resp.Body.String())
		}
		assert.Equal(t, "unable to get wallets", got.Message)
	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets", nil)
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
		err := h.WalletHandler(c)

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

	t.Run("given user filter by available wallet types should return list of wallets", func(t *testing.T) {
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
		err := h.WalletHandler(c)

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

	t.Run("given user filter by unavailable wallet types should return empty list of wallet", func(t *testing.T) {
		// Arrange
		resp, c, h, mock := testSetup(http.MethodGet, "/api/v1/wallets?wallet_type=Unknown", nil)
		expectedFilter := Wallet{
			WalletType: "Unknown",
		}
		want := []Wallet{}
		mock.wallets = want
		mock.ExpectToCall("GetWallets")

		// Act
		err := h.WalletHandler(c)

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

}

func TestUserWallet(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {

	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {

	})
}
