package main

import "testing"

func TestParseMapping(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input    string
		expected map[string]string
	}{
		"basic": {
			input: `
			/hello -> https://helloworld.com
			/world -> http://hellootherworld.com
			`,
			expected: map[string]string{
				"/hello": "https://helloworld.com",
				"/world": "http://hellootherworld.com",
			},
		},
	}
	for name, test := range tests {
		t.Run("name="+name, func(t *testing.T) {
			result, err := parseMapping(test.input)
			if err != nil {
				t.Error(err)
				return
			}
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d results, got %d\n", len(test.expected), len(result))
				return
			}
			for k, v := range test.expected {
				if v2, ok := result[k]; !ok {
					t.Errorf("Expected key %q in result, wasn't found\n", k)
				} else if v2 != v {
					t.Errorf("Expected key %q to be %q, was %q\n", k, v, v2)
				}
			}
		})
	}
}
