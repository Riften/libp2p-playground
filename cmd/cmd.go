package cmd

import (
	"context"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"main/host"
	"main/repo"
	"os"
)

const localhost = "http://localhost:8080"
type cmdsMap map[string]func() error
func Run() error {
	appCmd := kingpin.New("p2p",
		"p2p is a experimental toolbox of libp2p.")
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
	startRepoPath := startCmd.Arg("repo", "The path of repository.\n" +
		"The current path would be used as repository if not specified.\n").String()
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
		node, err := host.NewNode(context.Background(), cfg)
		if err != nil {
			fmt.Println("Error when create node: ", err)
			return err
		}
		node.Start()
		return nil
	}

	cmd := kingpin.MustParse(appCmd.Parse(os.Args[1:]))
	for key, value := range cmds {
		if key == cmd {
			return value()
		}
	}

	return nil
}