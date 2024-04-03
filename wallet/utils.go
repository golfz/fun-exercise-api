package wallet

func IsWalletTypeValid(walletType string) bool {
	for _, t := range AvailableWalletTypes {
		if t == walletType {
			return true
		}
	}
	return false
}
