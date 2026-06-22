package version

import "testing"

func TestPrunePseudoVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "v0.78.1-0.20260619214928-1849ddd1abcd",
			expected: "v0.78.1",
		},
		{
			input:    "v0.78.1-pre.0.20260619214928-1849ddd1abcd",
			expected: "v0.78.1",
		},
		{
			input:    "v0.78.0-2-g1849ddd1",
			expected: "v0.78.0",
		},
		{
			input:    "v0.78.1-1849ddd1",
			expected: "v0.78.1",
		},
		{
			input:    "v0.78.1-0.20260619214928-1849ddd1f173+dirty",
			expected: "v0.78.1",
		},
		{
			input:    "v1.0.0",
			expected: "v1.0.0",
		},
		{
			input:    "v1.0.0-alpha",
			expected: "v1.0.0-alpha",
		},
		{
			input:    "v1.0.0+dirty",
			expected: "v1.0.0",
		},
		{
			input:    "devel",
			expected: "devel",
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			actual := prunePseudoVersion(tc.input)
			if actual != tc.expected {
				t.Errorf("prunePseudoVersion(%q) = %q; want %q", tc.input, actual, tc.expected)
			}
		})
	}
}
