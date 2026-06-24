package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {

    cases := []struct {
        input     string
        expected []string
    }{
	{
		input:     "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	{
		input: "HELLO WoRlD!   ",
		expected: []string{"hello", "world!"},
	},
	{
		input: "FOOS	BALL",
		expected: []string{"foos", "ball"},
	},
    }

    for _, c := range cases {
	actual := cleanInput(c.input)
	// Check the length of the actual slice against the expected slice	
	// If they don't match, the test fails. Use t.Errorf to print an error message.

	if len(actual) != len(c.expected) {
		t.Errorf("Initial length check: actual %v, expected %v. FAIL.",len(actual), len(c.input))
		continue 
	}

	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
		// Check each word in the slice. If they don't match, the test fails.
		if word != expectedWord {
			t.Errorf("test %v, expected %v. FAIL. Your parents must be proud.", word, expectedWord)
		}
	}
    }

}




