package main

import "testing"

func TestURLHelper(t *testing.T) {
	tt := []struct {
		name         string
		input        []string
		expectedSegs []string
	}{
		{
			name:         "api",
			input:        []string{"/api/users", "api"},
			expectedSegs: []string{"api", "users"},
		},
		{
			name:         "users",
			input:        []string{"/api/users", "users"},
			expectedSegs: []string{"users"},
		},
		{
			name:         "index",
			input:        []string{"/api/users", ""},
			expectedSegs: []string{"", "api", "users"},
		},
		{
			name:         "index",
			input:        []string{"/", ""},
			expectedSegs: []string{""},
		},
		{
			name:         "not found",
			input:        []string{"/api/usersss", "users"},
			expectedSegs: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			segs := urlHelper(tc.input[0], tc.input[1])
			if len(segs) == 0 {
				if len(tc.expectedSegs) != 0 {
					t.Errorf("expected there to be no segs and got: %v", len(segs))
				}
			}
			for i, seg := range segs {
				if seg != tc.expectedSegs[i] {
					t.Errorf("expected seg %v to be %v got: %v", i, tc.expectedSegs[i], seg)
				}
			}
		})
	}

}
