package lib

import (
	"testing"
)

func TestShellescape(t *testing.T) {
	testCases := []struct {
		input, expected string
	}{
		{"", "''"},
		{"It's better to give than to receive", "It\\'s\\ better\\ to\\ give\\ than\\ to\\ receive"},
		{"$PATH", "\\$PATH"},
		{"hello world", "hello\\ world"},
		{`yq -i e ".clusters[0].cluster.proxy-url |=\"socks5://127.0.0.1:$TMUX_SSH_PORT\"" $KUBECONFIG`, "yq\\ -i\\ e\\ \\\".clusters\\[0\\].cluster.proxy-url\\ \\|\\=\\\\\\\"socks5://127.0.0.1:\\$TMUX_SSH_PORT\\\\\\\"\\\"\\ \\$KUBECONFIG"},
	}
	for _, tc := range testCases {
		actual := shellescape(tc.input)
		if actual != tc.expected {
			t.Errorf("shellescape(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}
