package time

import (
	"fmt"
	chain2 "github.com/hamster-shared/hamster-gateway/core/modules/chain"
	"github.com/hamster-shared/hamster-gateway/core/modules/config"
	"github.com/hamster-shared/hamster-gateway/core/modules/p2p"
	"github.com/ipfs/go-ipfs/core"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"time"
)

type StateService struct {
	cancel       func()
	Node         *core.IpfsNode
	reportClient chain2.ReportClient
	ctx          context.Context
	cm           *config.ConfigManager
}

func NewStateService(reportClient chain2.ReportClient, cm *config.ConfigManager) *StateService {
	return &StateService{
		reportClient: reportClient,
		cm:           cm,
	}
}

func (s *StateService) Start() {

	if s.Running() {
		return
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	node, err := p2p.RunDaemon(s.ctx)
	s.Node = node
	if err != nil {
		log.Error("run ipfs daemon fail")
		os.Exit(1)
	}

	cf, err := s.cm.GetConfig()
	if err != nil {
		log.Error("run ipfs daemon fail")
		os.Exit(1)
	}

	localAddress := fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", cf.PublicIp, cf.PublicPort, node.Identity.String())

	// 2: blockchain registration
	for {
		err := s.reportClient.Register(localAddress)
		if err != nil {
			log.Errorf("Blockchain registration failed, the reason for the failure： %s", err.Error())
			time.Sleep(time.Second * 30)
		} else {
			break
		}
	}

	// 3: healthcheck
	myTimer := time.NewTimer(time.Minute * 10) // start timer

	go func(ctx context.Context) {
		for {
			select {
			case <-myTimer.C:
				// health check
				s.reportClient.Heartbeat(localAddress)
				myTimer.Reset(time.Minute * 10) // reset timer
			case <-ctx.Done():
				return
			}
		}
	}(s.ctx)

}

func (s *StateService) Stop() {

	if !s.Running() {
		return
	}

	s.cancel()
	s.cancel = nil
}

func (s *StateService) Running() bool {
	return s.cancel != nil
}
