package env

import "testing"

func TestEnv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		LoadEnv()
	})

	t.Run("get", func(t *testing.T) {
		if env := GetEnv(); env.Mode != "release" {
			t.Fatalf("expect: %v, got: %v", "release", env.Mode)
		}
	})
}
