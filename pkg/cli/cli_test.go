package cli

import (
	"testing"

	"github.com/mcbattirola/gitnotes/pkg/gn"
	"github.com/stretchr/testify/assert"
)

func TestCheckInitParams(t *testing.T) {
	tt := []struct {
		name      string
		app       gn.GN
		expectErr bool
	}{
		{name: "it accepts np params", expectErr: false},
		{
			name: "it accepts no project and different branch",
			app: gn.GN{
				Branch: "another branch",
			},
			expectErr: false,
		},
		{
			name: "it accepts a different project and branch",
			app: gn.GN{
				Project: "another project",
				Branch:  "another branch",
			},
			expectErr: false,
		},
		{
			name: "it does not accept a different project if no branch is provided",
			app: gn.GN{
				Project: "another proj",
			},
			expectErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := checkInitParams(tc.app)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
