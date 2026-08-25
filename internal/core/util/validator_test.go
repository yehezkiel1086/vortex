package util

import "testing"

func TestValidator(t *testing.T) {
	t.Run("IP validation", func(t *testing.T) {
		if !IsValidIPv4("185.10.20.30") {
			t.Errorf("expected 185.10.20.30 to be valid IPv4")
		}
		if IsValidIPv4("999.10.20.30") {
			t.Errorf("expected 999.10.20.30 to be invalid IPv4")
		}
		if !IsValidIPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334") {
			t.Errorf("expected IPv6 to be valid")
		}
		if !IsPrivateOrLoopbackIP("127.0.0.1") {
			t.Errorf("expected 127.0.0.1 to be loopback")
		}
		if !IsPrivateOrLoopbackIP("192.168.1.1") {
			t.Errorf("expected 192.168.1.1 to be private")
		}
		if IsPrivateOrLoopbackIP("8.8.8.8") {
			t.Errorf("expected 8.8.8.8 to be public")
		}
	})

	t.Run("Domain validation", func(t *testing.T) {
		if !IsValidDomain("evil-domain.com") {
			t.Errorf("expected evil-domain.com to be valid domain")
		}
		if !IsValidDomain("sub.c2.attacker.org") {
			t.Errorf("expected sub.c2.attacker.org to be valid domain")
		}
		if IsValidDomain("185.10.20.30") {
			t.Errorf("expected IP to not be valid domain")
		}
		if IsValidDomain("invalid_domain") {
			t.Errorf("expected invalid_domain to be invalid")
		}
	})

	t.Run("URL validation", func(t *testing.T) {
		if !IsValidURL("https://evil-domain.com/payload.exe") {
			t.Errorf("expected https url to be valid")
		}
		if !IsValidURL("http://185.10.20.30:8080/malware") {
			t.Errorf("expected http url to be valid")
		}
		if IsValidURL("ftp://evil.com/file") {
			t.Errorf("expected ftp to not be valid http/https url")
		}
	})

	t.Run("Hash validation", func(t *testing.T) {
		sha256Hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		if !IsValidSHA256(sha256Hash) {
			t.Errorf("expected valid sha256")
		}
		if IsValidSHA256("short-hash") {
			t.Errorf("expected invalid sha256")
		}

		md5Hash := "d41d8cd98f00b204e9800998ecf8427e"
		if !IsValidMD5(md5Hash) {
			t.Errorf("expected valid md5")
		}
	})
}
