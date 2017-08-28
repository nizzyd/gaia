package stake

import (
	"github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/modules/coin"
)

// Tx
//--------------------------------------------------------------------------------

// register the tx type with its validation logic
// make sure to use the name of the handler as the prefix in the tx type,
// so it gets routed properly
const (
	ByteTxBond     = 0x55
	ByteTxUnbond   = 0x56
	ByteTxNominate = 0x57
	ByteTxModComm  = 0x58
	TypeTxBond     = name + "/bond"
	TypeTxUnbond   = name + "/unbond"
	TypeTxNominate = name + "/nominate"
	TypeTxModComm  = name + "/modComm" //modify commission rate
)

func init() {
	sdk.TxMapper.RegisterImplementation(TxBond{}, TypeTxBond, ByteTxBond)
	sdk.TxMapper.RegisterImplementation(TxUnbond{}, TypeTxUnbond, ByteTxUnbond)
	sdk.TxMapper.RegisterImplementation(TxNominate{}, TypeTxNominate, ByteTxNominate)
	sdk.TxMapper.RegisterImplementation(TxModComm{}, TypeTxModComm, ByteTxModComm)
}

//Verify interface at compile time
var _, _, _, _ sdk.TxInner = &TxBond{}, &TxUnbond{}, &TxNominate{}, &TxModComm{}

/////////////////////////////////////////////////////////////////
// TxBond

// TxBond - struct for bonding transactions
type TxBond struct{ TxBonding }

// NewTxBond - new TxBond
func NewTxBond(delegatee sdk.Actor, amount coin.Coin) sdk.Tx {
	return TxBond{TxBonding{
		Delegatee: delegatee,
		Amount:    amount,
	}}.Wrap()
}

// TxUnbond - struct for unbonding transactions
type TxUnbond struct{ TxBonding }

// NewTxUnbond - new TxUnbond
func NewTxUnbond(delegatee sdk.Actor, amount coin.Coin) sdk.Tx {
	return TxUnbond{TxBonding{
		Delegatee: delegatee,
		Amount:    amount,
	}}.Wrap()
}

// TxBonding - struct for bonding or unbonding transactions
type TxBonding struct {
	Delegatee sdk.Actor `json:"delegatee"`
	Amount    coin.Coin `json:"amount"`
}

// Wrap - Wrap a Tx as a Basecoin Tx
func (tx TxBonding) Wrap() sdk.Tx {
	return sdk.Tx{tx}
}

// ValidateBasic - Check the bonding coins, Validator is non-empty
func (tx TxBonding) ValidateBasic() error {
	if tx.Delegatee.Empty() {
		return errValidatorEmpty
	}
	coins := coin.Coins{tx.Amount}
	if !coins.IsValid() {
		return coin.ErrInvalidCoins()
	}
	if !coins.IsNonnegative() {
		return coin.ErrInvalidCoins()
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// TxNominate

// TxNominate - struct for all staking transactions
type TxNominate struct {
	Nominee    sdk.Actor `json:"nominee"`
	Amount     coin.Coin `json:"amount"`
	Commission Decimal   `json:"commission"`
}

// NewTxNominate - return a new transaction for validator self-nomination
func NewTxNominate(nominee sdk.Actor, amount coin.Coin, commission Decimal) sdk.Tx {
	return TxNominate{
		Nominee:    nominee,
		Amount:     amount,
		Commission: commission,
	}.Wrap()
}

// Wrap - Wrap a Tx as a Basecoin Tx
func (tx TxNominate) Wrap() sdk.Tx {
	return sdk.Tx{tx}
}

// ValidateBasic - Check coins as well as that the delegatee is actually a delegatee
func (tx TxNominate) ValidateBasic() error {
	if tx.Nominee.Empty() {
		return errValidatorEmpty
	}
	coins := coin.Coins{tx.Amount}
	if !coins.IsValid() {
		return coin.ErrInvalidCoins()
	}
	if !coins.IsNonnegative() {
		return coin.ErrInvalidCoins()
	}
	if tx.Commission.LT(NewDecimal(0, 1)) {
		return errCommissionNegative
	}
	if tx.Commission.GT(NewDecimal(1, 1)) {
		return errCommissionHuge
	}
	return nil
}

/////////////////////////////////////////////////////////////////
// TxModComm

// TxModComm - struct for all staking transactions
type TxModComm struct {
	Delegatee  sdk.Actor `json:"delegatee"`
	Commission Decimal   `json:"commission"`
}

// NewTxModComm - return a new counter transaction struct wrapped as a sdk transaction
func NewTxModComm(delegatee sdk.Actor, commission Decimal) sdk.Tx {
	return TxModComm{
		Delegatee:  delegatee,
		Commission: commission,
	}.Wrap()
}

// Wrap - Wrap a Tx as a Basecoin Tx
func (tx TxModComm) Wrap() sdk.Tx {
	return sdk.Tx{tx}
}

// ValidateBasic - Check coins as well as that the delegatee is actually a delegatee
func (tx TxModComm) ValidateBasic() error {
	if tx.Delegatee.Empty() {
		return errValidatorEmpty
	}
	if tx.Commission.LT(NewDecimal(0, 1)) {
		return errCommissionNegative
	}
	if tx.Commission.GT(NewDecimal(1, 1)) {
		return errCommissionHuge
	}
	return nil
}
