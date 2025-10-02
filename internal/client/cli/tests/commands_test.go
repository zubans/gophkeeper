package cli
import (
	"testing"
	"gophkeeper/internal/client/cli"
)
func TestParseCommand(t *testing.T) {
	_, err := cli.ParseCommand([]string{})
	if err == nil {
		t.Errorf("expected error for no command")
	}
	if _, err := cli.ParseCommand([]string{"unknown"}); err == nil {
		t.Errorf("expected error for unknown command")
	}
	if _, err := cli.ParseCommand([]string{"register", "u", "e@example.com"}); err == nil {
		t.Errorf("expected error for short username")
	}
}
