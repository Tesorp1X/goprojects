package util_test

import (
	"slices"
	"testing"

	"github.com/Tesorp1X/goprojects/01-todo-list/tests/util"
)

func TestAssertEqualNotes(t *testing.T) {
	t.Run("two identical notes", func(t *testing.T) {

	})
	t.Run("two different notes", func(t *testing.T) {

	})
}

func TestAssertEqualRawData(t *testing.T) {
	assertEqual := func(t *testing.T, a, b bool) {
		t.Helper()
		if a != b {
			t.Errorf("want: %t, but got: %t", a, b)
		}
	}
	var (
		a [][]string
		b [][]string
		c [][]string
		d [][]string
	)
	for range 5 {
		a = append(a, []string{"1", "2", "sdfs", "123"})
	}
	b = slices.Clone(a)
	c = append(c, []string{"1", "2", "3", "123"})
	d = slices.Insert(b, 3, []string{"a", "b"})
	t.Run("two identicals matrices", func(t *testing.T) {
		want := true
		got := util.AssertEqualRawData(a, b)
		assertEqual(t, want, got)
	})
	t.Run("length doesn't match", func(t *testing.T) {
		want := false
		got := util.AssertEqualRawData(a, c)
		assertEqual(t, want, got)
	})
	t.Run("one line is shorter", func(t *testing.T) {
		want := false
		got := util.AssertEqualRawData(a, d)
		assertEqual(t, want, got)
	})
}
