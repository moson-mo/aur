package aur

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBy_String(t *testing.T) {
	tests := []struct {
		name string
		by   By
		want string
	}{
		{
			name: "name", by: Name, want: "name",
		},
		{
			name: "namedesc", by: NameDesc, want: "name-desc",
		},
		{
			name: "maintainer", by: Maintainer, want: "maintainer",
		},
		{
			name: "depends", by: Depends, want: "depends",
		},
		{
			name: "makedepends", by: MakeDepends, want: "makedepends",
		},
		{
			name: "optdepends", by: OptDepends, want: "optdepends",
		},
		{
			name: "optdepends", by: OptDepends, want: "optdepends",
		},
		{
			name: "checkdepends", by: CheckDepends, want: "checkdepends",
		},
		{
			name: "default", by: None, want: "",
		},
		{
			name: "panic", by: 23, want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.by == 23 {
				assert.Panics(t, func() { _ = tt.by.String() })
			} else {
				got := tt.by.String()
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
