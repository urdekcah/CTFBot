package ui

import (
	"fmt"
	"strings"
	"sync"

	"ctfbot.urdekcah.ru/utils"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Model struct {
	textInput textinput.Model
	viewport  viewport.Model
	logs      []string
	logger    *log.Logger
	err       error
	width     int
	height    int
	mu        sync.Mutex
	logPool   *sync.Pool
}

func NewModel(logger *log.Logger) *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()

	return &Model{
		textInput: ti,
		logs:      []string{},
		logger:    logger,
		logPool: &sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.logger.Info("Quit....")
			return m, tea.Quit
		case tea.KeyEnter:
			m.processCommand(m.textInput.Value())
			m.textInput.SetValue("")
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		viewportHeight := m.height - verticalMarginHeight
		if viewportHeight < 1 {
			m.err = fmt.Errorf("Window Too Small")
		} else {
			m.err = nil
			m.viewport = viewport.New(m.width, viewportHeight)
			m.viewport.Width = m.width
			m.viewport.Height = viewportHeight
			m.viewport.SetContent(m.viewportContent())
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) processCommand(cmd string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sanitizedCmd := utils.SanitizeInput(cmd)
	m.logger.Info("Execute Command", "command", sanitizedCmd)
	m.addLogEntry(fmt.Sprintf("Executed: %s", sanitizedCmd))
}

func (m *Model) addLogEntry(entry string) {
	m.logs = append(m.logs, entry)
	if len(m.logs) > 1000 {
		m.logs = m.logs[1:]
	}
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n\n   %s\n\n", m.err)
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m *Model) headerView() string {
	title := titleStyle.Render("CTFBot")
	line := strings.Repeat("─", utils.Max(0, m.width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Model) footerView() string {
	info := statusMessageStyle("Ctrl+C to quit")
	line := strings.Repeat(" ", utils.Max(0, m.width-lipgloss.Width(info)))
	input := m.textInput.View()
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, line, info),
		input)
}

func (m *Model) viewportContent() string {
	sb := m.logPool.Get().(*strings.Builder)
	defer m.logPool.Put(sb)
	sb.Reset()

	for _, log := range m.logs {
		sb.WriteString(log)
		sb.WriteString("\n")
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render(sb.String())
}
