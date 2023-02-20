package assert

import "strings"

type TestingT interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

func Contains(t TestingT, s string, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected to contain %q:\n%s", substr, s)
	}
}

func Equals[E comparable](t TestingT, expected, actual E) {
	t.Helper()
	if expected != actual {
		t.Fatalf("\nexpected:\n'%+v'\nactual:\n'%+v'\n", expected, actual)
	}
}

func Nil(t TestingT, a any) {
	t.Helper()
	if a != nil {
		t.Fatalf("\nexpected nil, got %+v\n", a)
	}
}

func NotNil(t TestingT, a any) {
	t.Helper()
	if a == nil {
		t.Fatal("expected not nil")
	}
}

func SliceEquals[E comparable](t TestingT, expected, actual []E) {
	t.Helper()
	if actual == nil && expected == nil {
		return
	} else if actual != nil && expected != nil {
		expectedSize := len(expected)
		actualSize := len(actual)
		if len(expected) != len(actual) {
			t.Fatalf("slices are of different lengths:\nlen(expected)=%d\nlen(actual)=%d", expectedSize, actualSize)
		}
		for i := 0; i < expectedSize; i++ {
			if expected[i] != actual[i] {
				t.Fatalf("index %d of slices are not equal:\nexpected:\n%+v\nactual:\n%+v\n", i, expected[i], actual[i])
			}
		}
		return
	}
	t.Fatalf("slices are not equal:\nexpected:\n%v\nactual:\n%v\n", expected, actual)
}

func True(t TestingT, condition bool) {
	t.Helper()
	if !condition {
		t.Fatal("expected true")
	}
}
