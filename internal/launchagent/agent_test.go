package launchagent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteMenubarPlistRestartsOnlyAfterFailure(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	logs := filepath.Join(home, "Library", "Logs", "theirtime")
	if err := writeMenubarPlist("/usr/local/bin/theirtime", logs); err != nil {
		t.Fatalf("writeMenubarPlist() error = %v", err)
	}

	path := filepath.Join(home, "Library", "LaunchAgents", menubarLabel+".plist")
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read plist: %v", err)
	}

	plist := string(contents)
	if !strings.Contains(plist, "<key>KeepAlive</key>\n  <dict>\n    <key>SuccessfulExit</key>\n    <false/>\n  </dict>") {
		t.Fatalf("KeepAlive should only restart after unsuccessful exits:\n%s", plist)
	}
	if strings.Contains(plist, "<key>KeepAlive</key>\n  <true/>") {
		t.Fatalf("KeepAlive must not restart after clean Quit:\n%s", plist)
	}
}
