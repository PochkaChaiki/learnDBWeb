package testprepare

import (
	"encoding/xml"
	"fmt"
)

type Quiz struct {
	XMLName  xml.Name   `xml:"quiz"`
	Question []question `xml:"question"`
}

func makeQuizPtr(name string) *Quiz {
	return &Quiz{
		Question: []question{
			question{
				Type: "category",
				Category: category{
					Text: fmt.Sprintf("$course$/top/%s", name),
				},
				Info: makeInfo(),
			},
		},
	}
}

type question struct {
	Type            string       `xml:"type,attr"`
	Name            name         `xml:"name"`
	QuestionText    questionText `xml:"questiontext"`
	GeneralFeedback feedback     `xml:"generalfeedback"`
	DefaultGrade    float64      `xml:"defaultgrade"`
	Penalty         float64      `xml:"penalty"`
	Hidden          int          `xml:"hidden"`
	IDNumber        string       `xml:"idnumber"`
	UseCase         int          `xml:"usecase"`
	Answer          []answer     `xml:"answer"`
	Category        category     `xml:"category,omitempty"`
	Info            info         `xml:"info,omitempty"`
}

func makeQuestion(name string) question {
	return question{
		Type: "shortanswer",
		Name: makeName(name),
		QuestionText: questionText{
			Format: "html",
		},
		GeneralFeedback: feedback{
			Format: "html",
		},
		Answer:   make([]answer, 0),
		Penalty:  1 / 3,
		Hidden:   0,
		IDNumber: "",
		UseCase:  0,
	}
}

type name struct {
	Text string `xml:"text"`
}

func makeName(text string) name {
	return name{
		Text: text,
	}
}

type questionText struct {
	Format string `xml:"format,attr"`
	Text   string `xml:",cdata"`
}

type answer struct {
	Fraction float64  `xml:"fraction,attr"`
	Format   string   `xml:"format,attr"`
	Text     string   `xml:"text"`
	Feedback feedback `xml:"feedback"`
}

func makeAnswer(text string, feedback string) answer {
	return answer{
		Fraction: 100,
		Format:   "moodle_auto_format",
		Text:     text,
		Feedback: makeFeedback(feedback),
	}
}

type feedback struct {
	Format string `xml:"format,attr"`
	Text   string `xml:",cdata"`
}

func makeFeedback(text string) feedback {
	return feedback{
		Format: "html",
		Text:   text,
	}
}

type category struct {
	Text string `xml:"text"`
}

type info struct {
	Format string `xml:"format,attr"`
	Text   string `xml:"text"`
}

func makeInfo() info {
	return info{
		Format: "html",
		Text:   "",
	}
}
