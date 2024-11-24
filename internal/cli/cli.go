package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/limadavida/sql2api/internal/models"
)

var (
	// Definindo estilos para as cores
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7F00"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7F00")).Bold(true)
	cursorStyle   = focusedStyle
	noStyle       = lipgloss.NewStyle()
	helpStyle     = blurredStyle
	buttonHovered = lipgloss.NewStyle().Foreground(lipgloss.Color("#22FF04"))
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Background(lipgloss.Color("#FF7F00"))
	subTitlestyle = titleStyle
	TITLE         = "[ SQL2API ]"
	SUBTITLE      = "Magic with SQL! A fast backend framework to spin up production-ready APIs with just SQL."
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	config     models.ProjectConfig
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 6), // Alterado para 6 campos de entrada principais
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Project. Ex: TodoApi"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Path. Ex: ./ "
		case 2:
			t.Placeholder = "Author. Ex: Santos Dumont "
		case 3:
			t.Placeholder = "Server. Ex: 8080"
		case 4:
			t.Placeholder = "Database Name: todo.db"
		case 5:
			t.Placeholder = "Database Type: [sqlite3, postgresql, mysql, sqlserver, mongodb]"
		}
		m.inputs[i] = t
	}
	m.config = models.ProjectConfig{
		Models: make(map[string]models.ModelFiles),
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				// Capture final data and show summary
				m.saveConfig()

				// Retorna a model preenchida após o "enter"
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(TITLE) + "\n\n")
	b.WriteString(subTitlestyle.Render(SUBTITLE) + "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredStyle.Render("[ Save ]")
	if m.focusIndex == len(m.inputs) {
		button = buttonHovered.Render("[ Save ]")
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", button)
	return b.String()
}

func (m *model) saveConfig() {
	setValue := func(value, defaultValue string) string {
		if value == "" {
			return defaultValue
		}
		return value
	}
	m.config.Project = setValue(m.inputs[0].Value(), "TodoApiExample")
	m.config.RootDir = setValue(m.inputs[1].Value(), "./examples/TodoApiExample")
	m.config.Author = setValue(m.inputs[2].Value(), "Santos Dummont")
	m.config.Servers = setValue(m.inputs[3].Value(), "localhost")
	m.config.Database.Name = setValue(m.inputs[4].Value(), "todo_example.db")
	m.config.Database.Type = setValue(m.inputs[5].Value(), "sqlite3")
}

func RunSetup() models.ProjectConfig {

	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
	m := finalModel.(model)
	return m.config
}
