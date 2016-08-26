package logger

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

// FormatFileSize returns a string representation of a file size in bytes.
func FormatFileSize(sizeBytes int) string {
	if sizeBytes >= 1<<30 {
		return strconv.Itoa(sizeBytes/(1<<30)) + "gB"
	} else if sizeBytes >= 1<<20 {
		return strconv.Itoa(sizeBytes/(1<<20)) + "mB"
	} else if sizeBytes >= 1<<10 {
		return strconv.Itoa(sizeBytes/(1<<10)) + "kB"
	}
	return strconv.Itoa(sizeBytes) + "B"
}

// GetIP gets the origin/client ip for a request.
// X-FORWARDED-FOR is checked. If multiple IPs are included the first one is returned
// X-REAL-IP is checked. If multiple IPs are included the first one is returned
// Finally r.RemoteAddr is used
// Only benevolent services will allow access to the real IP.
func GetIP(r *http.Request) string {
	tryHeader := func(key string) (string, bool) {
		if headerVal := r.Header.Get(key); len(headerVal) > 0 {
			if !strings.ContainsRune(headerVal, ',') {
				return headerVal, true
			}
			return strings.SplitN(headerVal, ",", 2)[0], true
		}
		return "", false
	}

	for _, header := range []string{"X-FORWARDED-FOR", "X-REAL-IP"} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
