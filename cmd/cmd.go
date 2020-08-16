package cmd

import (
	"context"
	"fmt"
	"github.com/Riften/libp2p-playground/api"
	"github.com/Riften/libp2p-playground/global"
	"github.com/Riften/libp2p-playground/host"
	"github.com/Riften/libp2p-playground/repo"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
)

const defaultApiPort = "7891"

type method string // e.g. http.MethodGet

type params struct {
	args    []string
	opts    map[string]string
	payload io.Reader
	ctype   string
}

type cmdsMap map[string]func() error
func Run() error {
	appCmd := kingpin.New("p2p",
		"p2p is a experimental toolbox of libp2p.")
//	appCopy := appCmd.Flag("copy", "Whether copy to clipboard.").Short('c').Bool()
//	appApiPort := appCmd.Flag("api", "The port used to handle api.").Default(defaultApiPort).Int()

	cmds := make(cmdsMap)

	// ======== init
	initCmd := appCmd.Command("init", "Initialize the repository.")
	initRepoPath := initCmd.Arg("repo", "The path of repository.\n" +
		"The current path would be used as repository if not specified.\n").String()
	initTransport := initCmd.Flag("transport", "The transport protocol").Default("tcp").Short('t').String()

	cmds[initCmd.FullCommand()] = func() error {
		repoPath := *initRepoPath

		if repoPath == "" {
			fmt.Println("No repoPath specified. Used current directory as repo.")
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error when get pwd: ", err)
				return err
			}
			repoPath = pwd
		}

		config, err := repo.InitConfig(*initTransport)
		if err != nil {
			fmt.Println("Error when init config: ", err)
			return err
		}

		return repo.Write(repoPath, config)
	}

	// ======== start
	startCmd := appCmd.Command("start", "Start the playground node.")
	startLibp2p := startCmd.Flag("libp2p", "Whether to start the libp2p host.").Bool()
	startTcp := startCmd.Flag("tcp", "Whether to start the tcp service.").Bool()

	startRepoPath := startCmd.Arg("repo", "The path of repository.\n" +
		"The current path would be used as repository if not specified.").String()

	cmds[startCmd.FullCommand()] = func () error {
		node := host.EmptyNode(context.Background())

		// start libp2p host
		if *startLibp2p {
			repoPath := *startRepoPath
			if repoPath == "" {
				log.Println("No repoPath specified. Used current directory as repo.")
				pwd, err := os.Getwd()
				if err != nil {
					log.Println("Error when get pwd: ", err)
					return err
				}
				repoPath = pwd
			}
			cfg, err := repo.Read(repoPath)
			if err != nil {
				log.Println("Error when read config: ", err)
				return err
			}

			err = node.StartLibp2p(cfg)
			if err != nil {
				log.Println("Error when start libp2p host: ", err)
				return err
			}
		}

		// start tcp service
		if *startTcp {
			err := node.StartTcp()
			if err != nil {
				log.Println("Error when start tcp service: ", err)
				return err
			}
		}

		// start the api server
		r := api.InitRouter(node)
		r.Run(api.ApiPort)
		return nil
	}

	// ======== peer
	peerCmd := appCmd.Command("peer", "libp2p peer related command")
	peerCopy := peerCmd.Flag("copy", "Whether copy to clipboard.").Short('c').Bool()
	peerInfoCmd := peerCmd.Command("info", "Get the info of host peer.")
	//peerInfoCopy := peerCmd.Flag("copy", "Whether copy to clipboard.").Short('c').Bool()
	peerInfoOut := peerInfoCmd.Arg("outFile", "The path of output file.").String()
	cmds[peerInfoCmd.FullCommand()] = func() error {
		return peerInfo(*peerCopy, *peerInfoOut)
	}

	peerListCmd := peerCmd.Command("list", "List all the connected peers.")
	//peerListCopy := peerCmd.Flag("copy", "Whether copy to clipboard.").Short('c').Bool()
	cmds[peerListCmd.FullCommand()] = func() error {
		return peerList(*peerCopy)
	}

	peerConnectCmd := peerCmd.Command("connect", "Connect with another peer.")
	peerConnectId := peerConnectCmd.Arg("id", "peer id.").Required().String()
	peerConnectAddr := peerConnectCmd.Arg("address", "peer address").Required().String()
	cmds[peerConnectCmd.FullCommand()] = func() error {
		return peerConnect(*peerConnectId, *peerConnectAddr)
	}

	// ======== speed
	speedCmd := appCmd.Command("speed", "Speed test related commands.")
	speedSendCmd := speedCmd.Command("send", "Start send task to remote peer.")
	speedSendPeer := speedSendCmd.Arg("peer", "Remote peer.").Required().String()
	cmds[speedSendCmd.FullCommand()] = func () error {
		return speedSend(*speedSendPeer)
	}

	// ======== ipfs
	ipfsCmd := appCmd.Command("ipfs", "Ipfs related commands.")
	ipfsInitCmd := ipfsCmd.Command("init", "Init ipfs peer.")
	ipfsInitRepo := ipfsInitCmd.Arg("repo", "Repo path.").String()
	cmds[ipfsInitCmd.FullCommand()] = func () error {
		repoPath := *ipfsInitRepo
		if repoPath == "" {
			fmt.Println("No repoPath specified. Used current directory as repo.")
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error when get pwd: ", err)
				return err
			}
			repoPath = pwd
		}
		err := repo.InitIpfsRepo(repoPath)
		if err != nil {
			log.Println("Error when init ipfs repo: ", err)
			return err
		}
		return nil
	}

	ipfsStartCmd := ipfsCmd.Command("start", "Start ipfs node.")
	ipfsStartRepo := ipfsStartCmd.Arg("repo", "Repo path.").String()
	cmds[ipfsStartCmd.FullCommand()] = func () error {
		repoPath := *ipfsStartRepo
		if repoPath == "" {
			fmt.Println("No repoPath specified. Used current directory as repo.")
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error when get pwd: ", err)
				return err
			}
			repoPath = pwd
		}
		var err error
		global.Ipfs, err = host.NewIpfsNode(context.Background(), repoPath)
		if err != nil {
			log.Println("Error when start ipfs node: ", err)
		}
		return nil
	}

	// =============== tcp
	tcpCmd := appCmd.Command("tcp", "tcp related commands.")
	// tcpPort := tcpCmd.Flag("port", "tcp port").Default("10000").Int()
	tcpListenCmd := tcpCmd.Command("listen", "start listen tcp port")
	tcpListenPort := tcpListenCmd.Arg("port", "listening port").Default("10000").Int()
	cmds[tcpListenCmd.FullCommand()] = func () error {
		return tcpListen(*tcpListenPort)
	}
	tcpSendCmd := tcpCmd.Command("send", "start sending tcp sockets")
	tcpSendIp := tcpSendCmd.Arg("ip", "ip address of receiver").Required().String()
	tcpSendPort := tcpSendCmd.Arg("port", "port of receiver").Default("10000").Int()
	cmds[tcpSendCmd.FullCommand()] = func () error {
		return tcpSend(*tcpSendIp, *tcpSendPort)
	}


	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
	}

	return nil
}