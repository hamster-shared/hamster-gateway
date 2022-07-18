package chain

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/hamster-shared/hamster-gateway/core/modules/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterResource(t *testing.T) {

	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)

	assert.NoError(t, err)

	peerId := "abc123"

	err = cc.Register(peerId)
	assert.NoError(t, err)
}

func TestHealthbeat(t *testing.T) {

	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)

	assert.NoError(t, err)

	peerId := "abc123"

	err = cc.Heartbeat(peerId)
	assert.NoError(t, err)
}

func TestResource(t *testing.T) {
	cm := config.NewConfigManager()
	cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI(cfg.ChainApi)
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)

	resource, err := cc.GetResource(41)
	assert.NoError(t, err)
	fmt.Println(resource)
}

func TestGetEventRecord(t *testing.T) {
	cm := config.NewConfigManager()
	//cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI("ws://192.168.1.105:9944")
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)

	events, err := cc.GetEvent(202)

	fmt.Println(err)

	fmt.Println(len(events.Gateway_RegisterGatewayNodeSuccess))

	assert.NoError(t, err)
}

func TestGetState(t *testing.T) {
	cm := config.NewConfigManager()
	//cfg, _ := cm.GetConfig()
	substrateApi, err := gsrpc.NewSubstrateAPI("ws://127.0.0.1:9944")
	cc, err := NewChainClient(cm, substrateApi)
	assert.NoError(t, err)
	marketUser, err := cc.GetMarketUser()
	fmt.Println(marketUser)
	assert.NoError(t, err)
}
