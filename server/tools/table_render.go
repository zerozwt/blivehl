package main

import "strings"

func RenderTable(header []string, content [][]string) string {
	builder := strings.Builder{}

	colWidth := map[int]int{}
	updateWidth := func(line []string) {
		for idx, content := range line {
			if len(content)+2 > colWidth[idx] {
				colWidth[idx] = len(content) + 2
			}
		}
	}

	updateWidth(header)
	for _, line := range content {
		updateWidth(line)
	}

	renderLine := func(line []string, whitespace string) {
		for idx, content := range line {
			spaces := colWidth[idx] - len(content)
			before := spaces >> 1
			after := spaces - before
			builder.WriteString("|")
			builder.WriteString(strings.Repeat(whitespace, before))
			builder.WriteString(content)
			builder.WriteString(strings.Repeat(whitespace, after))
		}
		builder.WriteString("|\n")
	}

	renderLine(make([]string, len(header)), "-")
	renderLine(header, " ")
	renderLine(make([]string, len(header)), "-")
	for _, line := range content {
		renderLine(line, " ")
	}
	renderLine(make([]string, len(header)), "-")

	return builder.String()
}
