package terminal_test

import (
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
)

func TestProgressBar_String(t *testing.T) {
	testCases := []struct {
		width           int
		totalIncrements int
		increments      int
		expected        string
	}{
		{3, 27, 15, "█▋ "},
		{5, 51, 88, "█████"},
		{5, 68, 0, "     "},
		{5, 68, 1, "     "},
		{5, 68, 2, "▏    "},
		{5, 68, 67, "████▉"},
		{5, 68, 68, "█████"},
		{40, 243, 88, "██████████████▍                         "},
	}
	for i, tc := range testCases {
		pb := terminal.NewProgressBar("", "", "")
		pb.SetWidth(tc.width)
		pb.SetTotalIncrements(tc.totalIncrements)
		pb.Increment(tc.increments)
		if pb.String() != tc.expected {
			t.Fatalf("\ntest case #%d: %+v\n  actual: %q\nexpected: %q", i+1, tc, pb.String(), tc.expected)
		}
	}
}

func TestProgressBar_Increment(t *testing.T) {
	pb := terminal.NewProgressBar("", "", "")
	pb.SetWidth(5)
	pb.SetTotalIncrements(5)
	pb.Increment(4)
	pb.Increment(1)
	pb.Increment(1)
	if pb.String() != "█████" {
		t.Fatal()
	}
}
