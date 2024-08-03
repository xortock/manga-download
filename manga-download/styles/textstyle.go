package styles

import "github.com/charmbracelet/lipgloss"

func GetSuccessTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#16DE5C"))
}

func RenderSuccess(value ...string) string {
	return GetSuccessTextStyle().Render(value...)
}

func GetFailedTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#F52032"))
}

func RenderFailed(value ...string) string {
	return GetFailedTextStyle().Render(value...)
}