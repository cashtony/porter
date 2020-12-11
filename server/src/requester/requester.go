package requester

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	// DesktopUserAgent 电脑端浏览器标识
	DesktopUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
	// MobileUserAgent 移动端浏览器标识
	// MobileUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"
	MobileUserAgent = "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Mobile Safari/537.36"
)

var (
	// UserAgent 浏览器标识
	UserAgent = MobileUserAgent
	// DefaultClient 默认 http 客户端
	DefaultClient = NewHTTPClient()
)

type (
	ContentTyper interface {
		ContentType() string
	}

	// ContentLengther Content-Length 接口
	ContentLengther interface {
		ContentLength() int64
	}

	Lener interface {
		Len() int
	}

	// Lener64 返回64-bit长度接口
	Lener64 interface {
		Len() int64
	}
)

type HTTPClient struct {
	http.Client
	transport *http.Transport
	https     bool
	UserAgent string
}

func NewHTTPClient() *HTTPClient {
	h := &HTTPClient{
		Client: http.Client{
			Timeout: 300 * time.Second,
		},
		UserAgent: UserAgent,
	}
	h.Client.Jar, _ = cookiejar.New(nil)

	return h
}

func (h *HTTPClient) lazyInit() {
	if h.transport == nil {
		h.transport = &http.Transport{
			// Proxy:       proxyFunc,
			// DialContext: dialContext,
			// Dial:        dial,
			// DialTLS:     h.dialTLSFunc(),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !h.https,
			},
			TLSHandshakeTimeout:   10 * time.Second,
			DisableKeepAlives:     false,
			DisableCompression:    false, // gzip
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		}
		h.Client.Transport = h.transport
	}
}

// SetUserAgent 设置 UserAgent 浏览器标识
func (h *HTTPClient) SetUserAgent(ua string) {
	h.UserAgent = ua
}

// SetProxy 设置代理
func (h *HTTPClient) SetProxy(proxyAddr string) {
	h.lazyInit()
	var u *url.URL
	host, port, err := net.SplitHostPort(proxyAddr)
	if err == nil {
		u = &url.URL{
			Host: net.JoinHostPort(host, port),
		}
		return
	}
	u, err = url.Parse(proxyAddr)
	if err != nil {
		h.transport.Proxy = http.ProxyFromEnvironment
		return
	}

	h.transport.Proxy = http.ProxyURL(u)
}

// SetCookiejar 设置 cookie
func (h *HTTPClient) SetCookiejar(jar http.CookieJar) {
	h.Client.Jar = jar
}

// ResetCookiejar 清空 cookie
func (h *HTTPClient) ResetCookiejar() {
	h.Jar, _ = cookiejar.New(nil)
}

// SetHTTPSecure 是否启用 https 安全检查, 默认不检查
func (h *HTTPClient) SetHTTPSecure(b bool) {
	h.https = b
	h.lazyInit()
	if b {
		h.transport.TLSClientConfig = nil
	} else {
		h.transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: !b,
		}
	}
}

// SetKeepAlive 设置 Keep-Alive
func (h *HTTPClient) SetKeepAlive(b bool) {
	h.lazyInit()
	h.transport.DisableKeepAlives = !b
}

// SetGzip 是否启用Gzip
func (h *HTTPClient) SetGzip(b bool) {
	h.lazyInit()
	h.transport.DisableCompression = !b
}

// SetResponseHeaderTimeout 设置目标服务器响应超时时间
func (h *HTTPClient) SetResponseHeaderTimeout(t time.Duration) {
	h.lazyInit()
	h.transport.ResponseHeaderTimeout = t
}

// SetTLSHandshakeTimeout 设置tls握手超时时间
func (h *HTTPClient) SetTLSHandshakeTimeout(t time.Duration) {
	h.lazyInit()
	h.transport.TLSHandshakeTimeout = t
}

// SetTimeout 设置 http 请求超时时间, 默认30s
func (h *HTTPClient) SetTimeout(t time.Duration) {
	h.Client.Timeout = t
}
