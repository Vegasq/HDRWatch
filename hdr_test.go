package hdrwatch

import (
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	type args struct {
		searchFor string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Ubuntu search", args{"ubuntu"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(tt.args.searchFor)
			if !strings.Contains(got.TorrentResults[0].Title, tt.args.searchFor) {
				t.Errorf("Not found")
			}
		})
	}
}
