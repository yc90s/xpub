package main

import (
	"config"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"xterm"
)

type options struct {
	conf string
}

func init() {
	str := `───────────────────────────────────────────────────────────────────
─████████──████████─██████████████─██████──██████─██████████████───
─██░░░░██──██░░░░██─██░░░░░░░░░░██─██░░██──██░░██─██░░░░░░░░░░██───
─████░░██──██░░████─██░░██████░░██─██░░██──██░░██─██░░██████░░██───
───██░░░░██░░░░██───██░░██──██░░██─██░░██──██░░██─██░░██──██░░██───
───████░░░░░░████───██░░██████░░██─██░░██──██░░██─██░░██████░░████─
─────██░░░░░░██─────██░░░░░░░░░░██─██░░██──██░░██─██░░░░░░░░░░░░██─
───████░░░░░░████───██░░██████████─██░░██──██░░██─██░░████████░░██─
───██░░░░██░░░░██───██░░██─────────██░░██──██░░██─██░░██────██░░██─
─████░░██──██░░████─██░░██─────────██░░██████░░██─██░░████████░░██─
─██░░░░██──██░░░░██─██░░██─────────██░░░░░░░░░░██─██░░░░░░░░░░░░██─
─████████──████████─██████─────────██████████████─████████████████─
───────────────────────────────────────────────────────────────────
`
	fmt.Println(str)
}

func parseArgs(opt *options) {
	var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.StringVar(&opt.conf, "conf", "../config/server.json", "config file name")

	commandLine.Parse(os.Args[1:])
}

func main() {
	log.SetPrefix("[xpub:]")
	log.SetFlags(1)

	var opt options
	parseArgs(&opt)

	log.Println("start loading config...")
	if config.LoadConfig(opt.conf) == false {
		log.Fatalln("load config failed")
		return
	}
	log.Println("loding config complete")

	usage()

	ctxs := make([]*Server, 0, 0)
	for _, section := range config.Configurations {
		ctx := NewServer(section)
		defer ctx.Close()
		ctxs = append(ctxs, ctx)
	}

	xTerm := xterm.NewXTerm(">>> ", []string{"help", "list", "quit", "shell"})
	defer xTerm.Close()
ForEnd:
	for {
		line, err := xTerm.ReadLine()
		if err != nil {
			panic(err)
		}

		cmd := strings.ToLower(line)
		cmd = strings.TrimSpace(cmd)

		no, err := strconv.Atoi(cmd)
		if err == nil {
			if no >= 0 && no < len(ctxs) {
				ctxs[no].Active(xTerm)
			} else {
				log.Println("wrong server no.")
			}

			continue
		}

		switch cmd {
		case "help":
			usage()
		case "list":
			showAllServer(ctxs)
		case "quit":
			fmt.Println("bye.")
			break ForEnd
		default:
			if len(cmd) == 0 {
				continue
			}
			runCommand(cmd, ctxs)
		}
	}
}

func usage() {
	str := `Usage:
	list:	show all servers
	help:	display this help
	quit:	exit

Otherwise, You can enter the commands that you configure in your configuration file. 
Note that all servers containing this command will execute. And you can also enter 
a server num to connect to the server.
`
	fmt.Println(str)
}

func showAllServer(ctxs []*Server) {
	fmt.Printf("%-10s%-32s%-32s\n", "Num", "Name", "Host")
	for i, ctx := range ctxs {
		fmt.Printf("%-10d%-32s%-32s\n", i, ctx.Name, ctx.Host)
	}
}

func runCommand(cmd string, ctxs []*Server) {
	for _, ctx := range ctxs {
		fmt.Printf("========== %s ==========\n", ctx.Name)
		if err := ctx.Command(cmd); err != nil {
			log.Println(err)
		}
		fmt.Println()
	}
}
