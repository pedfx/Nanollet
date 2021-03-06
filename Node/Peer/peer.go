package Peer

import (
	"net"
	"strings"
	"errors"
	"strconv"
	"time"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
	"crypto/rand"
)

const (
	BetaPort = 54000
	LivePort = 7075
)

const (
	Timeout = 2 * time.Minute
)

var (
	ErrInvalidIP      = errors.New("invalid ip")
	ErrInvalidPort    = errors.New("invalid port")
	ErrIncompleteData = errors.New("invalid IP:PORT, both are needed")
)

type Peer struct {
	UDP *net.UDPAddr
	TCP *net.TCPAddr

	LastSeen  time.Time
	PublicKey Wallet.PublicKey
	Header    [8]byte
	Challenge [32]byte
}

func NewPeer(ip net.IP, port int) (peer *Peer) {
	peer = &Peer{
		UDP: &net.UDPAddr{
			IP:   ip,
			Port: port,
		},
		LastSeen: time.Now(),
	}

	rand.Read(peer.Challenge[:])

	return peer
}

func NewPeersFromString(hosts ...string) (peers []*Peer) {
	for _, host := range hosts {
		if ip, port, err := parseIP(host); err == nil {
			peers = append(peers, NewPeer(ip, port))
			continue
		}

		if ips, port, err := parseHost(host); err == nil {
			for _, ip := range ips {
				peers = append(peers, NewPeer(ip, port))
			}
		}
	}

	return peers
}

func (p *Peer) IsActive() bool {
	if p == nil {
		return false
	}

	if time.Since(p.LastSeen) > Timeout {
		return false
	}

	return true
}

func (p *Peer) IsKnow() bool {
	if p == nil {
		return false
	}

	if Util.IsEmpty(p.PublicKey[:]) {
		return false
	}

	return true
}

func (p *Peer) String() string {
	if p.UDP != nil {
		return p.UDP.String()
	}

	if p.TCP != nil {
		return p.TCP.String()
	}

	return ""
}

func parseHost(s string) (ips []net.IP, port int, err error) {
	split := strings.Split(s, ":")
	if len(split) < 2 {
		return ips, port, ErrIncompleteData
	}

	port, err = parsePort(split[1])
	if err != nil {
		return ips, port, err
	}

	ips, err = net.LookupIP(split[0])
	if err != nil {
		return ips, port, ErrInvalidIP
	}

	return ips, port, nil
}

func parseIP(s string) (ip net.IP, port int, err error) {
	split := strings.Split(s, ":")
	if len(split) < 2 {
		return ip, port, ErrIncompleteData
	}

	port, err = parsePort(split[1])
	if err != nil {
		return ip, port, err
	}

	ip = net.ParseIP(split[0])
	if ip == nil {
		return ip, port, ErrInvalidIP
	}

	return ip, port, nil
}

func parsePort(s string) (port int, err error) {
	port, err = strconv.Atoi(s)
	if err != nil || port > 65535 {
		return 0, ErrInvalidPort
	}

	return port, nil
}
