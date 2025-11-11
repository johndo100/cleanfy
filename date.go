// date.go
// --------
// Formats timestamps according to the selected style (iso, compact, month, short, withtime).

package main

import "time"

func formatDate(t time.Time, style string) string {
	switch style {
	case "iso":
		return t.Format("2006-01-02")
	case "compact":
		return t.Format("20060102")
	case "month":
		return t.Format("2006-01")
	case "short":
		return t.Format("060102")
	case "withtime":
		return t.Format("2006-01-02T15.04.05")
	default:
		return t.Format("2006-01-02")
	}
}
