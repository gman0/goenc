package goenc

import (
	"fmt"
	"github.com/gman0/goenc/enc"
	"github.com/gman0/goenc/p2p"
)

func PrintRequestInfo(id int, r *Request, p *p2p.Peer) {
	fmt.Printf("\n> %s #%d: '%s' %d bytes (%s) ; ", p.Conn.RemoteAddr().String(), id, r.Name, r.Size, r.Type)
	PrintFingerprint(p.Public)
}

func PrintFingerprint(key enc.RSAPublicKey) {
	fp, err := key.Fingerprint()
	if err != nil {
		fmt.Println("[ RSAPublicKey error:", err, "]")
		return
	}

	for i := range fp {
		fmt.Printf("%02x", fp[i])
		if i != len(fp)-1 {
			fmt.Print(":")
		}
	}
	fmt.Println()
}
