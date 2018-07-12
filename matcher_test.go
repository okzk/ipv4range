package ipv4range

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"testing"
)

func ExampleMatcher_Match() {
	matcher, _ := NewMatcher("10.10.0.0/16", "10.20.0.0/16", "10.30.0.0/16")

	fmt.Println(matcher.Match("10.10.10.10"))
	fmt.Println(matcher.Match("10.100.100.100"))

	// Output:
	// true
	// false
}

func TestMatcher_Match(t *testing.T) {
	matcher, err := NewMatcher("10.10.0.0/16", "10.10.10.0/24", "10.20.1.1", "192.168.1.10 - 192.168.1.20")
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]bool{
		"10.9.255.255": false, "10.10.0.0": true, "10.10.255.255": true, "10.11.0.0": false,
		"10.20.1.0": false, "10.20.1.1": true, "10.20.1.2": false,
		"192.168.1.9": false, "192.168.1.10": true, "192.168.1.20": true, "192.168.1.21": false,
	}

	for ip, expected := range data {
		actual := matcher.Match(ip)
		if actual != expected {
			t.Errorf("%s match failed. expected: %v, actual, %v", ip, expected, actual)
		}
	}
}

func generateRandomCIDRs(n int) []string {
	ret := make([]string, 0, n)
	for i := 0; i < n; i++ {
		v := rand.Int() & ((1 << 32) - 1)
		m := 32 - rand.Intn(16)
		v = v & ((1 << uint(m)) - 1)

		ip := make([]byte, 4)
		binary.BigEndian.PutUint32(ip, uint32(v))
		ip[0] = 10
		ret = append(ret, fmt.Sprintf("%v/%d", net.IP(ip), m))
	}

	return ret
}

func BenchmarkMatcher_Match(b *testing.B) {
	m, _ := NewMatcher(generateRandomCIDRs(1000)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Match("10.10.10.10")
	}
}
