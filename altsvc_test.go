package altsvc

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tTable := []struct {
		input    string
		expected []Service
	}{
		{
			input: `h2=":443"; ma=2592000;`,
			expected: []Service{
				{ProtocolID: "h2", AltAuthority: AltAuthority{Port: "443"}, MaxAge: 2592000},
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
