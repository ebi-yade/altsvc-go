package altsvc

import (
	"reflect"
	"testing"
)

// TODO: add further cases of both success and errors
func TestParse(t *testing.T) {
	tTable := []struct {
		input    string
		expected []Service
	}{
		{
			input: `clear`,
			expected: []Service{
				{Clear: true},
			},
		},
		{
			input: `h2=":443"; ma=2592000;`,
			expected: []Service{
				{ProtocolID: "h2", AltAuthority: AltAuthority{Port: "443"}, MaxAge: 2592000},
			},
		},
		{
			input: `h2=":443"; ma=2592000; persist=1`,
			expected: []Service{
				{ProtocolID: "h2", AltAuthority: AltAuthority{Port: "443"}, MaxAge: 2592000, Persist: 1},
			},
		},
		{
			input: `h2="alt.example.com:443", h2=":443"`,
			expected: []Service{
				{ProtocolID: "h2", AltAuthority: AltAuthority{Host: "alt.example.com", Port: "443"}},
				{ProtocolID: "h2", AltAuthority: AltAuthority{Port: "443"}},
			},
		},
		{
			input: `h3-25=":443"; ma=3600, h2=":443"; ma=3600`,
			expected: []Service{
				{ProtocolID: "h3-25", AltAuthority: AltAuthority{Port: "443"}, MaxAge: 3600},
				{ProtocolID: "h2", AltAuthority: AltAuthority{Port: "443"}, MaxAge: 3600},
			},
		},
	}

	for _, tCase := range tTable {
		svc, err := Parse(tCase.input)
		if err != nil {
			t.Errorf("failed to parse %s: %v\n", tCase.input, err)
		}
		if !reflect.DeepEqual(tCase.expected, svc) {
			t.Errorf(`expected "%v" but the result was %v\n`, tCase.expected, svc)
		}
	}
}
