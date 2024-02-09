package auction

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:   DefaultParams(),
		Auctions: []*Auction{},
	}
}

func NewGenesisState(params Params, auctions []*Auction) *GenesisState {
	return &GenesisState{
		Params:   params,
		Auctions: auctions,
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return nil
}
