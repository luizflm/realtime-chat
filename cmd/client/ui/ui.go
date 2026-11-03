package ui

import (
	"encoding/json"
	"net"
	"realtime-chat/internal/message"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// model é o estado da interface de chat (Bubble Tea).
type model struct {
	conn          net.Conn
	userName      string
	viewport      viewport.Model
	messages      []string
	textarea      textarea.Model
	senderStyle   lipgloss.Style
	incomingStyle lipgloss.Style
	err           error
}

func InitModel(conn net.Conn, userName string) model {
	ta := textarea.New()
	ta.Placeholder = "Envie uma mensagem..."
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Bem vindo ao chat!
Escreva uma mensagem e aperte enter para enviar.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	// Enter envia a mensagem; quebra de linha fica desabilitada.
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		conn:          conn,
		userName:      userName,
		textarea:      ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#5CE65C")),
		incomingStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#77CBDA")),
		err:           nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Mensagem recebida do servidor via goroutine em main.go.
	case message.Message:
		m.messages = append(m.messages, m.incomingStyle.Render(msg.SenderName+":")+string(msg.MessageText))

		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			msgText := m.textarea.Value()
			if strings.TrimSpace(msgText) == "" {
				return m, nil
			}

			request := message.ChatRequest{
				MessageText: msgText,
			}

			if err := json.NewEncoder(m.conn).Encode(request); err != nil {
				m.messages = append(m.messages, "System: Failed to send message.")
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
				m.viewport.GotoBottom()
				return m, nil
			}

			m.messages = append(m.messages, m.senderStyle.Render("Voce: ")+m.textarea.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()

			return m, nil

		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()

	}

	return m, nil
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.textarea.View())
	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}
