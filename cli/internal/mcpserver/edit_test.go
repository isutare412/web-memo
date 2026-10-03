package mcpserver

import (
	"strings"
	"testing"
)

func TestApplyEdit(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		oldStr     string
		newStr     string
		replaceAll bool
		want       string
		wantErr    string
	}{
		{name: "single", content: "a b c", oldStr: "b", newStr: "X", want: "a X c"},
		{name: "replace all", content: "a b b", oldStr: "b", newStr: "X", replaceAll: true, want: "a X X"},
		{name: "multiline", content: "l1\nl2\nl3", oldStr: "l2\n", newStr: "", want: "l1\nl3"},
		{name: "empty old", content: "abc", oldStr: "", newStr: "X", wantErr: "must not be empty"},
		{name: "empty old replace all", content: "abc", oldStr: "", newStr: "X", replaceAll: true, wantErr: "must not be empty"},
		{name: "no match", content: "abc", oldStr: "z", newStr: "X", wantErr: "0 matches"},
		{name: "ambiguous", content: "b b b", oldStr: "b", newStr: "X", wantErr: "3 places"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyEdit(tt.content, tt.oldStr, tt.newStr, tt.replaceAll)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
