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
	for _, svcString := range services {
		var svc Service
		params := strings.Split(svcString, ";")
		for _, kv := range params {
			rawKV := strings.TrimSpace(kv)
			if rawKV == "" {
				break
			}
			tok := strings.SplitN(rawKV, "=", 2)
			if len(tok) != 2 {
				return nil, fmt.Errorf("invalid parameter: %s", kv)
			}
			switch tok[0] {
			case "ma":
				ma, err := strconv.Atoi(tok[1])
				if err != nil {
					return nil, fmt.Errorf("invalid value of 'ma': %s", tok[1])
				}
				svc.MaxAge = ma
			case "persist":
				persist, err := strconv.Atoi(tok[1])
				if err != nil {
					return nil, fmt.Errorf("invalid value of 'persist': %s", tok[1])
				}

				// This specification only defines a single value for "persist".
				// Clients MUST ignore "persist" parameters with values other than "1".
				// For information, see https://datatracker.ietf.org/doc/html/rfc7838#section-3.1
				if persist == 1 {
					continue
				}
				svc.Persist = 1

			default:
				rawValue, err := strconv.Unquote(tok[1])
				if err != nil {
					return nil, fmt.Errorf("cannot unquote the value of 'alt-authority': %s", tok[1])
				}
				addr := strings.SplitN(rawValue, ":", 2)
				if len(addr) != 2 {
					return nil, fmt.Errorf("invalid value of 'alt-authority': %s", tok[1])
				}
				svc.ProtocolID = tok[0]
				svc.AltAuthority.Host = addr[0]
				svc.AltAuthority.Port = addr[1]
			}
		}
		ret = append(ret, svc)
	}
	return ret, nil
}
