package articles

import (
	"strings"
)

func searchExact(q string) string {
	return q
}

func searchAND(q string) string {
	return strings.ReplaceAll(q, "AND", "&")
}

func searchOR(q string) string {
	return strings.ReplaceAll(q, "OR", "|")
}

func searchNOT(q string) string {
	return strings.ReplaceAll(q, "NOT ", "& !")
}

func searchQn(q string) string {
	return "isQn" + q
}
