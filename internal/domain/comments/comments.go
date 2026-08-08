// Package comments provides domain logic for commenting provided task's solution
package comments

var (
	ErrEmptyAnswer     *Comment = New("empty answer")
	ErrIncorrectAnswer *Comment = New("incorrect answer")
	ErrSyntaxError     *Comment = New("syntax error")
	ErrInternalError   *Comment = New("internal error")
)

type Comment struct {
	comments []string
}

func New(str string) *Comment {
	return &Comment{
		comments: []string{str},
	}
}

func (c *Comment) String() string {
	if c == nil {
		return ""
	}

	if len(c.comments) == 0 {
		return ""
	}

	buf := make([]rune, 0, len(c.comments)*len(c.comments[0]))
	for _, comm := range c.comments {
		buf = append(buf, []rune(comm)...)
		buf = append(buf, '\n')
	}

	return string(buf)
}

func Join(comments ...any) *Comment {
	joinComms := make([]string, 0, len(comments))
	for _, comm := range comments {
		switch c := comm.(type) {
		case *Comment:
			if comm != nil {
				joinComms = append(joinComms, c.comments...)
			}
		case string:
			joinComms = append(joinComms, c)
		}
	}

	if len(joinComms) == 0 {
		return nil
	}

	return &Comment{
		comments: joinComms,
	}
}
