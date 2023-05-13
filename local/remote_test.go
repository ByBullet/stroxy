package local

import "testing"

func TestConnectProxyServer(t *testing.T) {
	connectProxyServerTest("wss://www.liuio.xyz:2000/test", "https://www.liuio.xyz:2000")
}
