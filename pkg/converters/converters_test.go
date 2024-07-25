package converters

import (
	"testing"
)

func TestStringAsInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"123\n", 123},
		{"4567\n", 4567},
		{"0\n", 0},
		{"-789\n", -789},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := StringAsInt(test.input)
			if result != test.expected {
				t.Errorf("StringAsInt(%q) = %d; expected %d", test.input, result, test.expected)
			}
		})
	}
}

func TestStringAsInt_Panic(t *testing.T) {
	tests := []string{
		"abc\n",
		"123abc\n",
		"\n",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("StringAsInt(%q) did not panic", input)
				}
			}()
			StringAsInt(input)
		})
	}
}
