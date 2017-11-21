package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/gman0/goenc"
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"
)

func main() {
	fPort := flag.Int("port", 8888, "port number")
	flag.Parse()

	kp := enc.GenerateKeyPair()
	srv := p2p.New(*fPort, func(c net.Conn) {
		fmt.Println("# New peer", c.RemoteAddr().String())
		p := p2p.NewPeer(c)

		if err := goenc.HandleClientConnection(p, kp); err != nil {
			panic(err)
		}

		c.Close()
	})

	go func() {
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("goenc ")
	for scanner.Scan() {
		parseInput(scanner.Text(), srv, kp)
		fmt.Print("goenc ")
	}

	srv.Shutdown()
}

func parseInput(in string, srv *p2p.Service, kp *enc.KeyPair) {
	if in == "" {
		return
	}

	cmd := strings.Split(in, " ")
	if len(cmd) < 1 {
		fmt.Println("# Input error")
		return
	}

	args := cmd[1:]

	switch cmd[0] {
	case "send":
		cmdSendFile(args, srv, kp)
	case "accept":
		cmdAccept(args)
	case "decline":
		cmdDecline(args)
	case "fingerprint":
		cmdFingerprint(kp)
	default:
		fmt.Println("# Command not found")
	}
}

func cmdSendFile(cmd []string, srv *p2p.Service, kp *enc.KeyPair) {
	if len(cmd) != 2 {
		fmt.Println("# Args := [ADDRESS:PORT] [FILEPATH]")
		return
	}

	addr := strings.Split(cmd[0], ":")
	if len(addr) != 2 {
		fmt.Println("# Address := [ADDRESS:PORT]")
		return
	}

	port, _ := strconv.ParseInt(addr[1], 10, 32)

	peerConf := p2p.ClientConfig{
		Address: addr[0],
		Port:    int(port),
		Handler: func(c net.Conn) {
			fmt.Println("# Connected to peer", c.RemoteAddr().String())
			p := p2p.NewPeer(c)

			if err := goenc.HandleServerConnection(cmd[1], p, kp); err != nil {
				panic(err)
			}

			c.Close()
		},
	}

	if err := srv.AddPeer(&peerConf); err != nil {
		fmt.Println("# Service error:", err)
		return
	}
}

func overwriteFile() bool {
	var answer string
	for answer != "n" && answer != "y" {
		fmt.Print("File already exists. Overwrite? [y/n]: ")
		fmt.Scanf("%s", &answer)
	}

	if answer == "y" {
		return true
	}

	return false
}

func cmdAccept(cmd []string) {
	if len(cmd) != 2 {
		fmt.Println("# Args := [REQUEST NUMBER] [DESTINATION]")
		return
	}

	id, err := strconv.Atoi(cmd[0])
	if err != nil {
		fmt.Println("# Error while parsing the request number")
		return
	}

	ar, err := goenc.GetAwaitingRequest(id)
	if err != nil {
		fmt.Println("#", err)
		return
	}

	dst := cmd[1]
	st, err := os.Stat(dst)
	if err == nil {
		if st.IsDir() {
			dst += "/" + ar.Req.Name
		} else {
			if !overwriteFile() {
				fmt.Println("# Won't overwrite, choose another file name")
				return
			}
		}
	} else {
		if !os.IsNotExist(err) {
			fmt.Println("# stat() failed:", err)
			return
		}
	}

	ar.RespChan <- goenc.Response{Accept: true, Dest: dst}
}

func cmdDecline(cmd []string) {
	if len(cmd) != 1 {
		fmt.Println("# Args := [REQUEST NUMBER]")
		return
	}
}

func cmdFingerprint(kp *enc.KeyPair) {
	goenc.PrintFingerprint(kp.Public)
}
