package context

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/hamster-shared/hamster-gateway/core/modules/chain"
	"github.com/hamster-shared/hamster-gateway/core/modules/config"
	"github.com/hamster-shared/hamster-gateway/core/modules/time"
	"github.com/hamster-shared/hamster-gateway/core/modules/utils"
)

// CoreContext the application context , wrapped with some bean
type CoreContext struct {
	Cm           *config.ConfigManager
	ReportClient chain.ReportClient
	SubstrateApi *gsrpc.SubstrateAPI
	TimerService *utils.TimerService
	StateService *time.StateService
}

func (c *CoreContext) GetConfig() *config.Config {
	cf, _ := c.Cm.GetConfig()
	return cf
}
