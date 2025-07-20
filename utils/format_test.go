package utils

import (
	"reflect"
	"testing"
)

func TestCallUser(t *testing.T) {
	tests := []struct {
		users    []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{"user1"}, "<@user1>"},
		{[]string{"user1", "user2"}, "<@user1><@user2>"},
		{[]string{"alice", "bob", "charlie"}, "<@alice><@bob><@charlie>"},
	}

	for _, test := range tests {
		result := CallUser(test.users)
		if result != test.expected {
			t.Errorf("CallUser(%v) = %s, want %s", test.users, result, test.expected)
		}
	}
}

func TestCallUserWithSeparator(t *testing.T) {
	tests := []struct {
		users     []string
		separator string
		expected  string
	}{
		{[]string{}, " ", ""},
		{[]string{"user1"}, " ", "<@user1>"},
		{[]string{"user1", "user2"}, " ", "<@user1> <@user2>"},
		{[]string{"alice", "bob"}, ", ", "<@alice>, <@bob>"},
		{[]string{"a", "b", "c"}, " | ", "<@a> | <@b> | <@c>"},
	}

	for _, test := range tests {
		result := CallUserWithSeparator(test.users, test.separator)
		if result != test.expected {
			t.Errorf("CallUserWithSeparator(%v, %s) = %s, want %s", test.users, test.separator, result, test.expected)
		}
	}
}

func TestFormatUserMention(t *testing.T) {
	tests := []struct {
		user     string
		expected string
	}{
		{"", ""},
		{"user1", "<@user1>"},
		{"alice", "<@alice>"},
	}

	for _, test := range tests {
		result := FormatUserMention(test.user)
		if result != test.expected {
			t.Errorf("FormatUserMention(%s) = %s, want %s", test.user, result, test.expected)
		}
	}
}

func TestFormatChannelMention(t *testing.T) {
	tests := []struct {
		channel  string
		expected string
	}{
		{"", ""},
		{"general", "<#general>"},
		{"123456789", "<#123456789>"},
	}

	for _, test := range tests {
		result := FormatChannelMention(test.channel)
		if result != test.expected {
			t.Errorf("FormatChannelMention(%s) = %s, want %s", test.channel, result, test.expected)
		}
	}
}

func TestFormatRoleMention(t *testing.T) {
	tests := []struct {
		role     string
		expected string
	}{
		{"", ""},
		{"admin", "<@&admin>"},
		{"123456789", "<@&123456789>"},
	}

	for _, test := range tests {
		result := FormatRoleMention(test.role)
		if result != test.expected {
			t.Errorf("FormatRoleMention(%s) = %s, want %s", test.role, result, test.expected)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		s        string
		maxLen   int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "he..."},
		{"hello world", 11, "hello world"},
		{"hello world", 8, "hello..."},
		{"hi", 2, "hi"},
		{"hello", 3, "hel"},
		{"hello", 1, "h"},
	}

	for _, test := range tests {
		result := TruncateString(test.s, test.maxLen)
		if result != test.expected {
			t.Errorf("TruncateString(%s, %d) = %s, want %s", test.s, test.maxLen, result, test.expected)
		}
	}
}

func TestPadString(t *testing.T) {
	tests := []struct {
		s        string
		length   int
		expected string
	}{
		{"hello", 10, "hello     "},
		{"hello world", 5, "hello world"},
		{"hi", 5, "hi   "},
		{"", 3, "   "},
	}

	for _, test := range tests {
		result := PadString(test.s, test.length)
		if result != test.expected {
			t.Errorf("PadString(%s, %d) = %q, want %q", test.s, test.length, result, test.expected)
		}
	}
}

func TestPadStringLeft(t *testing.T) {
	tests := []struct {
		s        string
		length   int
		expected string
	}{
		{"hello", 10, "     hello"},
		{"hello world", 5, "hello world"},
		{"hi", 5, "   hi"},
		{"", 3, "   "},
	}

	for _, test := range tests {
		result := PadStringLeft(test.s, test.length)
		if result != test.expected {
			t.Errorf("PadStringLeft(%s, %d) = %q, want %q", test.s, test.length, result, test.expected)
		}
	}
}

func TestCenterString(t *testing.T) {
	tests := []struct {
		s        string
		length   int
		expected string
	}{
		{"hello", 10, "  hello   "},
		{"hello world", 5, "hello world"},
		{"hi", 6, "  hi  "},
		{"hi", 5, " hi  "},
		{"", 4, "    "},
	}

	for _, test := range tests {
		result := CenterString(test.s, test.length)
		if result != test.expected {
			t.Errorf("CenterString(%s, %d) = %q, want %q", test.s, test.length, result, test.expected)
		}
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	tests := []struct {
		slice    []string
		expected []string
	}{
		{[]string{}, []string{}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "", "c"}, []string{"a", "c"}},
		{[]string{"", "", ""}, []string{}},
		{[]string{"a", "  ", "c"}, []string{"a", "c"}},
		{[]string{" a ", " b ", " c "}, []string{" a ", " b ", " c "}},
	}

	for _, test := range tests {
		result := RemoveEmptyStrings(test.slice)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("RemoveEmptyStrings(%v) = %v, want %v", test.slice, result, test.expected)
		}
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		s        string
		sep      string
		expected []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"a, b, c", ",", []string{"a", "b", "c"}},
		{"a,  ,c", ",", []string{"a", "c"}},
		{"", ",", []string{}},
		{"a", ",", []string{"a"}},
		{"a|b|c", "|", []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		result := SplitAndTrim(test.s, test.sep)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SplitAndTrim(%s, %s) = %v, want %v", test.s, test.sep, result, test.expected)
		}
	}
}

func TestJoinNonEmpty(t *testing.T) {
	tests := []struct {
		slice    []string
		sep      string
		expected string
	}{
		{[]string{"a", "b", "c"}, ",", "a,b,c"},
		{[]string{"a", "", "c"}, ",", "a,c"},
		{[]string{"", "", ""}, ",", ""},
		{[]string{"a"}, ",", "a"},
		{[]string{}, ",", ""},
	}

	for _, test := range tests {
		result := JoinNonEmpty(test.slice, test.sep)
		if result != test.expected {
			t.Errorf("JoinNonEmpty(%v, %s) = %s, want %s", test.slice, test.sep, result, test.expected)
		}
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		s          string
		substrings []string
		expected   bool
	}{
		{"hello world", []string{"hello", "foo"}, true},
		{"hello world", []string{"foo", "bar"}, false},
		{"hello world", []string{"world"}, true},
		{"hello world", []string{}, false},
		{"", []string{"hello"}, false},
	}

	for _, test := range tests {
		result := ContainsAny(test.s, test.substrings)
		if result != test.expected {
			t.Errorf("ContainsAny(%s, %v) = %t, want %t", test.s, test.substrings, result, test.expected)
		}
	}
}

func TestContainsAll(t *testing.T) {
	tests := []struct {
		s          string
		substrings []string
		expected   bool
	}{
		{"hello world", []string{"hello", "world"}, true},
		{"hello world", []string{"hello", "foo"}, false},
		{"hello world", []string{"world"}, true},
		{"hello world", []string{}, true},
		{"", []string{"hello"}, false},
	}

	for _, test := range tests {
		result := ContainsAll(test.s, test.substrings)
		if result != test.expected {
			t.Errorf("ContainsAll(%s, %v) = %t, want %t", test.s, test.substrings, result, test.expected)
		}
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"你好", "好你"},
	}

	for _, test := range tests {
		result := ReverseString(test.s)
		if result != test.expected {
			t.Errorf("ReverseString(%s) = %s, want %s", test.s, result, test.expected)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"hello", false},
		{" hello ", false},
	}

	for _, test := range tests {
		result := IsEmpty(test.s)
		if result != test.expected {
			t.Errorf("IsEmpty(%q) = %t, want %t", test.s, result, test.expected)
		}
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{"", false},
		{"   ", false},
		{"\t\n", false},
		{"hello", true},
		{" hello ", true},
	}

	for _, test := range tests {
		result := IsNotEmpty(test.s)
		if result != test.expected {
			t.Errorf("IsNotEmpty(%q) = %t, want %t", test.s, result, test.expected)
		}
	}
}
