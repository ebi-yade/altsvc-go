package altsvc

import (
	"fmt"
	"strconv"
	"strings"
)

// Service represents HTTP Alternative Services declared in RFC 7838
type Service struct {
	Clear        bool   // if true, it means the original value is the literal "clear"
	ProtocolID   string // ALPN protocol name
	AltAuthority AltAuthority

	// See https://datatracker.ietf.org/doc/html/rfc7838#section-3.1
	MaxAge  int
	Persist int
}

// AltAuthority is the set of Host and Port
type AltAuthority struct {
	Host string // an empty string means the alternative service is placed at the same host
	Port string
}

// Parse parses a string value and returns a slice of Services when no errors are raised
func Parse(s string) ([]Service, error) {
	ret := make([]Service, 0)
	if s == "clear" {
		ret = append(ret, Service{Clear: true})
		return ret, nil
	}

	services := strings.Split(s, ",")
	for i, svcString := range services {
		var svc Service
		params := strings.Split(svcString, ";")
		for j, kv := range params {
			rawKV := strings.TrimSpace(kv)
			k, v, ok := strings.Cut(rawKV, "=")
			if !ok {
				if rawKV == "" && j > 0 && j == len(params)-1 && i == len(services)-1 {
					// Note: assume the only trailing ";" is legal if the ";" does not have a valid parameter at the back
					break
				}
				return nil, fmt.Errorf("invalid parameter: %s", kv)
			}
			switch k {
			case "ma":
				ma, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("invalid value of 'ma': %s", v)
				}
				svc.MaxAge = ma
			case "persist":
				persist, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("invalid value of 'persist': %s", v)
				}

				// This specification only defines a single value for "persist".
				// Clients MUST ignore "persist" parameters with values other than "1".
				// For information, see https://datatracker.ietf.org/doc/html/rfc7838#section-3.1
				if persist != 1 {
					continue
				}
				svc.Persist = 1

			default:
				rawValue, err := strconv.Unquote(v)
				if err != nil {
					return nil, fmt.Errorf("cannot unquote the value of 'alt-authority': %s", v)
				}
				h, p, ok := strings.Cut(rawValue, ":")
				if !ok {
					return nil, fmt.Errorf("invalid value of 'alt-authority': %s", v)
				}
				svc.ProtocolID = k
				svc.AltAuthority.Host = h
				svc.AltAuthority.Port = p
			}
		}
		ret = append(ret, svc)
	}
	return ret, nil
}
