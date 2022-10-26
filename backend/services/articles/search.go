package articles

import (
	"fmt"
	"strings"
)

func searchExact(q string) string {
	return fmt.Sprintf("\"%v\"", q)
}

func searchAND(q string) string {
	terms := strings.Split(q, "AND")

	t := strings.Builder{}

	for _, v := range terms {
		t.WriteString(" +" + v)
	}
	return t.String()
}

func searchOR(q string) string {
	terms := strings.Split(q, "OR")

	t := strings.Builder{}

	for _, v := range terms {
		t.WriteString(" " + v)
	}
	return t.String()
}

func searchNOT(q string) string {
	terms := strings.Split(q, "NOT")

	t := strings.Builder{}

	t.WriteString("+" + terms[0])
	t.WriteString(" ")
	t.WriteString("-" + terms[1])

	return t.String()
}
