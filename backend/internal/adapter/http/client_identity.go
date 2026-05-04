package httpadapter

import (
	"net"
	"net/http"
	"strings"
)

type ClientIdentity struct {
	IP        string
	UserAgent string
}

func newClientIdentity(r *http.Request, trustProxyHeaders bool) ClientIdentity {
	return ClientIdentity{
		IP:        clientIP(r, trustProxyHeaders),
		UserAgent: r.UserAgent(),
	}
}

func clientIP(r *http.Request, trustProxyHeaders bool) string {
	if trustProxyHeaders {
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			if ip := firstValidForwardedIP(forwardedFor); ip != "" {
				return ip
			}
		}

		if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(realIP) != nil {
			return realIP
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func firstValidForwardedIP(forwardedFor string) string {
	parts := strings.Split(forwardedFor, ",")
	for _, part := range parts {
		ip := strings.TrimSpace(part)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	return ""
}
