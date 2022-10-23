package pkg

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

func WaitForExit(srv *dns.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	_ = srv.Shutdown()
}

func PrintHeader() string {
	return "\n><((((> dns c2 by syl <))))><\n"
}

func GetIp(remoteAddr string) string {
	ip, _, _ := net.SplitHostPort(remoteAddr)
	return ip
}
