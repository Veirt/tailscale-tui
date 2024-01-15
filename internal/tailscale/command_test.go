package tailscale

import (
	"testing"
)

func TestGetTailscaleUpFlags(t *testing.T) {
	flag := GetTailscaleUpFlags()

	if flag == nil {
		t.Fatalf("Failed to execute tailscale up --help")
	}

	if len(flag) < 1 {
		t.Fatalf("Flags are not parsed correctly. Expected at least one flag.")
	}

}
