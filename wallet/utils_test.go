package wallet

import "testing"

func TestIsWalletTypeValid(t *testing.T) {
	tests := []struct {
		name       string
		walletType string
		want       bool
	}{
		{
			name:       "Empty wallet type",
			walletType: "",
			want:       false,
		},
		{
			name:       "Valid wallet type",
			walletType: "Savings",
			want:       true,
		},
		{
			name:       "Invalid wallet type",
			walletType: "INVALID_TYPE",
			want:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsWalletTypeValid(test.walletType); got != test.want {
				t.Errorf("IsWalletTypeValid() = %v, want %v", got, test.want)
			}
		})
	}
}
