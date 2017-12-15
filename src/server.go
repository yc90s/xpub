package main

import (
	"bytes"
	"config"
	"errors"
	"fmt"
	"log"
	"sshhelper"
	"strings"
	"xterm"
)

type Server struct {
	config.Section
	addr        string
	commandsMap map[string]string
}

func NewServer(config config.Section) *Server {
	var cmdMap map[string]string
	cmdMap = make(map[string]string)
	for _, command := range config.Commands {
		cmdMap[command.Cmd] = command.Commandfile
	}

	var buf bytes.Buffer
	buf.WriteString(config.Host)
	buf.WriteString(":")
	buf.WriteString(config.Port)

	return &Server{Section: config, commandsMap: cmdMap, addr: buf.String()}
}

func (ctx *Server) Close() {
}

func (ctx *Server) Active(xTerm *xterm.XTerm) {
	oldPrompt := xTerm.SetPrompt("[" + ctx.Username + "@" + ctx.Name + " ~]>>> ")
	defer xTerm.SetPrompt(oldPrompt)

ForEnd:
	for {
		line, err := xTerm.ReadLine()
		if err != nil {
			panic(err)
		}

		cmd := strings.ToLower(line)
		cmd = strings.TrimSpace(cmd)
		switch cmd {
		case "help":
			serverUsage()
		case "quit":
			break ForEnd
		case "shell":
			if err := ctx.Shell(); err != nil {
				log.Println(err)
			}
		default:
			if err := ctx.Command(cmd); err != nil {
				log.Println(err)
			}
		}
	}
}

func (ctx *Server) Shell() error {
	sshclient := sshhelper.NewSSHClient()
	if err := sshclient.Connect(ctx.addr, ctx.Username, ctx.Passwd); err != nil {
		return errors.New("connect server failed")
	}
	defer sshclient.Close()
	return sshclient.StartShell()
}

func (ctx *Server) Command(cmd string) error {
	if len(cmd) == 0 {
		return nil
	}

	sshclient := sshhelper.NewSSHClient()
	if err := sshclient.Connect(ctx.addr, ctx.Username, ctx.Passwd); err != nil {
		return errors.New("connect server failed")
	}
	defer sshclient.Close()

	if file, ok := ctx.commandsMap[cmd]; ok {
		return sshclient.RunCommandFile(file)
	}
	//	return errors.New(cmd + " not exist.")
	return nil
}

func serverUsage() {
	str := `Usage:
	shell:	start an interactive shell with the remote host
	help :	display this help
	quit :	go back to the last one 	

Otherwise, You can enter the commands that you configure in your configuration file. Only this server will execute. 
`
	fmt.Println(str)
}
