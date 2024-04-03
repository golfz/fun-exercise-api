package wallet

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockWalletStorer struct {
	wallets      []Wallet
	err          error
	methodToCall map[string]bool
}

func (m *mockWalletStorer) Wallets() ([]Wallet, error) {
	m.methodToCall["Wallets"] = true
	return m.wallets, m.err
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

func TestWallet(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		mock := &mockWalletStorer{
			err: errors.New("unable to get wallets"),
		}
		mock.ExpectToCall("Wallets")
		h := New(mock)

		// Act
		err := h.WalletHandler(c)

		// Assert
		mock.Verify(t)
		if err != nil {
			t.Errorf("expected err to be nil, got %s", err)
		}
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status code to be 500, got %d", rec.Code)
		}
		if len(rec.Body.String()) == 0 {
			t.Errorf("expected response body to be not empty")
		}
	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
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
		mock := &mockWalletStorer{
			wallets: want,
		}
		mock.ExpectToCall("Wallets")
		h := New(mock)

		// Act
		err := h.WalletHandler(c)

		// Assert
		mock.Verify(t)
		if err != nil {
			t.Errorf("expected err to be nil, got %s", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("expected status code to be 200, got %d", rec.Code)
		}
		var got []Wallet
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("expected response body to be valid json, got %s", rec.Body.String())
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected response body to be %v, got %v", want, got)
		}
	})
}
