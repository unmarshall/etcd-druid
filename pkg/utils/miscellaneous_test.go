package utils

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	"k8s.io/utils/pointer"

	. "github.com/onsi/gomega"
)

type TestStringAlias string

const defaultTestStringAlias TestStringAlias = "default-test-string-alias"

var _ = Describe("TypeDeref", func() {
	var (
		testStringAliasValue TestStringAlias = "test-value"
	)
	DescribeTable("string type deref",
		func(actualVal *string, defaultVal string, expectedVal string) {
			Expect(TypeDeref[string](actualVal, defaultVal)).To(Equal(expectedVal))
		},
		Entry("nil actual string val should return default string val", nil, "bingo", "bingo"),
		Entry("non-nil actual string val should return the actual string val and not the default val", pointer.String("tringo"), "zingo", "tringo"),
	)
	DescribeTable("int32 type deref",
		func(actualVal *string, defaultVal string, expectedVal string) {
			Expect(TypeDeref[string](actualVal, defaultVal)).To(Equal(expectedVal))
		},
		Entry("nil actual string val should return default string val", nil, "bingo", "bingo"),
		Entry("non-nil actual string val should return the actual string val and not the default val", pointer.String("tringo"), "zingo", "tringo"),
	)
	DescribeTable("time.Duration struct type deref",
		func(actualVal *time.Duration, defaultVal time.Duration, expectedVal time.Duration) {
			Expect(TypeDeref[time.Duration](actualVal, defaultVal)).To(Equal(expectedVal))
		},
		Entry("nil actual duration val should return default duration val", nil, 2*time.Second, 2*time.Second),
		Entry("non-nil actual duration val should return the actual duration val and not the default val", pointer.Duration(5*time.Second), 2*time.Second, 5*time.Second),
	)
	DescribeTable("TestStringAlias type deref",
		func(actualVal *TestStringAlias, defaultVal TestStringAlias, expectedVal TestStringAlias) {
			Expect(TypeDeref[TestStringAlias](actualVal, defaultVal)).To(Equal(expectedVal))
		},
		Entry("nil actual TestStringAlias val should return default TestStringAlias val", nil, defaultTestStringAlias, defaultTestStringAlias),
		Entry("non-nil actual TestStringAlias val should return the actual TestStringAlias val and not the default val", &testStringAliasValue, defaultTestStringAlias, testStringAliasValue),
	)
})
