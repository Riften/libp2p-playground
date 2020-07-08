package cmd

import (
	"context"
	"fmt"
	"github.com/Riften/libp2p-playground/api"
	"github.com/Riften/libp2p-playground/host"
	"github.com/Riften/libp2p-playground/repo"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
)

const defaultApiPort = "7891"
var a *api.Api

type cmdsMap map[string]func() error
func Run() error {
	appCmd := kingpin.New("p2p",
		"p2p is a experimental toolbox of libp2p.")
	appApiPort := appCmd.Flag("api", "The port used to handle api.").Default(defaultApiPort).Int()

	cmds := make(cmdsMap)

	initCmd := appCmd.Command("init", "Initialize the repository.")
	initRepoPath := initCmd.Arg("repo", "The path of repository.\n" +
		"The current path would be used as repository if not specified.\n").String()

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

		config, err := repo.InitConfig()
		if err != nil {
			fmt.Println("Error when init config: ", err)
			return err
		}
		return repo.Write(repoPath, config)
	}

	startCmd := appCmd.Command("start", "Start the playground node.")
	startPort := startCmd.Arg("port", "Specified the port running playground.").Required().Int()
	startRepoPath := startCmd.Arg("repo", "The path of repository.\n" +
		"The current path would be used as repository if not specified.").String()
	cmds[startCmd.FullCommand()] = func () error {
		repoPath := *startRepoPath
		if repoPath == "" {
			fmt.Println("No repoPath specified. Used current directory as repo.")
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error when get pwd: ", err)
				return err
			}
			repoPath = pwd
		}
		cfg, err := repo.Read(repoPath)
		if err != nil {
			fmt.Println("Error when read config: ", err)
			return err
		}
		cfg.Port = *startPort
		node, err := host.NewNode(context.Background(), cfg)
		if err != nil {
			fmt.Println("Error when create node: ", err)
			return err
		}
		a = &api.Api{
			Node: node,
			Port: *appApiPort,
		}
		http.Handle("/", a)
		err = http.ListenAndServe(fmt.Sprintf(":%d", a.Port), nil)
		if err != nil {
			fmt.Printf("Failed to start api server: %v\n", err)
		}
		node.Start()

		return nil
	}

	exprCmd := appCmd.Command("expr", "Doing experiment.")
	exprStreamCmd := exprCmd.Command("stream", "Transport through stream.")
	cmds[exprStreamCmd.FullCommand()] = func() error {
		return api.SendRequest("expr", map[string]string{"a": "a1"}, *appApiPort)
	}

	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
	}

	return nil
}