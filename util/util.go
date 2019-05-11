package util

import (
	"encoding/base64"
	"errors"
	"log"
	"net"
	"regexp"
)

// Decode base64 string.
func Encode(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

// Encode string using base64.
func Decode(s string) string {
	bytes,_ := base64.URLEncoding.DecodeString(s)
	return string(bytes)
}

// IsDomain detect if value match the format of domain.
func IsDomain(host string) bool {
	match, err := regexp.MatchString(`\.[a-z]{2,}$`, host)
	if err != nil {
		return false
	}
	return match
}

// findIP returns the first ip matched detector.
func findIP(ips []net.IP, f func(ip net.IP) bool) net.IP {
	for _, ip := range ips {
		if f(ip) {
			return ip
		}
	}
	return nil
}

// Lookup return ip address of host.
func Lookup(host string) (net.IP, error) {
	addrs, err := net.LookupIP(host)
	if err != nil {
		log.Fatalln("error LookupHost")
	}
	ip := findIP(addrs, func(ip net.IP) bool {
		if ip != nil && ip.To4() != nil {
			return true
		}
		return false
	})
	if ip == nil {
		return ip, errors.New("not found ip")
	}
	return ip, nil
}
