package lib

import (
	"os"
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

func TestParseConfigPath(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	testCases := []struct {
		directory, project, expected string
	}{
		{"~/.tmuxinator", "gmux", homeDir + "/.tmuxinator/gmux.yml"},
		{"~/.tmuxinator", "gmux.yml", homeDir + "/.tmuxinator/gmux.yml"},
		{"~/.tmuxinator", "aa", homeDir + "/.tmuxinator/aa.yml"},
	}
	for _, tc := range testCases {
		actual := ParseConfigPath(tc.directory, tc.project)
		if actual != tc.expected {
			t.Errorf(
				"parseConfigPath(%q, %q) = %q; want %q",
				tc.directory, tc.project, actual, tc.expected)
		}
	}
}
