package comments

import (
	"slices"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input *Comment
		want  string
	}{
		{"print comment", &Comment{[]string{"1", "2"}}, "1\n2\n"},
		{"nil comment", nil, ""},
		{"no comments", &Comment{[]string{}}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.input.String()
			if str != tt.want {
				t.Errorf("want: %v, got: %v", tt.want, str)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name  string
		input []any
		want  *Comment
	}{
		{
			"all comments",
			[]any{New("1"), New("2"), New("3"), New("4")},
			&Comment{[]string{"1", "2", "3", "4"}},
		},
		{
			"all strings",
			[]any{"1", "2", "3", "4"},
			&Comment{[]string{"1", "2", "3", "4"}},
		},
		{
			"comments and strings",
			[]any{New("1"), New("2"), "3", "4"},
			&Comment{[]string{"1", "2", "3", "4"}},
		},
		{
			"0 items",
			[]any{},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Join(tt.input...)
			if tt.want == nil && got == nil {
				return
			}
			if tt.want == nil || got == nil {
				t.Errorf("want %v, got: %v", tt.want, got)
				return
			}
			if !slices.Equal(tt.want.comments, got.comments) {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}
