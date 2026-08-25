package util

import (
	"net/netip"
	"net/url"
	"regexp"
	"strings"
)

var (
	sha256Regex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	sha1Regex   = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
	md5Regex    = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	domainRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	emailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func IsValidIP(ip string) bool {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	return err == nil && addr.IsValid()
}

func IsValidIPv4(ip string) bool {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	return err == nil && addr.IsValid() && addr.Is4()
}

func IsValidIPv6(ip string) bool {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	return err == nil && addr.IsValid() && addr.Is6()
}

func IsPrivateOrLoopbackIP(ip string) bool {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil || !addr.IsValid() {
		return false
	}
	return addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsUnspecified()
}

func IsValidDomain(domain string) bool {
	domain = strings.TrimSpace(strings.ToLower(domain))
	if domain == "" || len(domain) > 253 {
		return false
	}
	// Avoid classifying valid IPs as domains
	if IsValidIP(domain) {
		return false
	}
	return domainRegex.MatchString(domain)
}

func IsValidURL(rawURL string) bool {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return false
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func IsValidSHA256(hash string) bool {
	return sha256Regex.MatchString(strings.TrimSpace(hash))
}

func IsValidSHA1(hash string) bool {
	return sha1Regex.MatchString(strings.TrimSpace(hash))
}

func IsValidMD5(hash string) bool {
	return md5Regex.MatchString(strings.TrimSpace(hash))
}

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}
