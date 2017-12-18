package sshhelper

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var output io.Writer

func init() {
	output = colorable.NewColorableStdout()
}

type SSHClient struct {
	conn *ssh.Client
}

func NewSSHClient() *SSHClient {
	return &SSHClient{conn: nil}
}

func (client *SSHClient) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
}

func (client *SSHClient) Connect(addr, user, passwd string) error {
	if client.conn != nil {
		return nil
	}
	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

func (client *SSHClient) RunCommandFile(commandFile string) error {
	session, err := client.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	termWidth, termHeight, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		return err
	}

	w, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer w.Close()

	session.Stdout = output
	session.Stderr = output

	if err := session.Shell(); err != nil {
		return err
	}
	defer session.Wait()

	in := make(chan string)
	errchan := make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if exitErr, ok := err.(error); ok {
					errchan <- exitErr
				}
			}
		}()
		time.Sleep(200 * time.Millisecond)
		for cmd := range in {
			n, err := w.Write([]byte(cmd + "\n"))
			if err != nil {
				panic(err)
			} else if n == 0 {
				panic(errors.New("connect failed"))
			}
		}
	}()

	f, err := os.Open(commandFile)
	if err != nil {
		return err
	}
	defer f.Close()

	commands := bufio.NewReader(f)
	for {
		command, err := commands.ReadString('\n')
		command = strings.TrimSpace(command)
		if command != "" && !strings.HasPrefix(command, "//") {
			select {
			case in <- command:
			case err := <-errchan:
				close(in)
				return err
			}
		}

		if err != nil {
			in <- "exit"
			close(in)
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}

func (client *SSHClient) StartShell() error {
	session, err := client.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = output
	session.Stderr = output
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		return err
	}
	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}
	return nil
}
