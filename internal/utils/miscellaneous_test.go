// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"testing"

	. "github.com/onsi/gomega"
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
			t.Parallel()
			g.Expect(MergeMaps(tc.sourceMaps...)).To(Equal(tc.expected))
		})
	}
}
