package listener

import (
	"crypto/tls"
	"net"

	"github.com/Sirupsen/logrus"
)

func Init(addr string, tlsConfig *tls.Config) (ls []net.Listener, err error) {
	l, err := initTCPSocket(addr, tlsConfig)
	if err != nil {
		return nil, err
	}
	ls = append(ls, l)

	return
}

func initTCPSocket(addr string, tlsConfig *tls.Config) (l net.Listener, err error) {
	if tlsConfig == nil || tlsConfig.ClientAuth != tls.RequireAndVerifyClientCert {
		logrus.Warn("/!\\ DON'T BIND ON ANY IP ADDRESS WITHOUT setting -tlsverify IF YOU DON'T KNOW WHAT YOU'RE DOING /!\\")
	}
	if l, err = NewTCPSocket(addr, tlsConfig); err != nil {
		return nil, err
	}
	return
}
