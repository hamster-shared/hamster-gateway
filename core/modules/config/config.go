package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var packageLock sync.Mutex

const (
	CONFIG_DIR_NAME         = ".hamster-gateway"
	CONFIG_DEFAULT_FILENAME = "config"
	SWARM_KEY               = "/key/swarm/psk/1.0.0/\n/base16/\n55158d9b6b7e5a8e41aa8b34dd057ff1880e38348613d27ae194ad7c5b9670d7"
)

// Config  config parameter
type Config struct {
	ApiPort      int    `json:"apiPort"`      // API port number
	ChainApi     string `json:"chainApi"`     // blockchain address
	SeedOrPhrase string `json:"seedOrPhrase"` // blockchain account seed or mnemonic
	PublicIp     string `json:"publicIp"`
	PublicPort   int    `json:"publicPort"`
	PeerId       string `json:"peer_id"`
}

type ConfigFlag string

const DONE ConfigFlag = "done"
const NONE ConfigFlag = "none"

// VmOption vm configuration information
type VmOption struct {
	Cpu        uint64 `json:"cpu"`
	Mem        uint64 `json:"mem"`
	Disk       uint64 `json:"disk"`
	System     string `json:"system"`
	Image      string `json:"image"`
	AccessPort int    `json:"accessPort"`
	// virtualization type,docker/kvm
	Type string `json:"type"`
}

// Identity p2p identity token structure
type Identity struct {
	PeerID   string
	PrivKey  string `json:",omitempty"`
	SwarmKey string `json:"swarm_key"`
}

// PublicKey public key information
type PublicKey struct {
	Key string `json:"key"`
}

type ConfigManager struct {
	configPath string
}

type ChainRegInfo struct {
	ResourceIndex   uint64 `json:"resourceIndex"`
	OrderIndex      uint64 `json:"orderIndex"`
	AgreementIndex  uint64 `json:"agreementIndex"`
	RenewOrderIndex uint64 `json:"renewOrderIndex"`
	Working         string `json:"working"`
	Price           uint64 `json:"price"`
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configPath: DefaultConfigPath(),
	}
}

func NewConfigManagerWithPath(path string) *ConfigManager {
	return &ConfigManager{
		configPath: path,
	}
}

func DefaultConfigPath() string {
	return strings.Join([]string{DefaultConfigDir(), CONFIG_DEFAULT_FILENAME}, string(os.PathSeparator))
}

func DefaultConfigDir() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return CONFIG_DIR_NAME + "."
	}
	if err != nil {
		log.Error(err)
	}
	dir := strings.Join([]string{userHomeDir, CONFIG_DIR_NAME}, string(os.PathSeparator))
	if err != nil {
		log.Error(err)
	}
	return dir
}

func (cm *ConfigManager) GetConfig() (*Config, error) {
	packageLock.Lock()

	var cfg Config
	f, err := os.Open(cm.configPath)
	defer f.Close()

	if err != nil {
		cfg = Config{
			ApiPort:      8888,
			SeedOrPhrase: "tomato denial broccoli video correct spring link oval ostrich category stereo beauty",
			ChainApi:     "ws://127.0.0.1:9944",
			PublicIp:     "127.0.0.1",
			PublicPort:   4001,
		}
		packageLock.Unlock()
		err = cm.Save(&cfg)
		return &cfg, err
	}

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failure to decode config: %s", err)
	}

	defer packageLock.Unlock()
	return &cfg, nil
}

func (cm *ConfigManager) Save(config *Config) error {
	packageLock.Lock()
	defer packageLock.Unlock()
	f, err := os.OpenFile(cm.configPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0766)

	if err != nil {
		log.Println("hamster-gateway not initialized, run `hamster-gateway init`")
		err := os.MkdirAll(filepath.Dir(cm.configPath), os.ModeDir)
		if err != nil {
			log.Error(err)
		}

		err = os.Chmod(filepath.Dir(cm.configPath), os.ModePerm)
		if err != nil {
			log.Error(err)
		}
		f, err = os.OpenFile(cm.configPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0766)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		return err
	}
	return nil
}

func CreateIdentity() (Identity, error) {
	ident := Identity{
		SwarmKey: SWARM_KEY,
	}

	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return ident, err
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}
