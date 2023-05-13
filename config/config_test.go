package config

import "testing"

func TestGroup_SetConfigArg(t *testing.T) {
	g := Group{
		LocalPort: 1234,
	}

	g.SetConfigArg("Port", "3456")
}
