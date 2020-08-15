package repo

import (
	"crypto/rand"
	"fmt"
	"github.com/Riften/libp2p-playground/util"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/libp2p/go-libp2p-core/crypto"
	"log"
	"github.com/ipfs/go-ipfs/plugin/loader"
)

type ErrRepoExists struct {
	path string
}

func (e *ErrRepoExists) Error() string {
	return fmt.Sprintf("repo %s not empty")
}

func LoadPlugins(repoPath string) (*loader.PluginLoader, error) {
	// check if repo is accessible before loading plugins
	_, err := util.CheckPermissions(repoPath)
	if err != nil {
		return nil, err
	}

	plugins, err := loader.NewPluginLoader(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error loading preloaded plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return nil, nil
	}

	if err := plugins.Inject(); err != nil {
		return nil, nil
	}
	return plugins, nil
}

func InitIpfsRepo(repoPath string) error {
	// create repo if not exists:
	err := util.CreateDirIfNotExist(repoPath)
	if err != nil {
		log.Println("Error when check repo directory: ", err)
		return err
	}

	if fsrepo.IsInitialized(repoPath) {
		return &ErrRepoExists{path: repoPath}
	}

	_, err = LoadPlugins(repoPath)
	if err != nil {
		log.Println("Error when load plugins: ", err)
	}

	// Generate Identity
	r:= rand.Reader
	privK, pubK, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, r)
	if err != nil {
		log.Println("Error when generate key pair: ", err)
		return err
	}
	ident, err := util.IdentityFromKey(privK, pubK)
	if err != nil {
		log.Println("Error when generate identity from key pair: ", err)
		return err
	}
	conf, err := config.InitWithIdentity(ident)
	if err != nil {
		log.Println("Error when create config: ", err)
		return err
	}

	_, err = LoadPlugins(repoPath)
	if err != nil {
		return err
	}

	err = fsrepo.Init(repoPath, conf)
	if err != nil {
		log.Println("Error when init ipfs repo: ", err)
		return err
	}

	return nil
}