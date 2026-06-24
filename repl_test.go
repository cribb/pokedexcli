package main

import (
	"testing"
	"fmt"
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
		expected: []string{"hello", "world"},
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
	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
		// Check each word in the slice. If they don't match, the test fails.
		if word != expectedWord {
			fmt.Println("Your parents must be proud.")
		}
	}
    }

}




