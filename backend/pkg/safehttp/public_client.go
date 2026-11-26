// Package safehttp 提供服务端访问不可信外部 URL 时的统一网络边界。
package safehttp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

const (
	publicHTTPDialTimeout           = 5 * time.Second
	publicHTTPKeepAlive             = 30 * time.Second
	publicHTTPTLSHandshakeTimeout   = 5 * time.Second
	publicHTTPResponseHeaderTimeout = 15 * time.Second
	publicHTTPIdleConnTimeout       = 90 * time.Second
	publicHTTPMaxRedirects          = 3
)

// publicNonInternetPrefixes 补充 Go 标准库仍可能判为全局单播、但不应由公网请求访问的特殊用途网段。
var publicNonInternetPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001::/32"),
	netip.MustParsePrefix("2001:2::/48"),
	netip.MustParsePrefix("2001:10::/28"),
	netip.MustParsePrefix("2001:20::/28"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("2002::/16"),
}

// ipResolver 描述每次校验和拨号前使用的 DNS 解析能力。
type ipResolver interface {
	LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
}

// contextDialer 描述通过校验后的 IP 的最终建连能力。
type contextDialer interface {
	DialContext(context.Context, string, string) (net.Conn, error)
}

// validatingTransport 在每次 RoundTrip 前校验 URL；重定向产生的新请求也无法绕过初始目标策略。
type validatingTransport struct {
	base           *http.Transport
	resolver       ipResolver
	privateTargets privateTargetPolicy
}

// RoundTrip 校验协议、用户信息、主机和全部 DNS 结果后再交给受控 Transport。
func (t *validatingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("外部 HTTP 请求不能为空")
	}
	if err := validatePublicURL(req.Context(), req.URL, t.resolver, t.privateTargets); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

// CloseIdleConnections 释放底层 Transport 的空闲连接。
func (t *validatingTransport) CloseIdleConnections() {
	t.base.CloseIdleConnections()
}

// NewPublicClient 创建仅允许公网目标的 HTTP 客户端。
// 客户端禁用环境代理，在实际拨号前重新校验并固定 IP，同时只允许同源且不降级的有限重定向。
func NewPublicClient(timeout time.Duration) *http.Client {
	dialer := &net.Dialer{Timeout: publicHTTPDialTimeout, KeepAlive: publicHTTPKeepAlive}
	return newPublicClientWithNetwork(timeout, net.DefaultResolver, dialer)
}

// newPublicClientWithNetwork 注入 DNS 和拨号能力，确保 DNS 重绑定防线可以在无外网测试中确定性验证。
func newPublicClientWithNetwork(timeout time.Duration, resolver ipResolver, dialer contextDialer, policies ...privateTargetPolicy) *http.Client {
	var policy privateTargetPolicy
	if len(policies) > 0 {
		policy = policies[0]
	}
	transport := &http.Transport{
		// 代理会在进程外再次解析目标，无法保证这里验证的 IP 就是最终连接 IP，因此必须禁用。
		Proxy: nil,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialPublicAddress(ctx, network, address, resolver, dialer, policy)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          128,
		MaxIdleConnsPerHost:   16,
		MaxConnsPerHost:       32,
		IdleConnTimeout:       publicHTTPIdleConnTimeout,
		TLSHandshakeTimeout:   publicHTTPTLSHandshakeTimeout,
		ResponseHeaderTimeout: publicHTTPResponseHeaderTimeout,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
	}
	client := &http.Client{
		Timeout: timeout,
		Transport: &validatingTransport{
			base:           transport,
			resolver:       resolver,
			privateTargets: policy,
		},
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= publicHTTPMaxRedirects {
			return fmt.Errorf("外部 HTTP 重定向次数超过 %d", publicHTTPMaxRedirects)
		}
		if len(via) == 0 || !sameRedirectAuthority(via[0].URL, req.URL) {
			return fmt.Errorf("外部 HTTP 请求不允许跨源重定向")
		}
		if strings.EqualFold(via[0].URL.Scheme, "https") && !strings.EqualFold(req.URL.Scheme, "https") {
			return fmt.Errorf("外部 HTTP 请求不允许 HTTPS 降级")
		}
		return validatePublicURL(req.Context(), req.URL, resolver, policy)
	}
	return client
}

// dialPublicAddress 在实际连接使用的 DNS 结果上再次执行公网校验，并直接拨号到已验证 IP。
func dialPublicAddress(
	ctx context.Context,
	network string,
	address string,
	resolver ipResolver,
	dialer contextDialer,
	policies ...privateTargetPolicy,
) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("解析外部 HTTP 地址失败: %w", err)
	}
	ips, err := resolvePublicIPs(ctx, host, resolver, policies...)
	if err != nil {
		return nil, err
	}
	dialErrors := make([]error, 0, len(ips))
	for _, ip := range ips {
		connection, dialErr := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
		if dialErr == nil {
			return connection, nil
		}
		dialErrors = append(dialErrors, dialErr)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
	}
	return nil, fmt.Errorf("连接外部 HTTP 主机失败: %w", errors.Join(dialErrors...))
}

// validatePublicURL 校验不可信 URL 的协议、凭据、主机和全部解析结果。
func validatePublicURL(ctx context.Context, target *url.URL, resolver ipResolver, policies ...privateTargetPolicy) error {
	if target == nil || (target.Scheme != "http" && target.Scheme != "https") {
		return fmt.Errorf("外部 HTTP 仅支持 HTTP/HTTPS 地址")
	}
	if target.User != nil {
		return fmt.Errorf("外部 HTTP URL 不允许包含用户信息")
	}
	host := strings.TrimSpace(target.Hostname())
	if host == "" {
		return fmt.Errorf("外部 HTTP URL 缺少主机名")
	}
	_, err := resolvePublicIPs(ctx, host, resolver, policies...)
	return err
}

// resolvePublicIPs 解析并校验全部地址；任一 DNS 记录受限时整体拒绝，避免轮询命中内网地址。
func resolvePublicIPs(ctx context.Context, host string, resolver ipResolver, policies ...privateTargetPolicy) ([]net.IP, error) {
	host = strings.Trim(strings.TrimSpace(host), "[]")
	if strings.EqualFold(host, "localhost") {
		return nil, fmt.Errorf("不允许访问本地、私网或特殊用途地址")
	}
	if ip := net.ParseIP(host); ip != nil {
		if isForbiddenPublicIP(ip) && !allowsPrivateTarget(host, ip, policies...) {
			return nil, fmt.Errorf("不允许访问本地、私网或特殊用途地址")
		}
		return []net.IP{ip}, nil
	}
	if resolver == nil {
		return nil, fmt.Errorf("外部 HTTP DNS 解析器未配置")
	}
	addresses, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("解析外部 HTTP 主机失败: %w", err)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("外部 HTTP 主机没有可用地址")
	}
	ips := make([]net.IP, 0, len(addresses))
	for _, address := range addresses {
		if isForbiddenPublicIP(address.IP) && !allowsPrivateTarget(host, address.IP, policies...) {
			return nil, fmt.Errorf("外部 HTTP 主机解析到不允许的地址")
		}
		ips = append(ips, address.IP)
	}
	return ips, nil
}

// isForbiddenPublicIP 判断地址是否不属于可直接访问的公网目标。
func isForbiddenPublicIP(ip net.IP) bool {
	if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		return true
	}
	address = address.Unmap()
	if !address.IsGlobalUnicast() || address.IsPrivate() {
		return true
	}
	for _, prefix := range publicNonInternetPrefixes {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}

// sameRedirectAuthority 只允许原始请求同主机同端口跳转，防止携带 OAuth 凭据跨源泄漏。
func sameRedirectAuthority(source *url.URL, target *url.URL) bool {
	return source != nil && target != nil && strings.EqualFold(source.Host, target.Host)
}
