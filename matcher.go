package ipv4range

import (
	"encoding/binary"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
)

type ipv4Range struct {
	from uint32
	to   uint32
}

var sep = regexp.MustCompile("\\s*-\\s*")

func newIPv4Range(s string) (ipv4Range, error) {
	if strings.Contains(s, "/") {
		_, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			return ipv4Range{}, err
		}
		if len(ipNet.IP) != net.IPv4len {
			return ipv4Range{}, fmt.Errorf("invalid IPv4 CIDR: %s", s)
		}
		from := binary.BigEndian.Uint32(ipNet.IP)
		mask := binary.BigEndian.Uint32(ipNet.Mask)
		to := from | ((^uint32(0)) ^ mask)
		return ipv4Range{from: from, to: to}, nil
	}

	part := sep.Split(s, 2)
	ip := net.ParseIP(part[0]).To4()
	if ip == nil {
		return ipv4Range{}, fmt.Errorf("invalid IPv4: %s", s)
	}
	from := binary.BigEndian.Uint32(ip)
	if len(part) == 1 {
		return ipv4Range{from: from, to: from}, nil
	}
	ip = net.ParseIP(part[1]).To4()
	if ip == nil {
		return ipv4Range{}, fmt.Errorf("invalid IPv4 range: %s", s)
	}
	to := binary.BigEndian.Uint32(ip)
	return ipv4Range{from: from, to: to}, nil
}

// Matcher provides fast IPv4 address matching
type Matcher struct {
	list []ipv4Range
}

// NewMatcher creates a Matcher with the ip ranges.
// IP range expression should be one of following formats:
//   - CIDR(ex. "10.0.0.0/8")
//   - raw IP(ex. "127.0.0.1")
//   - IP range(ex. "192.168.1.1 - 192.168.1.100")
func NewMatcher(rangeList ...string) (*Matcher, error) {
	parsed := make([]ipv4Range, 0, len(rangeList))
	for _, s := range rangeList {
		r, err := newIPv4Range(s)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, r)
	}
	sort.Slice(parsed, func(i, j int) bool {
		if parsed[i].from == parsed[j].from {
			return parsed[i].to < parsed[j].to
		}
		return parsed[i].from < parsed[j].from
	})

	merged := make([]ipv4Range, 0, len(rangeList))
	for _, r := range parsed {
		last := len(merged) - 1
		if last < 0 {
			merged = append(merged, r)
			continue
		}
		if merged[last].to+1 < r.from {
			merged = append(merged, r)
			continue
		}
		if merged[last].to < r.to {
			merged[last].to = r.to
		}
	}
	return &Matcher{list: merged}, nil
}

// Match reports whether the Matcher's ip ranges includes ip
func (m *Matcher) Match(ipStr string) bool {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return false
	}
	v := binary.BigEndian.Uint32(ip)
	list := m.list
	for len(list) != 0 {
		n := len(list) / 2
		if list[n].from > v {
			list = list[0:n]
		} else if list[n].to < v {
			list = list[n+1:]
		} else {
			return true
		}
	}
	return false
}
