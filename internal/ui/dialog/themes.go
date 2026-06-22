package dialog

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/ibranraeen/casspr/internal/ui/common"
	"github.com/ibranraeen/casspr/internal/ui/list"
	"github.com/ibranraeen/casspr/internal/ui/styles"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/sahilm/fuzzy"
)

const (
	// ThemesID is the identifier for the theme picker dialog.
	ThemesID             = "themes"
	themesDialogMaxWidth  = 50
	themesDialogMaxHeight = 12
)

// ThemeOption represents a color theme option.
type ThemeOption struct {
	ID          string
	Title       string
	Description string
}

// AllThemes lists all available themes in order.
var AllThemes = []ThemeOption{
	{ID: "default", Title: "Default (Auto)", Description: "Automatically select theme based on LLM provider"},
	{ID: "light", Title: "Light", Description: "Bright and clean light-background theme"},
	{ID: "nord", Title: "Nord", Description: "Cool arctic slate-blue and frost theme"},
	{ID: "dracula", Title: "Dracula", Description: "Classic dark vampire purple-pink theme"},
	{ID: "cyberpunk", Title: "Cyberpunk", Description: "High contrast retro-neon cyan and hot pink"},
}

// Themes represents a dialog for selecting color theme.
type Themes struct {
	com   *common.Common
	help  help.Model
	list  *list.FilterableList
	input textinput.Model

	keyMap struct {
		Select   key.Binding
		Next     key.Binding
		Previous key.Binding
		UpDown   key.Binding
		Close    key.Binding
	}
}

// ThemeItem represents a theme list item.
type ThemeItem struct {
	*list.Versioned
	theme     ThemeOption
	isCurrent bool
	t         *styles.Styles
	m         fuzzy.Match
	cache     map[int]string
	focused   bool
}

// Finished implements list.Item.
func (t *ThemeItem) Finished() bool {
	return true
}

var (
	_ Dialog   = (*Themes)(nil)
	_ ListItem = (*ThemeItem)(nil)
)

// NewThemes creates a new theme picker dialog.
func NewThemes(com *common.Common) *Themes {
	t := &Themes{com: com}

	h := help.New()
	h.Styles = com.Styles.DialogHelpStyles()
	t.help = h

	t.list = list.NewFilterableList()
	t.list.Focus()

	t.input = textinput.New()
	t.input.SetVirtualCursor(false)
	t.input.Placeholder = "Type to filter"
	t.input.SetStyles(com.Styles.TextInput)
	t.input.Focus()

	t.keyMap.Select = key.NewBinding(
		key.WithKeys("enter", "ctrl+y"),
		key.WithHelp("enter", "confirm"),
	)
	t.keyMap.Next = key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓", "next item"),
	)
	t.keyMap.Previous = key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑", "previous item"),
	)
	t.keyMap.UpDown = key.NewBinding(
		key.WithKeys("up", "down"),
		key.WithHelp("↑/↓", "choose"),
	)
	t.keyMap.Close = CloseKey

	t.setItems()
	return t
}

// ID implements Dialog.
func (t *Themes) ID() string {
	return ThemesID
}

// HandleMsg implements [Dialog].
func (t *Themes) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, t.keyMap.Close):
			return ActionClose{}
		case key.Matches(msg, t.keyMap.Previous):
			t.list.Focus()
			if t.list.IsSelectedFirst() {
				t.list.SelectLast()
				t.list.ScrollToBottom()
				break
			}
			t.list.SelectPrev()
			t.list.ScrollToSelected()
		case key.Matches(msg, t.keyMap.Next):
			t.list.Focus()
			if t.list.IsSelectedLast() {
				t.list.SelectFirst()
				t.list.ScrollToTop()
				break
			}
			t.list.SelectNext()
			t.list.ScrollToSelected()
		case key.Matches(msg, t.keyMap.Select):
			selectedItem := t.list.SelectedItem()
			if selectedItem == nil {
				break
			}
			themeItem, ok := selectedItem.(*ThemeItem)
			if !ok {
				break
			}
			return ActionSelectTheme{Theme: themeItem.theme.ID}
		default:
			var cmd tea.Cmd
			t.input, cmd = t.input.Update(msg)
			value := t.input.Value()
			t.list.SetFilter(value)
			t.list.ScrollToTop()
			t.list.SetSelected(0)
			return ActionCmd{cmd}
		}
	}
	return nil
}

// Cursor returns the cursor position relative to the dialog.
func (t *Themes) Cursor() *tea.Cursor {
	return InputCursor(t.com.Styles, t.input.Cursor())
}

// Draw implements [Dialog].
func (t *Themes) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	sty := t.com.Styles
	width := max(0, min(themesDialogMaxWidth, area.Dx()))
	height := max(0, min(themesDialogMaxHeight, area.Dy()))
	innerWidth := width - sty.Dialog.View.GetHorizontalFrameSize()
	heightOffset := sty.Dialog.Title.GetVerticalFrameSize() + titleContentHeight +
		sty.Dialog.InputPrompt.GetVerticalFrameSize() + inputContentHeight +
		sty.Dialog.HelpView.GetVerticalFrameSize() +
		sty.Dialog.View.GetVerticalFrameSize()

	t.input.SetWidth(innerWidth - sty.Dialog.InputPrompt.GetHorizontalFrameSize() - 1)
	t.list.SetSize(innerWidth, height-heightOffset)
	t.help.SetWidth(innerWidth)

	rc := NewRenderContext(sty, width)
	rc.Title = "Select Theme"
	inputView := sty.Dialog.InputPrompt.Render(t.input.View())
	rc.AddPart(inputView)

	visibleCount := len(t.list.FilteredItems())
	if t.list.Height() >= visibleCount {
		t.list.ScrollToTop()
	} else {
		t.list.ScrollToSelected()
	}

	listView := sty.Dialog.List.Height(t.list.Height()).Render(t.list.Render())
	rc.AddPart(listView)
	rc.Help = t.help.View(t)

	view := rc.Render()

	cur := t.Cursor()
	DrawCenterCursor(scr, area, view, cur)
	return cur
}

// ShortHelp implements [help.KeyMap].
func (t *Themes) ShortHelp() []key.Binding {
	return []key.Binding{
		t.keyMap.UpDown,
		t.keyMap.Select,
		t.keyMap.Close,
	}
}

// FullHelp implements [help.KeyMap].
func (t *Themes) FullHelp() [][]key.Binding {
	m := [][]key.Binding{}
	slice := []key.Binding{
		t.keyMap.Select,
		t.keyMap.Next,
		t.keyMap.Previous,
		t.keyMap.Close,
	}
	for i := 0; i < len(slice); i += 4 {
		end := min(i+4, len(slice))
		m = append(m, slice[i:end])
	}
	return m
}

func (t *Themes) setItems() {
	cfg := t.com.Config()
	currentTheme := "default"
	if cfg != nil && cfg.Options != nil && cfg.Options.TUI != nil && cfg.Options.TUI.Theme != "" {
		currentTheme = cfg.Options.TUI.Theme
	}

	items := make([]list.FilterableItem, 0, len(AllThemes))
	selectedIndex := 0
	for i, theme := range AllThemes {
		item := &ThemeItem{
			Versioned: list.NewVersioned(),
			theme:     theme,
			isCurrent: theme.ID == currentTheme,
			t:         t.com.Styles,
		}
		items = append(items, item)
		if theme.ID == currentTheme {
			selectedIndex = i
		}
	}

	t.list.SetItems(items...)
	t.list.SetSelected(selectedIndex)
	t.list.ScrollToSelected()
}

// Filter returns the filter value for the theme item.
func (t *ThemeItem) Filter() string {
	return t.theme.Title
}

// ID returns the unique identifier for the theme.
func (t *ThemeItem) ID() string {
	return t.theme.ID
}

// SetFocused sets the focus state of the theme item.
func (t *ThemeItem) SetFocused(focused bool) {
	if t.focused == focused {
		return
	}
	t.cache = nil
	t.focused = focused
	if t.Versioned != nil {
		t.Bump()
	}
}

// SetMatch sets the fuzzy match for the theme item.
func (t *ThemeItem) SetMatch(m fuzzy.Match) {
	if sameFuzzyMatch(t.m, m) {
		return
	}
	t.cache = nil
	t.m = m
	if t.Versioned != nil {
		t.Bump()
	}
}

// Render returns the string representation of the theme item.
func (t *ThemeItem) Render(width int) string {
	info := ""
	if t.isCurrent {
		info = "current"
	}
	st := ListItemStyles{
		ItemBlurred:     t.t.Dialog.NormalItem,
		ItemFocused:     t.t.Dialog.SelectedItem,
		InfoTextBlurred: t.t.Dialog.ListItem.InfoBlurred,
		InfoTextFocused: t.t.Dialog.ListItem.InfoFocused,
	}
	return renderItem(st, t.theme.Title, info, t.focused, width, t.cache, &t.m)
}
