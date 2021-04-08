package errors

import "testing"

func TestFormatting(t *testing.T) {
	type test struct {
		expected string
		e        *Error
	}
	cases := []test{
		{
			expected: "inner",
			e: &Error{
				Msg: "inner",
				Err: nil,
			},
		},
		{
			expected: "outer:\n  inner",
			e: &Error{
				Msg: "outer",
				Err: &Error{
					Msg: "inner",
					Err: nil,
				},
			},
		},
	}
	for _, c := range cases {
		msg := c.e.Error()
		if msg != c.expected {
			t.Errorf("expected '%s', got: '%s'", c.expected, msg)
		}

	}
}
