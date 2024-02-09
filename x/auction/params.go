package auction

import (
	"errors"
	fmt "fmt"
	"strings"
	time "time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultCommitsDuration = 5 * time.Minute
	DefaultRevealsDuration = 5 * time.Minute
	DefaultCommitFee       = sdk.Coin{Amount: sdkmath.NewInt(10), Denom: sdk.DefaultBondDenom}
	DefaultRevealFee       = sdk.Coin{Amount: sdkmath.NewInt(10), Denom: sdk.DefaultBondDenom}
	DefaultMinimumBid      = sdk.Coin{Amount: sdkmath.NewInt(1000), Denom: sdk.DefaultBondDenom}
)

func NewParams(commitsDuration time.Duration, revealsDuration time.Duration, commitFee sdk.Coin, revealFee sdk.Coin, minimumBid sdk.Coin) Params {
	return Params{
		CommitsDuration: commitsDuration,
		RevealsDuration: revealsDuration,
		CommitFee:       commitFee,
		RevealFee:       revealFee,
		MinimumBid:      minimumBid,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		CommitsDuration: DefaultCommitsDuration,
		RevealsDuration: DefaultRevealsDuration,
		CommitFee:       DefaultCommitFee,
		RevealFee:       DefaultRevealFee,
		MinimumBid:      DefaultMinimumBid,
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("CommitsDuration: %s\n", p.CommitsDuration.String()))
	sb.WriteString(fmt.Sprintf("RevealsDuration: %s\n", p.RevealsDuration.String()))
	sb.WriteString(fmt.Sprintf("CommitFee: %s\n", p.CommitFee.String()))
	sb.WriteString(fmt.Sprintf("RevealFee: %s\n", p.RevealFee.String()))
	sb.WriteString(fmt.Sprintf("MinimumBid: %s\n", p.MinimumBid.String()))
	return sb.String()
}

// Validate a set of params.
func (p Params) Validate() error {
	if err := validateCommitsDuration(p.CommitsDuration); err != nil {
		return err
	}

	if err := validateRevealsDuration(p.RevealsDuration); err != nil {
		return err
	}

	if err := validateCommitFee(p.CommitFee); err != nil {
		return err
	}

	if err := validateRevealFee(p.RevealFee); err != nil {
		return err
	}

	if err := validateMinimumBid(p.MinimumBid); err != nil {
		return err
	}

	return nil
}

func validateCommitsDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return errors.New("commits duration cannot be negative")
	}

	return nil
}

func validateRevealsDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return errors.New("reveals duration cannot be negative")
	}

	return nil
}

func validateCommitFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("commit fee must be positive")
	}

	return nil
}

func validateRevealFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("reveal fee must be positive")
	}

	return nil
}

func validateMinimumBid(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Amount.IsNegative() {
		return errors.New("minimum bid must be positive")
	}

	return nil
}
