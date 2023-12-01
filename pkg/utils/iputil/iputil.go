package iputil

import (
	"net"
	"net/http"
)

// Define http headers.
const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	XClientIP     = "x-client-ip"
)

// 本地非回环 IP 地址是指在本地网络中不指向回环地址（loopback address）的 IP 地址。回环地址是一种特殊的 IP 地址，用于在同一设备上模拟网络通信，通常是指 127.0.0.1（IPv4）或 ::1（IPv6）。它们被用于本地主机内部进行网络通信，数据不会离开主机。
// 本地非回环 IP 地址是指在本地网络中用于识别设备或与其他设备进行通信的 IP 地址。在局域网或广域网环境中，非回环 IP 地址可用于设备之间的通信，这些地址与其他设备的地址相对应，可以用于数据包的传输和接收

// GetLocalIP  获取本地非回环 IP
func GetLocalIP() string {

	// 返回所有网络接口的地址信息。
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	// 遍历这些地址，检查地址类型是否不是回环地址（!ipnet.IP.IsLoopback()），并且是 IPv4 地址（ipnet.IP.To4() != nil）。
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// RemoteIP 用于获取 HTTP 请求的远程 IP 地址。它尝试从请求头中获取可能包含客户端 IP 地址的字段，然后选择一个作为远程 IP 地址返回。
func RemoteIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XClientIP); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
