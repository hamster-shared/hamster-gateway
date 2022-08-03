package time

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	chain2 "github.com/hamster-shared/hamster-gateway/core/modules/chain"
	"github.com/hamster-shared/hamster-gateway/core/modules/config"
	"github.com/hamster-shared/hamster-gateway/core/modules/p2p"
	"github.com/ipfs/go-ipfs/core"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

type StateService struct {
	cancel       func()
	Node         *core.IpfsNode
	reportClient chain2.ReportClient
	ctx          context.Context
	cm           *config.ConfigManager
}

func NewStateService(cm *config.ConfigManager) *StateService {
	return &StateService{
		cm: cm,
	}
}

func (s *StateService) Start() error {

	if s.Running() {
		return nil
	}
	cf, err := s.cm.GetConfig()

	substrateApi, err := gsrpc.NewSubstrateAPI(cf.ChainApi)
	if err != nil {
		log.Error(err)
		return err
	}
	reportClient, err := chain2.NewChainClient(s.cm, substrateApi)
	if err != nil {
		log.Error(err)
		return err
	}

	s.reportClient = reportClient

	s.ctx, s.cancel = context.WithCancel(context.Background())

	node, err := p2p.RunDaemon(s.ctx)
	if err != nil {
		log.Error("run ipfs daemon fail")
		return err
	}

	s.Node = node

	cf.PeerId = node.Identity.String()
	_ = s.cm.Save(cf)

	localAddress := fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", cf.PublicIp, cf.PublicPort, node.Identity.String())

	// 2: blockchain registration
	marketUser, err := s.reportClient.GetMarketUser()
	fmt.Println(marketUser)
	if err != nil {
		err := s.reportClient.CrateMarketAccount()
		if err != nil {
			return err
		}
	}
	err = s.reportClient.Register(localAddress)

	if err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Minute)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				_ = s.reportClient.Heartbeat(localAddress)
			}
		}
	}(s.ctx)

	return nil
}

func (s *StateService) Stop() {

	if !s.Running() {
		return
	}

	s.cancel()
	_ = s.Node.Close()
	s.cancel = nil
	s.Node = nil
}

func (s *StateService) Running() bool {
	return s.cancel != nil
}
