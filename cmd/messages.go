package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Aayush9029/tmail/internal/api"
	"github.com/Aayush9029/tmail/internal/config"
	"github.com/Aayush9029/tmail/internal/ui"
)

// --- BubbleTea list item ---

type msgItem struct {
	summary api.MessageSummary
	rank    int
}

func (m msgItem) Title() string {
	from := m.summary.From.Address
	if m.summary.From.Name != "" {
		from = m.summary.From.Name
	}
	marker := " "
	if !m.summary.Seen {
		marker = "●"
	}
	return fmt.Sprintf("%s %d. %s", marker, m.rank, from)
}

func (m msgItem) Description() string {
	return fmt.Sprintf("%s  %s", m.summary.Subject, formatDate(m.summary.CreatedAt))
}

func (m msgItem) FilterValue() string {
	return m.summary.From.Address + " " + m.summary.Subject
}

// --- BubbleTea model ---

type inboxModel struct {
	list     list.Model
	client   *api.Client
	msgs     []api.MessageSummary
	quitting bool
	opening  bool
}

func newInboxModel(client *api.Client, msgs []api.MessageSummary, address string) inboxModel {
	items := make([]list.Item, len(msgs))
	for i, m := range msgs {
		items[i] = msgItem{summary: m, rank: i + 1}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("📬 %s", address)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("36")).
		Bold(true).
		MarginLeft(2)
	l.SetShowStatusBar(true)
	l.DisableQuitKeybindings()

	return inboxModel{
		list:   l,
		client: client,
		msgs:   msgs,
	}
}

func (m inboxModel) Init() tea.Cmd { return nil }

func (m inboxModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(msgItem); ok {
				m.opening = true
				full, err := m.client.GetMessage(item.summary.ID)
				if err == nil {
					openMessageInBrowser(full)
				}
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m inboxModel) View() string {
	if m.quitting || m.opening {
		return ""
	}
	return m.list.View()
}

// --- Entry point ---

func Messages() {
	acct, err := config.Load()
	if err != nil {
		ui.Fatalf("%v", err)
	}

	client := api.NewWithToken(acct.Token)
	msgs, err := client.ListMessages()
	if err != nil {
		ui.Fatalf("failed to fetch messages: %v", err)
	}

	if len(msgs) == 0 {
		ui.Header("tmail")
		fmt.Printf("  %sinbox: %s%s\n\n", ui.Dim, acct.Address, ui.Reset)
		ui.Dimf("no messages yet")
		fmt.Println()
		return
	}

	// Non-interactive: plain table (for piping / Claude Code)
	if !ui.IsTTY() {
		printMessageTable(msgs, acct.Address)
		return
	}

	// Interactive: BubbleTea list
	m := newInboxModel(client, msgs, acct.Address)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		ui.Fatalf("TUI error: %v", err)
	}
}

func printMessageTable(msgs []api.MessageSummary, address string) {
	w := ui.TermWidth()

	ui.Header("tmail")
	fmt.Printf("  %sinbox: %s%s\n\n", ui.Dim, address, ui.Reset)

	overhead := 18
	avail := w - overhead
	if avail < 30 {
		avail = 30
	}
	fromW := avail * 3 / 10
	if fromW < 12 {
		fromW = 12
	}
	subjW := avail - fromW
	if subjW < 12 {
		subjW = 12
	}

	hdrFmt := fmt.Sprintf("  %%s%%-%ds %%-%ds %%-%ds %%s%%s\n", 4, fromW, subjW)
	fmt.Printf(hdrFmt, ui.Dim, "#", "FROM", "SUBJECT", "DATE", ui.Reset)
	fmt.Printf("  %s%s%s\n", ui.Dim, strings.Repeat("─", w-4), ui.Reset)

	rowFmt := fmt.Sprintf(" %%s%%-%dd %%-%ds %%-%ds %%s\n", 4, fromW, subjW)
	for i, m := range msgs {
		from := m.From.Address
		if m.From.Name != "" {
			from = m.From.Name
		}
		from = truncate(from, fromW)
		subject := truncate(m.Subject, subjW)
		date := formatDate(m.CreatedAt)

		marker := " "
		if !m.Seen {
			marker = ui.Cyan + "●" + ui.Reset
		}
		fmt.Printf(rowFmt, marker, i+1, from, subject, date)
	}
	fmt.Println()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func formatDate(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso[:10]
	}
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04")
	}
	return t.Format("Jan 02")
}

func openMessageInBrowser(msg *api.MessageFull) {
	htmlContent := strings.Join(msg.HTML, "\n")
	if htmlContent == "" {
		htmlContent = "<pre>" + msg.Text + "</pre>"
	}

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "tmail-message.html")
	if err := os.WriteFile(tmpFile, []byte(htmlContent), 0o644); err != nil {
		ui.Fatalf("failed to write temp file: %v", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-a", "Safari", tmpFile)
	default:
		cmd = exec.Command("xdg-open", tmpFile)
	}
	if err := cmd.Start(); err != nil {
		ui.Fatalf("failed to open browser: %v", err)
	}
}
