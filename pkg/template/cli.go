package template

import "github.com/charmbracelet/lipgloss"

func GetCliTemplateFns() map[string]any {
	colorCache := make(map[string]func(string) string)

	return map[string]any{
		"bold":      lipgloss.NewStyle().Bold(true).Render,
		"italic":    lipgloss.NewStyle().Italic(true).Render,
		"underline": lipgloss.NewStyle().Underline(true).Render,
		"color": func(text, color string) string {
			if color == "" {
				return text
			}
			if _, ok := colorCache[color]; !ok {
				styleFn := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
				colorCache[color] = func(s string) string {
					return styleFn(s)
				}
			}

			return colorCache[color](text)
		},
		"ternary": func(condition bool, a, b any) any {
			if condition {
				return a
			}
			return b
		},
		"rangeseq": func(n int) []int {
			seq := make([]int, n)
			for i := 0; i < n; i++ {
				seq[i] = i
			}
			return seq
		},
	}
}
