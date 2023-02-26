package discover

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/core/proc"
)

const (
	allEths  = "0.0.0.0"
	envPodIP = "POD_IP"
)

type EtcdConf struct {
	discov.EtcdConf
	Metadata string
}

func (c EtcdConf) HasMetadata() bool {
	return c.Metadata != ""
}

func RegisterRpc(conf EtcdConf, ListenOn string) error {
	var pubOpts []discov.PubOption
	if conf.HasAccount() {
		pubOpts = append(pubOpts, discov.WithPubEtcdAccount(conf.User, conf.Pass))
	}
	if conf.HasTLS() {
		pubOpts = append(pubOpts, discov.WithPubEtcdTLS(conf.CertFile, conf.CertKeyFile,
			conf.CACertFile, conf.InsecureSkipVerify))
	}
	var metadata url.Values
	if conf.HasMetadata() {
		var err error
		metadata, err = url.ParseQuery(conf.Metadata)
		if err != nil {
			return err
		}
	}

	pubListenOn := figureOutListenOn(ListenOn)
	value := fmt.Sprintf("%s?%s", pubListenOn, metadata.Encode())
	pubClient := discov.NewPublisher(conf.Hosts, conf.Key, value, pubOpts...)
	proc.AddShutdownListener(func() {
		pubClient.Stop()
	})

	return pubClient.KeepAlive()
}

func MustRegisterRpc(conf EtcdConf, ListenOn string) {
	logx.Must(RegisterRpc(conf, ListenOn))
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
