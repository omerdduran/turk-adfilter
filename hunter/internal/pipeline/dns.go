package pipeline

import (
	"context"
	"net"
	"time"
)

// DNSResult, bir domain'in DNS doğrulama sonucu.
type DNSResult struct {
	Active bool   // yalnız "real" IP + quorum sağlanınca true
	IP     string // ilk gerçek (yoksa ilk herhangi) IP
	Class  string // real | park | sinkhole | block | ""
}

// lookuper, tek bir upstream resolver soyutlaması (test için mock'lanabilir).
type lookuper interface {
	lookupIP(ctx context.Context, host string) ([]net.IP, error)
}

// Resolver, bir domain'i doğrular.
type Resolver interface {
	Resolve(ctx context.Context, domain string) DNSResult
}

// QuorumResolver, birden çok public resolver'a sorup en az `need` tanesinden
// GERÇEK IP gelirse domaini aktif sayar. Registrar park / sinkhole / BTK-blok
// cevaplarını ayırt eder — tek resolver'a güvenmenin false-positive selini keser.
type QuorumResolver struct {
	resolvers []lookuper
	need      int
}

// NewQuorumResolver, sunucu adreslerinden bir resolver kurar. Özel değer
// "system" sistem resolver'ını kullanır (bazı ortamlarda özel UDP:53 bloklu olur;
// ancak sistem resolver DNS-tabanlı engelleyici barındırıyorsa yanıltıcı olabilir).
func NewQuorumResolver(servers []string, timeout time.Duration) *QuorumResolver {
	rs := make([]lookuper, 0, len(servers))
	for _, s := range servers {
		if s == "system" {
			rs = append(rs, &systemResolver{timeout: timeout})
		} else {
			rs = append(rs, newUDPResolver(s, timeout))
		}
	}
	// Majority quorum (2-of-3): tek resolver bloklanır/yanılırsa keşif sessizce
	// sıfırlanmasın, ama tek resolver da yeterli olmasın (false-positive koruması).
	need := len(rs)/2 + 1
	if need < 1 {
		need = 1
	}
	return &QuorumResolver{resolvers: rs, need: need}
}

// Resolve, quorum mantığıyla domaini doğrular.
func (q *QuorumResolver) Resolve(ctx context.Context, domain string) DNSResult {
	realCount := 0
	var firstReal, firstAny, anyClass string
	for _, r := range q.resolvers {
		ips, err := r.lookupIP(ctx, domain)
		if err != nil || len(ips) == 0 {
			continue
		}
		gotReal := false
		for _, ip := range ips {
			cls := classifyIP(ip)
			if firstAny == "" {
				firstAny, anyClass = ip.String(), cls
			}
			if cls == ClassReal && !gotReal {
				gotReal = true
				if firstReal == "" {
					firstReal = ip.String()
				}
			}
		}
		if gotReal {
			realCount++
		}
	}
	if realCount >= q.need {
		return DNSResult{Active: true, IP: firstReal, Class: ClassReal}
	}
	return DNSResult{Active: false, IP: firstAny, Class: anyClass}
}

// IP sınıfları.
const (
	ClassReal     = "real"
	ClassSinkhole = "sinkhole"
	ClassPark     = "park"
	ClassBlock    = "block"
	ClassInternal = "internal" // özel/dahili ağ — SSRF riski, aktif SAYILMAZ
)

// cgnatNet, RFC 6598 carrier-grade NAT bloğu (100.64.0.0/10).
var cgnatNet = mustCIDR("100.64.0.0/10")

func mustCIDR(s string) *net.IPNet {
	_, n, _ := net.ParseCIDR(s)
	return n
}

// isInternalIP, SSRF açısından tehlikeli (dahili/özel/link-local/CGNAT) IP'leri
// tespit eder. Bir saldırgan domain register edip bunlardan birine çözerek botu
// bulut metadata / iç servislere yönlendiremesin diye.
func isInternalIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		cgnatNet.Contains(ip)
}

// btkBlockIPs, BTK/ISS engel uyarı sayfası IP'leri. Bunlara çözülme "domain
// gerçek ve TR-hedefli" demektir ama içerik bloklu → aktif SAYILMAZ (bot Almanya
// sunucusunda; pratikte nadir, yine de güvenli sınıflandırma).
var btkBlockIPs = map[string]bool{
	"195.175.254.2":  true,
	"193.192.98.42":  true,
	"193.192.98.43":  true,
	"88.255.94.250":  true,
}

// parkIPs, bilinen registrar park/holding IP'leri (mirror için "henüz aktif değil").
var parkIPs = map[string]bool{
	"91.195.240.94":   true, // Sedo park
	"208.91.197.27":   true, // Above.com park
	"3.64.163.50":     true,
}

func classifyIP(ip net.IP) string {
	if ip == nil || ip.IsUnspecified() || ip.IsLoopback() {
		return ClassSinkhole
	}
	if isInternalIP(ip) {
		return ClassInternal // SSRF koruması: dahili IP asla "real" sayılmaz
	}
	s := ip.String()
	if btkBlockIPs[s] {
		return ClassBlock
	}
	if parkIPs[s] {
		return ClassPark
	}
	return ClassReal
}

// udpResolver, sabit bir upstream'e UDP ile soran net.Resolver sarmalı.
type udpResolver struct{ r *net.Resolver }

func newUDPResolver(server string, timeout time.Duration) *udpResolver {
	d := &net.Dialer{Timeout: timeout}
	return &udpResolver{r: &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			return d.DialContext(ctx, "udp", server)
		},
	}}
}

func (u *udpResolver) lookupIP(ctx context.Context, host string) ([]net.IP, error) {
	return ipAddrs(u.r.LookupIPAddr(ctx, host))
}

// systemResolver, işletim sisteminin varsayılan resolver'ını kullanır.
type systemResolver struct{ timeout time.Duration }

func (s *systemResolver) lookupIP(ctx context.Context, host string) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return ipAddrs(net.DefaultResolver.LookupIPAddr(ctx, host))
}

func ipAddrs(addrs []net.IPAddr, err error) ([]net.IP, error) {
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, 0, len(addrs))
	for _, a := range addrs {
		ips = append(ips, a.IP)
	}
	return ips, nil
}
