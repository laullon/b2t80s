package ui

import "strings"

type Control interface {
	GetRegisters() string
	GetOutput() string
}

type RegTable [][]*RegText

func (t RegTable) Dump() string {
	var sb strings.Builder
	sb.WriteString("<table>")
	for _, row := range t {
		sb.WriteString("<tr>")
		for _, reg := range row {
			sb.WriteString("<td>")
			sb.WriteString(reg.label)
			sb.WriteString("</td><td>")
			sb.WriteString(reg.value)
			sb.WriteString("</td>")
		}
		sb.WriteString("</tr>")
	}
	sb.WriteString("</table>")
	return sb.String()
}

type RegText struct {
	label string
	value string
	style string
}

func NewRegText(label string) *RegText {
	rt := &RegText{
		label: label,
	}
	return rt
}

func (rt *RegText) Update(text string) {
	if rt.value != text {
		rt.value = text
		rt.style = `style="updated"`
	} else {
		rt.style = ``
	}
}
