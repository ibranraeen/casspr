package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	celestialBlue  = lipgloss.Color("#2563EB")
	spaceNavy      = lipgloss.Color("#1E3A8A")
	starlightBlue  = lipgloss.Color("#60A5FA")
	silverWhite    = lipgloss.Color("#F3F4F6")
	cosmicSilver   = lipgloss.Color("#D1D5DB")
	nebulaGray     = lipgloss.Color("#9CA3AF")
	deepSpace      = lipgloss.Color("#0B1120")
	midnightNebula = lipgloss.Color("#111827")
	deepGray       = lipgloss.Color("#4B5563")
	slateGray      = lipgloss.Color("#1E293B")
	darkBlueGray   = lipgloss.Color("#0F172A")
)

// ThemeForProvider returns the Styles associated with the given provider
// ID. Unknown or empty provider IDs yield the default Charmtone Pantera
// theme.
func ThemeForProvider(providerID string) Styles {
	switch providerID {
	case "hyper":
		return HypercrushObsidiana()
	default:
		return CharmtonePantera()
	}
}

// CharmtonePantera returns the Charmtone dark theme. It's the default style
// for the UI.
func CharmtonePantera() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   celestialBlue,
		secondary: starlightBlue,
		accent:    spaceNavy,
		keyword:   cosmicSilver,

		fgBase:       silverWhite,
		fgMoreSubtle: nebulaGray,
		fgSubtle:     cosmicSilver,
		fgMostSubtle: deepGray,

		onPrimary: deepSpace,

		bgBase:         deepSpace,
		bgLeastVisible: midnightNebula,
		bgLessVisible:  darkBlueGray,
		bgMostVisible:  spaceNavy,

		separator: slateGray,

		destructive:       lipgloss.Color("#EF4444"),
		error:             lipgloss.Color("#DC2626"),
		warningSubtle:     lipgloss.Color("#FBBF24"),
		warning:           lipgloss.Color("#F59E0B"),
		denied:            lipgloss.Color("#F97316"),
		busy:              lipgloss.Color("#EAB308"),
		info:              starlightBlue,
		infoMoreSubtle:    celestialBlue,
		infoMostSubtle:    spaceNavy,
		success:           lipgloss.Color("#10B981"),
		successMoreSubtle: lipgloss.Color("#047857"),
		successMostSubtle: lipgloss.Color("#065F46"),
	})

	// Bang ! prompt overrides - use Celestial colors.
	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(silverWhite).
		Background(celestialBlue)
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(celestialBlue)
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(spaceNavy)

	// Shell bar/prompt overrides - use Celestial colors.
	s.Messages.ShellBarFocused = s.Messages.ShellBarFocused.
		BorderForeground(celestialBlue)
	s.Messages.ShellBarBlurred = s.Messages.ShellBarBlurred.
		BorderForeground(cosmicSilver)
	s.Messages.ShellPrompt = s.Messages.ShellPrompt.
		Foreground(starlightBlue)
	s.Messages.ShellPromptBlurred = s.Messages.ShellPromptBlurred.
		Foreground(starlightBlue)

	return s
}

// HypercrushObsidiana returns the Hypercrush dark theme.
func HypercrushObsidiana() Styles {
	return CharmtonePantera()
}
