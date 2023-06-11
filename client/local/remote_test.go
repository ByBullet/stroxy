package local

import "testing"

func TestConnectProxyServer(t *testing.T) {
	connectProxyServerTest("wss://sp.masterliu.fun:65533/test", "https://sp.masterliu.fun:65533")

	// connectProxyServerTest("wss://bybullet-curly-potato-7wggv996rjp3p7wv-65533.preview.app.github.dev/test", "https://bybullet-curly-potato-7wggv996rjp3p7wv-65533.preview.app.github.dev")
}
