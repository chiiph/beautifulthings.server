package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"beautifulthings/account"
	"beautifulthings/server"
)

// remote <server>
// signup <username> <password>
// signin <username> <password>
// signout
// set <date> <input>
// enumerate <from> <to>
// echo <string>

const (
	NoRemote = iota
	Remoted
	Tokened
)

const (
	Unknown   = "UNKNOWN"
	Remote    = "remote"
	SignUp    = "signup"
	SignIn    = "signin"
	SignOut   = "signout"
	Set       = "set"
	Enumerate = "enumerate"
	Echo      = "echo"
)

func parseCommand(line string) (string, []string) {
	parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(parts) == 0 {
		return Unknown, nil
	}

	c := parts[0]
	switch c {
	case Echo:
		fallthrough
	case Remote:
		if len(parts) != 2 {
			return Unknown, nil
		}
		return c, parts[1:]
	case SignIn:
		fallthrough
	case Set:
		fallthrough
	case Enumerate:
		fallthrough
	case SignUp:
		args := strings.SplitN(parts[1], " ", 2)
		if len(args) != 2 {
			return Unknown, nil
		}
		return c, args
	case SignOut:
		return c, nil
	default:
		return Unknown, nil
	}
}

func commands(ps1 string, reader *bufio.Reader, print bool) (string, []string) {
	fmt.Print(ps1)
	text, err := reader.ReadString('\n')
	if err == io.EOF {
		if text == "" {
			return Unknown, nil
		}
	}
	if print {
		fmt.Print(text)
		if text[len(text)-1] != '\n' {
			fmt.Println()
		}
	}
	return parseCommand(text)
}

func getInput(ps1 string) (string, []string) {
	reader := bufio.NewReader(os.Stdin)
	return commands(ps1, reader, false)
}

func printErr(msg string, err error) {
	fmt.Printf("ERROR %s", msg)
	if err != nil {
		fmt.Printf(": %s", err)
	}
	fmt.Println()
}

func main() {
	var a *account.Account
	var s server.Server
	var err error

	inputFunc := getInput
	handleError := printErr
	if len(os.Args) > 1 {
		handleError = func(msg string, err error) {
			printErr(msg, err)
			os.Exit(1)
		}
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(f)

		inputFunc = func(ps1 string) (string, []string) {
			return commands(ps1, reader, true)
		}
	}

	state := NoRemote
	token := ""
	ps1 := "> "
	remote := ""
	user := ""

	for {
		cmd, args := inputFunc(ps1)
		switch cmd {
		case Remote:
			if state != NoRemote && state != Remoted {
				handleError("Can't set remote without signing out", nil)
				continue
			}
			remote = args[0]
			s = server.NewRemoteRest(remote)
			ps1 = fmt.Sprintf("s=%s > ", remote)
			state = Remoted
		case SignIn:
			if state != Remoted {
				handleError("Can't set remote without setting the remote or signing out", nil)
				continue
			}
			a, err = account.New(args[0], args[1])
			if err != nil {
				handleError("Creating server", err)
				continue
			}
			b, err := a.Bytes()
			if err != nil {
				handleError("Serializing account", err)
				continue
			}
			traw, err := s.SignIn(b)
			if err != nil {
				handleError("Signing in", err)
				continue
			}
			token, err = a.Decrypt(traw)
			if err != nil {
				handleError("Decrypting token", err)
				continue
			}
			user = args[0]
			ps1 = fmt.Sprintf("u=%s s=%s > ", user, remote)
			state = Tokened
		case SignUp:
			if state != Remoted {
				handleError("Can't set remote without setting the remote or signing out", nil)
				continue
			}
			a, err = account.New(args[0], args[1])
			if err != nil {
				handleError("Creating server", err)
				continue
			}
			b, err := a.Bytes()
			if err != nil {
				handleError("Serializing account", err)
				continue
			}
			err = s.SignUp(b)
			if err != nil {
				handleError("Signing up", err)
				continue
			}
			traw, err := s.SignIn(b)
			if err != nil {
				handleError("Signing in", err)
				continue
			}
			token, err = a.Decrypt(traw)
			if err != nil {
				handleError("Decrypting token", err)
				continue
			}
			user = args[0]
			ps1 = fmt.Sprintf("u=%s s=%s > ", user, remote)
			state = Tokened
		case Set:
			if state != Tokened {
				handleError("Can't run set without being signed in", nil)
				continue
			}
			ct, err := a.Encrypt(args[1])
			if err != nil {
				handleError("Encrypting input", err)
				continue
			}
			err = s.Set(token, args[0], ct)
			if err != nil {
				handleError(fmt.Sprintf("Setting %s", args[0]), err)
				continue
			}
		case Enumerate:
			if state != Tokened {
				handleError("Can't run set without being signed in", nil)
				continue
			}
			items, err := s.Enumerate(token, args[0], args[1])
			if err != nil {
				handleError(fmt.Sprintf("Enumerating from:%s to:%s", args[0], args[1]), err)
				continue
			}

			for i, item := range items {
				m, err := a.Decrypt(item.Content)
				if err != nil {
					handleError("Decryptng entry", err)
					continue
				}
				fmt.Printf("%d:\n  Date: %s\n  Content: %s\n", i, item.Date, m)
			}
		case SignOut:
			if state != Tokened {
				handleError("Can't run set without being signed in", nil)
				continue
			}
			state = Remoted
			s = nil
			user = ""
			ps1 = "> "
			token = ""
		case Unknown:
			handleError("Unknown command", nil)
		}
	}
}
