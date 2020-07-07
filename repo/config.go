package repo

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"io/ioutil"
	"os"
	"path"
)

const defaultPort = 40101
const defaultApiPort = 40102

type Config struct {
	Pubkey []byte
	PrivKey []byte
	Port int
	ApiPort int
}

func InitConfig() (*Config, error) {
	r:= rand.Reader
	privK, pubK, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, r)
	//privK, pubK, err :=
	if err != nil {
		fmt.Printf("Error occur when generate key pair\n%s\n", err.Error())
		return nil, err
	}
	privK_m, err := crypto.MarshalPrivateKey(privK)
	if err != nil {
		fmt.Printf("Error occur when marshal private key\n%s\n", err.Error())
		return nil, err
	}
	pubK_m, err := crypto.MarshalPublicKey(pubK)
	if err != nil {
		fmt.Printf("Error occur when marshal public key\n%s\n", err.Error())
		return nil, err
	}
	return &Config{
		PrivKey:privK_m,
		Pubkey:pubK_m,
		Port:defaultPort,
		ApiPort:defaultApiPort,
	}, nil
}

// Read reads config from disk
func Read(repoPath string) (*Config, error) {

	data, err := ioutil.ReadFile(path.Join(repoPath, "config"))
	if err != nil {
		return nil, err
	}

	var conf *Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// Write replaces the on-disk version of config with the given one
func Write(repoPath string, conf *Config) error {
	f, err := os.Create(path.Join(repoPath, "config"))
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
