package main

import (
	"github.com/facundobatista/go-telegram/logging"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"gowork/telegram-connector/telegram"
)

func show_incoming(origin, message string) {
	fmt.Printf("<---[%s] %q\n", origin, message)
}

func main() {
	// check parameters
	var verbose1 = flag.Bool("v", false, "Be verbose")
	var verbose2 = flag.Bool("vv", false, "Be very verbose")
	flag.Parse()
	//if len(flag.Args()) < 2 {
	//	log.Fatal("Usage: tgram [{-v|-vv}] <path-to-telegram-cli> <path-to-server.pub>")
	//}
	tgCliPath := "/snap/bin/telegram-cli"
	tgPubPath := "/home/ubuntu/snap/telegram-cli/server.pub"

	// convert verbose flags to log level
	var loglevel int
	if *verbose2 {
		loglevel = logging.LevelDebug
	} else if *verbose1 {
		loglevel = logging.LevelInfo
	} else {
		loglevel = logging.LevelError
	}

	// start Telegram backend
	fmt.Printf("Hello! Starting backend...\n")
	telegramCli, err := telegram.New(tgCliPath, tgPubPath, show_incoming, loglevel)
	if err != nil {
		log.Fatal(err)
	}

	// start dialog with user
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Done! Allowed: quit, send, list-contacts\n")

	// main user interface loop
	shouldQuit := false
	for 1 == 1 {
		fmt.Printf(">> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		text = strings.TrimSpace(text)
		tokens := strings.Split(text, " ")
		fmt.Printf("=== user: %q\n", tokens)
		switch tokens[0] {
		case "quit":
			shouldQuit = true
		case "list-contacts":
			contacts := telegramCli.ListContacts()
			for _, v := range contacts {
				fmt.Print(v + "\n")
			}
		case "send":
			dest := tokens[1]
			msg := strings.Join(tokens[2:], " ")
			telegramCli.SendMessage(dest, msg)
		}
		if shouldQuit {
			break
		}
	}

	// clean up and die
	fmt.Printf("Quitting\n")
	telegramCli.Quit()
}
