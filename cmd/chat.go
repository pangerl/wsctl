package cmd

import (
	"context"
	"fmt"
	"os"
	"time"
	"vhagar/chat"
	"vhagar/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// Bubbletea 聊天 TUI model
type aiResponseMsg struct {
	reply string
	err   error
}

type loadingTickMsg struct{}

type chatModel struct {
	messages     []string
	textInput    textinput.Model
	history      []string
	historyIndex int
	ctx          context.Context
	loading      bool
	loadingFrame string
}

var loadingFrames = []string{".", "..", "..."}

const loadingInterval = 200 // ms

func initialChatModel(ctx context.Context) chatModel {
	ti := textinput.New()
	ti.Placeholder = "请输入你的问题"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	return chatModel{
		messages:     []string{"你是一个乐于助人、无害的AI助手。"},
		textInput:    ti,
		history:      []string{},
		historyIndex: -1,
		ctx:          ctx,
		loading:      false,
		loadingFrame: loadingFrames[0],
	}
}

func (m chatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil // loading时禁用输入
		}
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}
			if input == "exit" || input == "EXIT" || input == "Exit" {
				return m, tea.Quit
			}
			m.messages = append(m.messages, "你: "+input)
			m.history = append(m.history, input)
			m.historyIndex = len(m.history)
			m.textInput.SetValue("")
			m.loading = true
			m.loadingFrame = loadingFrames[0]
			// 启动动画tick，并异步调用AI
			return m, tea.Batch(loadingTick(), callAI(m.ctx, input, m.messages))
		case "up":
			if len(m.history) > 0 && m.historyIndex > 0 {
				m.historyIndex--
				m.textInput.SetValue(m.history[m.historyIndex])
			}
			return m, nil
		case "down":
			if len(m.history) > 0 && m.historyIndex < len(m.history)-1 {
				m.historyIndex++
				m.textInput.SetValue(m.history[m.historyIndex])
			} else {
				m.historyIndex = len(m.history)
				m.textInput.SetValue("")
			}
			return m, nil
		}
	case loadingTickMsg:
		if m.loading {
			// 切换动画帧
			idx := 0
			for i, f := range loadingFrames {
				if f == m.loadingFrame {
					idx = (i + 1) % len(loadingFrames)
					break
				}
			}
			m.loadingFrame = loadingFrames[idx]
			return m, loadingTick()
		}
		return m, nil
	case aiResponseMsg:
		m.loading = false
		if msg.err != nil {
			m.messages = append(m.messages, "AI: [出错] "+msg.err.Error())
		} else {
			m.messages = append(m.messages, "AI: "+msg.reply)
		}
		m.messages = append(m.messages, "") // 每轮问答后插入空行
		return m, nil
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m chatModel) View() string {
	s := "AI 聊天（Ctrl+C 退出）\n\n"
	for _, msg := range m.messages {
		s += msg + "\n"
	}
	if m.loading {
		s += "\nAI 正在思考" + m.loadingFrame + "\n"
	}
	s += "\n" + m.textInput.View()
	return s
}

func loadingTick() tea.Cmd {
	return tea.Tick(loadingInterval*1e6, func(time.Time) tea.Msg {
		return loadingTickMsg{}
	})
}

func callAI(ctx context.Context, input string, messages []string) tea.Cmd {
	return func() tea.Msg {
		// 这里简单拼接历史消息，实际可根据你的AI接口调整
		reply, err := chat.ChatWithAI(ctx, []any{
			map[string]any{"role": "user", "content": input},
		})
		return aiResponseMsg{reply: reply, err: err}
	}
}

func RunChatTUI() {
	ctx := context.Background()
	p := tea.NewProgram(initialChatModel(ctx))
	if err := p.Start(); err != nil {
		fmt.Println("出错:", err)
		os.Exit(1)
	}
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "AI 聊天命令",
	Long:  `与 AI 进行基础对话的命令。`,
	Run: func(cmd *cobra.Command, args []string) {
		aiCfg := &config.Config.AI
		if !aiCfg.Enable || aiCfg.Provider == "" {
			fmt.Println("AI 聊天功能未启用，请检查 config.toml 配置。")
			return
		}
		// 使用 Bubbletea TUI 聊天界面
		RunChatTUI()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
