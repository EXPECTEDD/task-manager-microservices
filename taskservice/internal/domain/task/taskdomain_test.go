package taskdomain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTaskDomain(t *testing.T) {
	tests := []struct {
		testName string

		ProjectId   uint32
		Description string
		Deadline    time.Time

		expErr error
	}{
		{
			testName: "Success",

			ProjectId:   1,
			Description: "desc",
			Deadline:    time.Now(),

			expErr: nil,
		}, {
			testName: "Invalid project id",

			ProjectId:   0,
			Description: "desc",
			Deadline:    time.Now(),

			expErr: ErrInvalidProjectId,
		}, {
			testName: "Invalid description",

			ProjectId:   1,
			Description: strings.Repeat("a", 256),
			Deadline:    time.Now(),

			expErr: ErrInvalidDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			td, err := NewTaskDomain(tt.ProjectId, tt.Description, tt.Deadline)
			require.Equal(t, tt.expErr, err)
			if err == nil {
				require.Equal(t, uint32(0), td.Id)
				require.Equal(t, tt.ProjectId, td.ProjectId)
				require.Equal(t, tt.Description, td.Description)
				require.Equal(t, tt.Deadline, td.Deadline)
			}
		})
	}
}
