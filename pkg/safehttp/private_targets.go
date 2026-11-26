package safehttp

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

// privateTargetPolicy 将部署配置中的精确域名绑定到私网 CIDR，不能由请求参数扩权。
type privateTargetPolicy map[string][]netip.Prefix

// NewPublicClientWithPrivateTargets 为可信部署配置增加私网例外，空配置保留公网客户端行为。
func NewPublicClientWithPrivateTargets(timeout time.Duration, targets map[string][]string) (*http.Client, error) {
	policy, err := parsePrivateTargets(targets)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{Timeout: publicHTTPDialTimeout, KeepAlive: publicHTTPKeepAlive}
	return newPublicClientWithNetwork(timeout, net.DefaultResolver, dialer, policy), nil
}

// parsePrivateTargets 拒绝通配域名和非私网范围，避免误配置扩大出站访问边界。
func parsePrivateTargets(targets map[string][]string) (privateTargetPolicy, error) {
	policy := privateTargetPolicy{}
	for host, cidrs := range targets {
		normalized := strings.ToLower(strings.TrimSuffix(host, "."))
		if !validPrivateTargetHost(normalized) || len(cidrs) == 0 {
			return nil, fmt.Errorf("私网目标需要精确域名和非空 CIDR")
		}
		if _, exists := policy[normalized]; exists {
			return nil, fmt.Errorf("私网目标域名重复")
		}
		for _, cidr := range cidrs {
			prefix, err := netip.ParsePrefix(cidr)
			if err != nil {
				return nil, fmt.Errorf("私网目标 CIDR 无效")
			}
			prefix = prefix.Masked()
			minimumBits := 7
			if prefix.Addr().Is4() {
				minimumBits = 8
				if netip.MustParsePrefix("172.16.0.0/12").Contains(prefix.Addr()) {
					minimumBits = 12
				}
				if netip.MustParsePrefix("192.168.0.0/16").Contains(prefix.Addr()) {
					minimumBits = 16
				}
			}
			if !prefix.Addr().IsPrivate() || prefix.Addr().Is4In6() || prefix.Bits() < minimumBits {
				return nil, fmt.Errorf("仅允许普通私网 CIDR")
			}
			policy[normalized] = append(policy[normalized], prefix)
		}
	}
	return policy, nil
}

// validPrivateTargetHost 只接受 DNS 主机名，禁止 URL、端口、IP 和通配符混入配置。
func validPrivateTargetHost(host string) bool {
	if host == "localhost" || len(host) > 253 || !strings.Contains(host, ".") || net.ParseIP(host) != nil {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, ch := range label {
			if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-') {
				return false
			}
		}
	}
	return true
}

// allowsPrivateTarget 在每次解析和实际拨号时同时匹配域名与私网地址，不放行本机或元数据地址。
func allowsPrivateTarget(host string, ip net.IP, policies ...privateTargetPolicy) bool {
	address, ok := netip.AddrFromSlice(ip)
	if !ok || len(policies) == 0 {
		return false
	}
	address = address.Unmap()
	if !address.IsPrivate() {
		return false
	}
	for _, prefix := range policies[0][strings.ToLower(strings.TrimSuffix(host, "."))] {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}
