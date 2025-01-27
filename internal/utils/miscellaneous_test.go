package utils

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestMergeStringMaps(t *testing.T) {
	testCases := []struct {
		name       string
		sourceMaps []map[string]string
		expected   map[string]string
	}{
		{
			name:       "nil source maps",
			sourceMaps: nil,
			expected:   nil,
		},
		{
			name: "one of the source maps is nil",
			sourceMaps: []map[string]string{
				{"key1": "value1"},
				nil,
			},
			expected: map[string]string{"key1": "value1"},
		},
		{
			name: "source maps are empty",
			sourceMaps: []map[string]string{
				{}, {},
			},
			expected: map[string]string{},
		},
		{
			name: "source maps are neither nil nor empty",
			sourceMaps: []map[string]string{
				{"key1": "value1"},
				{"key2": "value2"},
				{"key1": "value3"},
			},
			expected: map[string]string{
				"key1": "value3",
				"key2": "value2",
			},
		},
	}

	g := NewWithT(t)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g.Expect(MergeMaps(tc.sourceMaps...)).To(Equal(tc.expected))
		})
	}
}
