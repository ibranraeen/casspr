package dialog

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ibranraeen/casspr/internal/question"
	"github.com/ibranraeen/casspr/internal/ui/common"
	"github.com/ibranraeen/casspr/internal/ui/styles"
	uv "github.com/charmbracelet/ultraviolet"
)

// AskQuestionID is the identifier for the ask question dialog.
const AskQuestionID = "ask_question"

// AskQuestion represents a dialog for clarifying questions.
type AskQuestion struct {
	com          *common.Common
	req          question.QuestionRequest
	selectedIdx  int  // Currently focused item (0 to len(options)-1, plus custom, plus submit button)
	submitChose  bool // True if Submit button is focused/selected

	// State for options
	checkedOptions map[int]bool // Key is option index, value is checked status (for multi-select)
	customInput    textinput.Model

	width        int
	help         help.Model
	keyMap       askQuestionKeyMap
}

type askQuestionKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Tab    key.Binding
	Toggle key.Binding // Space to toggle checkboxes
	Select key.Binding // Enter to submit/select
	Close  key.Binding
}

func defaultAskQuestionKeyMap() askQuestionKeyMap {
	return askQuestionKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Close: CloseKey,
	}
}

var _ Dialog = (*AskQuestion)(nil)

// NewAskQuestion creates a new clarifying questions dialog.
func NewAskQuestion(com *common.Common, req question.QuestionRequest) *AskQuestion {
	h := help.New()
	h.Styles = com.Styles.DialogHelpStyles()

	t := com.Styles
	ti := textinput.New()
	ti.Placeholder = "Write custom response..."
	ti.SetStyles(t.TextInput)
	ti.Focus()

	q := &AskQuestion{
		com:            com,
		req:            req,
		checkedOptions: make(map[int]bool),
		customInput:    ti,
		help:           h,
		keyMap:         defaultAskQuestionKeyMap(),
		width:          75,
	}

	return q
}

// ID implements [Dialog].
func (*AskQuestion) ID() string {
	return AskQuestionID
}

// ToolCallID returns the tool call ID associated with this dialog.
func (q *AskQuestion) ToolCallID() string {
	return q.req.ToolCallID
}

func (q *AskQuestion) numItems() int {
	count := len(q.req.Options)
	if q.req.AllowCustom {
		count++
	}
	return count
}

func (q *AskQuestion) isCustomIndex(idx int) bool {
	return q.req.AllowCustom && idx == len(q.req.Options)
}

// HandleMsg implements [Dialog].
func (q *AskQuestion) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, q.keyMap.Close):
			// Close / cancel denies or skips the question
			return ActionClose{}

		case key.Matches(msg, q.keyMap.Up):
			if q.submitChose {
				q.submitChose = false
				q.selectedIdx = q.numItems() - 1
			} else if q.selectedIdx > 0 {
				q.selectedIdx--
			}

		case key.Matches(msg, q.keyMap.Down):
			if !q.submitChose {
				if q.selectedIdx < q.numItems()-1 {
					q.selectedIdx++
				} else {
					q.submitChose = true
				}
			}

		case key.Matches(msg, q.keyMap.Tab):
			if q.submitChose {
				q.submitChose = false
				q.selectedIdx = 0
			} else if q.selectedIdx < q.numItems()-1 {
				q.selectedIdx++
			} else {
				q.submitChose = true
			}

		case key.Matches(msg, q.keyMap.Toggle):
			if !q.submitChose && q.req.IsMultiSelect && !q.isCustomIndex(q.selectedIdx) {
				q.checkedOptions[q.selectedIdx] = !q.checkedOptions[q.selectedIdx]
			}

		case key.Matches(msg, q.keyMap.Select):
			if q.submitChose {
				return q.submitResponse()
			}
			if !q.req.IsMultiSelect {
				// Single select - selecting option submits directly or selects it.
				// If custom write-in is focused, user must select submit button or
				// press enter. Let's make enter on custom write-in or single option submit.
				if q.isCustomIndex(q.selectedIdx) {
					// User pressed enter on custom input. Submit!
					return q.submitResponse()
				}
				// Selecting single-select option selects and submits immediately.
				q.checkedOptions[q.selectedIdx] = true
				return q.submitResponse()
			} else {
				// Multi-select - toggle checkbox
				if !q.isCustomIndex(q.selectedIdx) {
					q.checkedOptions[q.selectedIdx] = !q.checkedOptions[q.selectedIdx]
				} else {
					// Focusing custom, enter moves to submit
					q.submitChose = true
				}
			}

		case !q.isCustomIndex(q.selectedIdx) && !q.submitChose && len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9':
			optIdx := int(msg.String()[0] - '1')
			if optIdx < len(q.req.Options) {
				q.selectedIdx = optIdx
				if !q.req.IsMultiSelect {
					q.checkedOptions[optIdx] = true
					return q.submitResponse()
				} else {
					q.checkedOptions[optIdx] = !q.checkedOptions[optIdx]
				}
			}

		default:
			// If focusing custom write-in, forward keys to textinput
			if q.isCustomIndex(q.selectedIdx) && !q.submitChose {
				var cmd tea.Cmd
				q.customInput, cmd = q.customInput.Update(msg)
				if cmd != nil {
					return ActionCmd{cmd}
				}
			}
		}
	}
	return nil
}

func (q *AskQuestion) submitResponse() Action {
	var selected []string
	for i, opt := range q.req.Options {
		if q.checkedOptions[i] {
			selected = append(selected, opt)
		}
	}

	customAnswer := ""
	if q.req.AllowCustom {
		customAnswer = q.customInput.Value()
		// If custom answer has text, and they haven't explicitly checked custom,
		// we treat custom as implicitly included.
		if customAnswer != "" {
			selected = append(selected, "Custom")
		}
	}

	return ActionQuestionResponse{
		Request: q.req,
		Response: question.QuestionResponse{
			QuestionID:      q.req.ID,
			SelectedOptions: selected,
			CustomAnswer:    customAnswer,
		},
	}
}

// Draw implements [Dialog].
func (q *AskQuestion) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := q.com.Styles

	dialogStyle := t.Dialog.View.Width(q.width).Padding(0, 1)
	contentWidth := q.width - t.Dialog.View.GetHorizontalFrameSize() - 2

	header := q.renderHeader(contentWidth)
	optionsView := q.renderOptions(contentWidth)
	submitBtn := q.renderSubmitButton(contentWidth)
	helpView := q.help.View(q)

	parts := []string{
		header,
		"",
		optionsView,
		"",
		submitBtn,
		"",
		helpView,
	}

	innerContent := lipgloss.JoinVertical(lipgloss.Left, parts...)
	view := dialogStyle.Render(innerContent)

	var cur *tea.Cursor
	if q.isCustomIndex(q.selectedIdx) && !q.submitChose {
		// Calculate cursor position for text input
		cur = &tea.Cursor{}
		if inputCursor := q.customInput.Cursor(); inputCursor != nil {
			cur.X = inputCursor.X
		}
		// Adjust relative to customInput position
		lines := strings.Split(innerContent, "\n")
		var customLineIdx int
		for i, l := range lines {
			if strings.Contains(l, "Custom:") || strings.Contains(l, q.customInput.Placeholder) {
				customLineIdx = i
				break
			}
		}
		cur.Y = customLineIdx

		// Calculate the prefix width
		prefix := "○ "
		if focused := !q.submitChose && q.isCustomIndex(q.selectedIdx); focused {
			prefix = "● "
		}
		customLabel := prefix + "Custom: "

		// Offset by the label width
		cur.X += lipgloss.Width(customLabel)

		// Offset by dialogStyle padding, borders, and margins
		cur.X += dialogStyle.GetBorderLeftSize() + dialogStyle.GetPaddingLeft() + dialogStyle.GetMarginLeft()
		cur.Y += dialogStyle.GetBorderTopSize() + dialogStyle.GetPaddingTop() + dialogStyle.GetMarginTop()
	}

	DrawCenterCursor(scr, area, view, cur)
	return cur
}

func (q *AskQuestion) renderHeader(contentWidth int) string {
	t := q.com.Styles

	title := common.DialogTitle(t, "Clarifying Question", contentWidth-t.Dialog.Title.GetHorizontalFrameSize(), t.Dialog.TitleGradFromColor, t.Dialog.TitleGradToColor)
	title = t.Dialog.Title.Render(title)

	questionText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.CheckIcon)).
		Width(contentWidth).
		Render(q.req.Question)

	return lipgloss.JoinVertical(lipgloss.Left, title, "", questionText)
}

func (q *AskQuestion) renderOptions(contentWidth int) string {
	t := q.com.Styles
	var rendered []string

	normalStyle := t.Dialog.NormalItem.Width(contentWidth)
	selectedStyle := t.Dialog.SelectedItem.Width(contentWidth)

	for i, opt := range q.req.Options {
		focused := !q.submitChose && q.selectedIdx == i
		checked := q.checkedOptions[i]

		// Option prefix symbol with shortcut number hint
		var prefix string
		if q.req.IsMultiSelect {
			if checked {
				prefix = fmt.Sprintf("[✓] %d. ", i+1)
			} else {
				prefix = fmt.Sprintf("[ ] %d. ", i+1)
			}
		} else {
			if focused {
				prefix = fmt.Sprintf("● %d. ", i+1)
			} else {
				prefix = fmt.Sprintf("○ %d. ", i+1)
			}
		}

		// Text alignment & indents for option descriptions
		lines := strings.Split(opt, "\n")
		titlePart := prefix + lines[0]
		var renderedItem string

		if len(lines) > 1 {
			// Construct description lines with faint style
			var descParts []string
			for _, l := range lines[1:] {
				// Indent descriptions
				descParts = append(descParts, "    "+l)
			}
			descText := strings.Join(descParts, "\n")

			if focused {
				selDescStyle := selectedStyle.Faint(true)
				renderedItem = selectedStyle.Render(titlePart) + "\n" + selDescStyle.Render(descText)
			} else {
				normDescStyle := normalStyle.Faint(true)
				renderedItem = normalStyle.Render(titlePart) + "\n" + normDescStyle.Render(descText)
			}
		} else {
			if focused {
				renderedItem = selectedStyle.Render(titlePart)
			} else {
				renderedItem = normalStyle.Render(titlePart)
			}
		}
		rendered = append(rendered, renderedItem)
	}

	if q.req.AllowCustom {
		focused := !q.submitChose && q.isCustomIndex(q.selectedIdx)
		prefix := "○ "
		if focused {
			prefix = "● "
		}

		customLabel := prefix + "Custom: "
		q.customInput.SetWidth(max(0, contentWidth-lipgloss.Width(customLabel)-2))
		customView := customLabel + q.customInput.View()

		if focused {
			rendered = append(rendered, selectedStyle.Render(customView))
		} else {
			rendered = append(rendered, normalStyle.Render(customView))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (q *AskQuestion) renderSubmitButton(contentWidth int) string {
	t := q.com.Styles
	buttons := []common.ButtonOpts{
		{Text: "Submit", Selected: q.submitChose},
	}

	content := common.ButtonGroup(t, buttons, "  ")
	return lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Right).
		Render(content)
}

// ShortHelp implements [help.KeyMap].
func (q *AskQuestion) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		q.keyMap.Up,
		q.keyMap.Down,
		q.keyMap.Select,
		q.keyMap.Close,
	}

	if q.req.IsMultiSelect {
		bindings = append(bindings, q.keyMap.Toggle)
	}

	return bindings
}

// FullHelp implements [help.KeyMap].
func (q *AskQuestion) FullHelp() [][]key.Binding {
	return [][]key.Binding{q.ShortHelp()}
}
