package conf

import (
	"testing"
)

func TestConfMarshalling(t *testing.T) {
	conf := &BirdConf{
		Filename: "bird.conf",
		Statements: []ConfStatement{
			&IncludeStatement{
				Filename: "bird.conf",
			},
			&Proto{
				Type: Bgp,
				Name: "bgp-1",
			},
		},
	}

	// Marshal the configuration to a string
	m, err := conf.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal configuration: %v", err)
	}

	t.Logf("%s", string(m))
}
