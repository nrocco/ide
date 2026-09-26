package tools

import (
	"fmt"
	"iter"
	"math/big"
	"math/bits"
	"net/netip"
	"strconv"
	"strings"
)

// Special purpose ranges, matching the tables used by Python's ipaddress module
var (
	ipv4PrivateNetworks = mustParsePrefixes(
		"0.0.0.0/8", "10.0.0.0/8", "127.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12",
		"192.0.0.0/24", "192.0.0.170/31", "192.0.2.0/24", "192.168.0.0/16", "198.18.0.0/15",
		"198.51.100.0/24", "203.0.113.0/24", "240.0.0.0/4", "255.255.255.255/32",
	)
	ipv4PrivateExceptions  = mustParsePrefixes("192.0.0.9/32", "192.0.0.10/32")
	ipv4SharedAddressSpace = netip.MustParsePrefix("100.64.0.0/10")
	ipv4Multicast          = netip.MustParsePrefix("224.0.0.0/4")
	ipv4Loopback           = netip.MustParsePrefix("127.0.0.0/8")
	ipv4LinkLocal          = netip.MustParsePrefix("169.254.0.0/16")
	ipv4Reserved           = netip.MustParsePrefix("240.0.0.0/4")
	ipv4Unspecified        = netip.MustParsePrefix("0.0.0.0/32")

	ipv6PrivateNetworks = mustParsePrefixes(
		"::1/128", "::/128", "::ffff:0:0/96", "64:ff9b:1::/48", "100::/64", "2001::/23",
		"2001:db8::/32", "2002::/16", "3fff::/20", "fc00::/7", "fe80::/10",
	)
	ipv6PrivateExceptions = mustParsePrefixes(
		"2001:1::1/128", "2001:1::2/128", "2001:3::/32", "2001:4:112::/48", "2001:20::/28", "2001:30::/28",
	)
	ipv6ReservedNetworks = mustParsePrefixes(
		"::/8", "100::/8", "200::/7", "400::/6", "800::/5", "1000::/4", "4000::/3", "6000::/3",
		"8000::/3", "a000::/3", "c000::/3", "e000::/4", "f000::/5", "f800::/6", "fe00::/9",
	)
	ipv6Multicast   = netip.MustParsePrefix("ff00::/8")
	ipv6Loopback    = netip.MustParsePrefix("::1/128")
	ipv6LinkLocal   = netip.MustParsePrefix("fe80::/10")
	ipv6SiteLocal   = netip.MustParsePrefix("fec0::/10")
	ipv6Unspecified = netip.MustParsePrefix("::/128")
	ipv6Mapped      = netip.MustParsePrefix("::ffff:0:0/96")
	ipv6SixToFour   = netip.MustParsePrefix("2002::/16")
	ipv6Teredo      = netip.MustParsePrefix("2001::/32")
)

func mustParsePrefixes(prefixes ...string) []netip.Prefix {
	result := make([]netip.Prefix, len(prefixes))
	for i, p := range prefixes {
		result[i] = netip.MustParsePrefix(p)
	}
	return result
}

// IPNetwork is an IPv4 or IPv6 network in CIDR notation. It remembers the
// address it was created from, host bits included.
type IPNetwork struct {
	addr   netip.Addr
	prefix netip.Prefix
}

// NewIPNetwork parses an address with an optional prefix length, e.g.
// 192.168.1.10/24 or 2001:db8::/112. IPv4 networks can also be given with a
// netmask or hostmask, e.g. 10.0.0.0/255.0.0.0. Without a prefix length the
// network holds a single address.
func NewIPNetwork(cidr string) (IPNetwork, error) {
	addrPart, lenPart, hasLen := strings.Cut(strings.TrimSpace(cidr), "/")

	addr, err := netip.ParseAddr(addrPart)
	if err != nil {
		return IPNetwork{}, fmt.Errorf("%q does not appear to be an IPv4 or IPv6 address", addrPart)
	}
	if addr.Zone() != "" {
		return IPNetwork{}, fmt.Errorf("zone identifiers are not supported: %q", addrPart)
	}

	prefixLen := addr.BitLen()
	if hasLen {
		prefixLen, err = parsePrefixLen(lenPart, addr)
		if err != nil {
			return IPNetwork{}, err
		}
	}

	return IPNetwork{addr: addr, prefix: netip.PrefixFrom(addr, prefixLen).Masked()}, nil
}

func newIPNetwork(addr netip.Addr, prefixLen int) IPNetwork {
	prefix := netip.PrefixFrom(addr, prefixLen).Masked()
	return IPNetwork{addr: prefix.Addr(), prefix: prefix}
}

// String returns the network in CIDR notation, without host bits
func (n IPNetwork) String() string {
	return n.prefix.String()
}

// Version returns 4 or 6
func (n IPNetwork) Version() int {
	return n.Address().Version()
}

// Address returns the address the network was created from
func (n IPNetwork) Address() IPAddress {
	return IPAddress{n.addr}
}

// PrefixLength returns the number of network bits
func (n IPNetwork) PrefixLength() int {
	return n.prefix.Bits()
}

// Netmask returns the network mask, e.g. 255.255.255.0
func (n IPNetwork) Netmask() IPAddress {
	return IPAddress{maskAddr(n.prefix.Bits(), n.maxBits(), false)}
}

// Hostmask returns the inverted network mask (wildcard), e.g. 0.0.0.255
func (n IPNetwork) Hostmask() IPAddress {
	return IPAddress{maskAddr(n.prefix.Bits(), n.maxBits(), true)}
}

// NetworkAddress returns the first address of the network
func (n IPNetwork) NetworkAddress() IPAddress {
	return IPAddress{n.prefix.Addr()}
}

// LastAddress returns the last address of the network
func (n IPNetwork) LastAddress() IPAddress {
	last, _ := intToAddr(new(big.Int).Sub(n.offset(n.NumAddresses()), big.NewInt(1)), n.maxBits())
	return IPAddress{last}
}

// Broadcast returns the broadcast address. Only IPv4 networks larger than /31
// have one.
func (n IPNetwork) Broadcast() (IPAddress, bool) {
	if n.Version() != 4 || n.prefix.Bits() >= 31 {
		return IPAddress{}, false
	}
	return n.LastAddress(), true
}

// FirstUsable returns the first address that can be assigned to a host
func (n IPNetwork) FirstUsable() IPAddress {
	first, _, _, _ := n.usable()
	return first
}

// LastUsable returns the last address that can be assigned to a host
func (n IPNetwork) LastUsable() IPAddress {
	_, last, _, _ := n.usable()
	return last
}

// NumAddresses returns the number of addresses in the network
func (n IPNetwork) NumAddresses() *big.Int {
	return pow2(n.maxBits() - n.prefix.Bits())
}

// NumUsable returns the number of addresses that can be assigned to hosts
func (n IPNetwork) NumUsable() *big.Int {
	_, _, count, _ := n.usable()
	return count
}

// UsableNote explains which addresses are excluded from the usable range
func (n IPNetwork) UsableNote() string {
	_, _, _, note := n.usable()
	return note
}

func (n IPNetwork) usable() (IPAddress, IPAddress, *big.Int, string) {
	first, last, total := n.NetworkAddress(), n.LastAddress(), n.NumAddresses()
	if n.Version() == 4 {
		if n.prefix.Bits() >= 31 { // /31 point-to-point (RFC 3021), /32 single host
			return first, last, total, "RFC 3021 (/31) or single host (/32)"
		}
		return IPAddress{first.addr.Next()}, IPAddress{last.addr.Prev()}, new(big.Int).Sub(total, big.NewInt(2)), "network and broadcast excluded"
	}
	// IPv6 has no broadcast. First address is the Subnet-Router anycast (RFC 4291 2.6.1).
	if n.prefix.Bits() >= 127 { // /127 point-to-point (RFC 6164), /128 single host
		return first, last, total, "RFC 6164 (/127) or single host (/128)"
	}
	return IPAddress{first.addr.Next()}, last, new(big.Int).Sub(total, big.NewInt(1)), "Subnet-Router anycast (first address) excluded, RFC 4291 2.6.1"
}

// Flags returns the special purpose properties that apply to the whole network,
// e.g. private, global, multicast or loopback
func (n IPNetwork) Flags() []string {
	first, last := n.prefix.Addr(), n.LastAddress().addr
	both := func(is func(netip.Addr) bool) bool { return is(first) && is(last) }

	type flag struct {
		name  string
		value bool
	}
	var flags []flag
	if n.Version() == 4 {
		private := isPrivateNetwork(first, last, ipv4PrivateNetworks, ipv4PrivateExceptions)
		flags = []flag{
			{"private", private},
			{"global", !(ipv4SharedAddressSpace.Contains(first) && ipv4SharedAddressSpace.Contains(last)) && !private},
		}
	} else {
		private := isPrivateNetwork(first, last, ipv6PrivateNetworks, ipv6PrivateExceptions)
		flags = []flag{
			{"private", private},
			{"global", !private},
		}
	}
	flags = append(flags,
		flag{"multicast", both(inRanges(ipv4Multicast, ipv6Multicast))},
		flag{"loopback", both(inRanges(ipv4Loopback, ipv6Loopback))},
		flag{"link_local", both(inRanges(ipv4LinkLocal, ipv6LinkLocal))},
		flag{"reserved", both(inRanges(ipv4Reserved, ipv6ReservedNetworks...))},
		flag{"unspecified", both(inRanges(ipv4Unspecified, ipv6Unspecified))},
	)
	if n.Version() == 6 {
		flags = append(flags, flag{"site_local", ipv6SiteLocal.Contains(first) && ipv6SiteLocal.Contains(last)})
	}

	result := []string{}
	for _, f := range flags {
		if f.value {
			result = append(result, f.name)
		}
	}
	return result
}

// Previous returns the network of the same size right before this one
func (n IPNetwork) Previous() (IPNetwork, bool) {
	return n.neighbour(new(big.Int).Neg(n.NumAddresses()))
}

// Next returns the network of the same size right after this one
func (n IPNetwork) Next() (IPNetwork, bool) {
	return n.neighbour(n.NumAddresses())
}

func (n IPNetwork) neighbour(delta *big.Int) (IPNetwork, bool) {
	addr, ok := intToAddr(n.offset(delta), n.maxBits())
	if !ok {
		return IPNetwork{}, false
	}
	return newIPNetwork(addr, n.prefix.Bits()), true
}

// Supernet returns the network that is one bit larger. A /0 has no supernet.
func (n IPNetwork) Supernet() (IPNetwork, bool) {
	if n.prefix.Bits() == 0 {
		return IPNetwork{}, false
	}
	return newIPNetwork(n.prefix.Addr(), n.prefix.Bits()-1), true
}

// Contains reports if other lies completely within this network
func (n IPNetwork) Contains(other IPNetwork) bool {
	return other.prefix.Bits() >= n.prefix.Bits() && n.prefix.Contains(other.prefix.Addr())
}

// Overlaps reports if this network and other have any address in common
func (n IPNetwork) Overlaps(other IPNetwork) bool {
	return n.prefix.Overlaps(other.prefix)
}

// Subnets splits the network into subnets of the given prefix length. The
// number of subnets can be huge, so they are generated lazily.
func (n IPNetwork) Subnets(prefixLen int) (iter.Seq[IPNetwork], error) {
	if prefixLen < n.prefix.Bits() || prefixLen > n.maxBits() {
		return nil, fmt.Errorf("prefix length %d must be between %d and %d", prefixLen, n.prefix.Bits(), n.maxBits())
	}

	step := pow2(n.maxBits() - prefixLen)
	return func(yield func(IPNetwork) bool) {
		start, end := addrToInt(n.prefix.Addr()), n.offset(n.NumAddresses())
		for i := start; i.Cmp(end) < 0; i = new(big.Int).Add(i, step) {
			addr, _ := intToAddr(i, n.maxBits())
			if !yield(newIPNetwork(addr, prefixLen)) {
				return
			}
		}
	}, nil
}

// offset returns the network address plus delta as an integer
func (n IPNetwork) offset(delta *big.Int) *big.Int {
	return new(big.Int).Add(addrToInt(n.prefix.Addr()), delta)
}

func (n IPNetwork) maxBits() int {
	return n.addr.BitLen()
}

// IPAddress is a single IPv4 or IPv6 address
type IPAddress struct {
	addr netip.Addr
}

// String returns the address in its canonical notation, IPv6 compressed
func (a IPAddress) String() string {
	return a.addr.String()
}

// Version returns 4 or 6
func (a IPAddress) Version() int {
	if a.addr.Is4() {
		return 4
	}
	return 6
}

// Hex returns the address as a hexadecimal number, e.g. 0xffffff00
func (a IPAddress) Hex() string {
	return fmt.Sprintf("0x%0*x", a.addr.BitLen()/4, addrToInt(a.addr))
}

// Binary returns the address in binary notation, per octet for IPv4 and per
// group for IPv6
func (a IPAddress) Binary() string {
	b := a.addr.AsSlice()
	if a.addr.Is4() {
		octets := make([]string, len(b))
		for i := range b {
			octets[i] = fmt.Sprintf("%08b", b[i])
		}
		return strings.Join(octets, ".")
	}
	groups := make([]string, len(b)/2)
	for i := range groups {
		groups[i] = fmt.Sprintf("%08b%08b", b[2*i], b[2*i+1])
	}
	return strings.Join(groups, ":")
}

// Exploded returns the IPv6 address without zero compression. IPv4-mapped
// addresses keep their dotted IPv4 notation, like Python's ipaddress does.
// IPv4 addresses are returned as is.
func (a IPAddress) Exploded() string {
	if a.addr.Is4() {
		return a.addr.String()
	}
	b := a.addr.As16()
	groups := make([]string, 8)
	for i := range groups {
		groups[i] = fmt.Sprintf("%02x%02x", b[2*i], b[2*i+1])
	}
	if a.addr.Is4In6() {
		return strings.Join(groups[:6], ":") + ":" + a.addr.Unmap().String()
	}
	return strings.Join(groups, ":")
}

// IPv4Mapped returns the IPv4 address of an IPv4-mapped IPv6 address (::ffff:0:0/96)
func (a IPAddress) IPv4Mapped() (IPAddress, bool) {
	if !ipv6Mapped.Contains(a.addr) {
		return IPAddress{}, false
	}
	return IPAddress{a.addr.Unmap()}, true
}

// SixToFour returns the IPv4 address embedded in a 6to4 address (2002::/16)
func (a IPAddress) SixToFour() (IPAddress, bool) {
	if !ipv6SixToFour.Contains(a.addr) {
		return IPAddress{}, false
	}
	b := a.addr.As16()
	return IPAddress{netip.AddrFrom4([4]byte(b[2:6]))}, true
}

// Teredo returns the server and client IPv4 addresses embedded in a Teredo
// address (2001::/32)
func (a IPAddress) Teredo() (server, client IPAddress, ok bool) {
	if !ipv6Teredo.Contains(a.addr) {
		return IPAddress{}, IPAddress{}, false
	}
	b := a.addr.As16()
	server = IPAddress{netip.AddrFrom4([4]byte(b[4:8]))}
	client = IPAddress{netip.AddrFrom4([4]byte{^b[12], ^b[13], ^b[14], ^b[15]})}
	return server, client, true
}

func parsePrefixLen(text string, addr netip.Addr) (int, error) {
	if text != "" && strings.Trim(text, "0123456789") == "" {
		prefixLen, err := strconv.Atoi(text)
		if err != nil || prefixLen > addr.BitLen() {
			return 0, fmt.Errorf("invalid prefix length %q", text)
		}
		return prefixLen, nil
	}

	if addr.Is4() {
		if mask, err := netip.ParseAddr(text); err == nil && mask.Is4() {
			m := uint32(addrToInt(mask).Uint64())
			if ones := bits.LeadingZeros32(^m); m == ^uint32(0)<<(32-ones) {
				return ones, nil
			}
			if zeros := bits.LeadingZeros32(m); m == ^uint32(0)>>zeros {
				return zeros, nil
			}
		}
		return 0, fmt.Errorf("invalid prefix length, netmask or hostmask %q", text)
	}

	return 0, fmt.Errorf("invalid prefix length %q", text)
}

// isPrivateNetwork reports if the range first-last falls within a single private
// network and does not touch any of the exceptions
func isPrivateNetwork(first, last netip.Addr, networks, exceptions []netip.Prefix) bool {
	for _, p := range exceptions {
		if p.Contains(first) || p.Contains(last) {
			return false
		}
	}
	for _, p := range networks {
		if p.Contains(first) && p.Contains(last) {
			return true
		}
	}
	return false
}

// inRanges returns a predicate that checks an address against the IPv4 or IPv6
// ranges. IPv4-mapped IPv6 addresses are checked as their IPv4 address.
func inRanges(ipv4 netip.Prefix, ipv6 ...netip.Prefix) func(netip.Addr) bool {
	return func(addr netip.Addr) bool {
		if addr.Is4() || addr.Is4In6() {
			return ipv4.Contains(addr.Unmap())
		}
		for _, p := range ipv6 {
			if p.Contains(addr) {
				return true
			}
		}
		return false
	}
}

func maskAddr(prefixLen, maxBits int, host bool) netip.Addr {
	mask := new(big.Int).Sub(pow2(maxBits), pow2(maxBits-prefixLen))
	if host {
		mask.Sub(pow2(maxBits-prefixLen), big.NewInt(1))
	}
	addr, _ := intToAddr(mask, maxBits)
	return addr
}

func pow2(n int) *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), uint(n))
}

func addrToInt(addr netip.Addr) *big.Int {
	return new(big.Int).SetBytes(addr.AsSlice())
}

func intToAddr(n *big.Int, maxBits int) (netip.Addr, bool) {
	if n.Sign() < 0 || n.BitLen() > maxBits {
		return netip.Addr{}, false
	}
	buf := make([]byte, maxBits/8)
	n.FillBytes(buf)
	addr, _ := netip.AddrFromSlice(buf)
	return addr, true
}
