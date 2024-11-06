package wallet

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestWallet_CheckBalance(t *testing.T) {
	tests := []struct {
		name    string
		wallet  *Wallet
		amount  decimal.Decimal
		wantErr bool
	}{
		{
			name: "sufficient balance",
			wallet: &Wallet{
				ID:      1,
				Balance: decimal.NewFromFloat(100),
			},
			amount:  decimal.NewFromFloat(50),
			wantErr: false,
		},
		{
			name: "insufficient balance",
			wallet: &Wallet{
				ID:      1,
				Balance: decimal.NewFromFloat(100),
			},
			amount:  decimal.NewFromFloat(150),
			wantErr: true,
		},
		{
			name: "exact balance",
			wallet: &Wallet{
				ID:      1,
				Balance: decimal.NewFromFloat(100),
			},
			amount:  decimal.NewFromFloat(100),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.wallet.CheckBalance(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
