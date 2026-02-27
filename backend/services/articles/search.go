package articles

import (
	"strings"
)

func searchExact(q string) string {
	// Join words with & so to_tsquery treats them as AND, not a syntax error
	return strings.Join(strings.Fields(q), " & ")
}

func searchAND(q string) string {
	// "climate AND change" → "climate & change"
	parts := strings.Split(q, " AND ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return strings.Join(parts, " & ")
}

func searchOR(q string) string {
	// "climate OR change" → "climate | change"
	parts := strings.Split(q, " OR ")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return strings.Join(parts, " | ")
}

func searchNOT(q string) string {
	// "climate NOT war" → "climate & !war"
	parts := strings.SplitN(q, " NOT ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]) + " & !" + strings.TrimSpace(parts[1])
	}
	return q
}

func searchQn(q string) string {
	return "isQn" + q
}
