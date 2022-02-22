package main

import (
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	libp2p2 "github.com/ipfs/go-ipfs/core/node/libp2p"

	"github.com/ipfs/go-ipfs/core"
	logging "github.com/ipfs/go-log"
)

// log is the command logger
var log = logging.Logger("cmd/ipfs")

const SwarmKey = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"

func main() {

	ipfsPath, err := fsrepo.BestKnownPath()
	plugins, err := loader.NewPluginLoader(ipfsPath)
	if err != nil {
		log.Errorf("error loading plugins: %s", err)
		os.Exit(0)
	}

	if err := plugins.Initialize(); err != nil {
		log.Errorf("error initializing plugins: %s", err)
		os.Exit(0)
	}

	if err := plugins.Inject(); err != nil {
		log.Errorf("error initializing plugins: %s", err)
		os.Exit(0)
	}

	ctx := context.Background()

	if !fsrepo.IsInitialized(ipfsPath) {
		identity, err := config.CreateIdentity(os.Stdout, []options.KeyGenerateOption{
			options.Key.Type(options.Ed25519Key),
		})
		if err != nil {
			log.Error("create identity error : ", err)
			os.Exit(1)
		}
		conf, err := config.InitWithIdentity(identity)
		if err != nil {
			log.Error("InitWithIdentity error: ", err)
			os.Exit(1)
		}

		conf.Bootstrap = []string{}
		conf.Swarm.EnableAutoRelay = true
		conf.Swarm.EnableRelayHop = true
		err = fsrepo.Init(ipfsPath, conf)
		if err != nil {
			log.Error("fsrepo  init fail : ", err)
			os.Exit(1)
		}
		// init swarm key
		err = ioutil.WriteFile(filepath.Join(ipfsPath, "swarm.key"), []byte(SwarmKey), fs.ModePerm)
		if err != nil {
			log.Error("init swarm.key fail", err)
			os.Exit(1)
		}
	}
	repo, err := fsrepo.Open(ipfsPath)
	if err != nil {
		log.Error("fsrepo is not initialization: ", err)
		os.Exit(0)
	}
	ncfg := &core.BuildCfg{
		Repo:                        repo,
		Permanent:                   true,
		Online:                      true,
		DisableEncryptedConnections: false,
		ExtraOpts: map[string]bool{
			"pubsub": false,
			"ipnsps": false,
		},
		Routing: libp2p2.DHTOption,
	}

	node, err := core.NewNode(ctx, ncfg)
	if err != nil {
		log.Error("error from node construction: ", err)
		os.Exit(0)
	}
	node.IsDaemon = true

	printSwarmAddrs(node)
	select {}
	//defer func() {
	//	// We wait for the node to close first, as the node has children
	//	// that it will wait for before closing, such as the API server.
	//	node.Close()
	//
	//	select {
	//	case <-context.Done():
	//		log.Info("Gracefully shut down daemon")
	//	default:
	//	}
	//}()
}

// printSwarmAddrs prints the addresses of the host
func printSwarmAddrs(node *core.IpfsNode) {
	if !node.IsOnline {
		fmt.Println("Swarm not listening, running in offline mode.")
		return
	}

	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		log.Errorf("failed to read listening addresses: %s", err)
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Strings(lisAddrs)
	for _, addr := range lisAddrs {
		fmt.Printf("Swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Strings(addrs)
	for _, addr := range addrs {
		fmt.Printf("Swarm announcing %s\n", addr)
	}

}
