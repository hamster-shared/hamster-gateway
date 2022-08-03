package chain

// ReportClient data reporting interface
type ReportClient interface {
	// Register  resource registration
	Register(string) error
	// Heartbeat protocol heartbeat report
	Heartbeat(localhostAddress string) error

	GetMarketUser() (MarketUser, error)

	CrateMarketAccount() error
}
