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

// ThemeForConfig returns the Styles associated with the given theme name.
// If the theme is empty or "default" (or unknown), it falls back to the provider theme.
func ThemeForConfig(themeName string, providerID string) Styles {
	switch themeName {
	case "light":
		return CharmtoneLight()
	case "nord":
		return NordTheme()
	case "dracula":
		return DraculaTheme()
	case "cyberpunk", "neon":
		return CyberpunkTheme()
	default:
		return ThemeForProvider(providerID)
	}
}

// CharmtoneLight returns the Charmtone light theme.
func CharmtoneLight() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   celestialBlue,
		secondary: spaceNavy,
		accent:    starlightBlue,
		keyword:   deepGray,

		fgBase:       lipgloss.Color("#111827"),
		fgMoreSubtle: lipgloss.Color("#4B5563"),
		fgSubtle:     lipgloss.Color("#374151"),
		fgMostSubtle: lipgloss.Color("#9CA3AF"),

		onPrimary: lipgloss.Color("#FFFFFF"),

		bgBase:         lipgloss.Color("#FFFFFF"),
		bgLeastVisible: lipgloss.Color("#F9FAFB"),
		bgLessVisible:  lipgloss.Color("#F3F4F6"),
		bgMostVisible:  lipgloss.Color("#E5E7EB"),

		separator: lipgloss.Color("#E5E7EB"),

		destructive:       lipgloss.Color("#EF4444"),
		error:             lipgloss.Color("#DC2626"),
		warningSubtle:     lipgloss.Color("#FEF3C7"),
		warning:           lipgloss.Color("#D97706"),
		denied:            lipgloss.Color("#EA580C"),
		busy:              lipgloss.Color("#CA8A04"),
		info:              starlightBlue,
		infoMoreSubtle:    celestialBlue,
		infoMostSubtle:    spaceNavy,
		success:           lipgloss.Color("#10B981"),
		successMoreSubtle: lipgloss.Color("#059669"),
		successMostSubtle: lipgloss.Color("#A7F3D0"),
	})

	// Bang prompt overrides
	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(celestialBlue)
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(celestialBlue)
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(spaceNavy)

	// Shell bar/prompt overrides
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

// NordTheme returns the Nord style theme.
func NordTheme() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   lipgloss.Color("#88C0D0"),
		secondary: lipgloss.Color("#81A1C1"),
		accent:    lipgloss.Color("#5E81AC"),
		keyword:   lipgloss.Color("#D8DEE9"),

		fgBase:       lipgloss.Color("#E5E9F0"),
		fgMoreSubtle: lipgloss.Color("#4C566A"),
		fgSubtle:     lipgloss.Color("#D8DEE9"),
		fgMostSubtle: lipgloss.Color("#434C5E"),

		onPrimary: lipgloss.Color("#2E3440"),

		bgBase:         lipgloss.Color("#2E3440"),
		bgLeastVisible: lipgloss.Color("#3B4252"),
		bgLessVisible:  lipgloss.Color("#434C5E"),
		bgMostVisible:  lipgloss.Color("#4C566A"),

		separator: lipgloss.Color("#3B4252"),

		destructive:       lipgloss.Color("#BF616A"),
		error:             lipgloss.Color("#BF616A"),
		warningSubtle:     lipgloss.Color("#EBCB8B"),
		warning:           lipgloss.Color("#EBCB8B"),
		denied:            lipgloss.Color("#D08770"),
		busy:              lipgloss.Color("#EBCB8B"),
		info:              lipgloss.Color("#88C0D0"),
		infoMoreSubtle:    lipgloss.Color("#81A1C1"),
		infoMostSubtle:    lipgloss.Color("#5E81AC"),
		success:           lipgloss.Color("#A3BE8C"),
		successMoreSubtle: lipgloss.Color("#A3BE8C"),
		successMostSubtle: lipgloss.Color("#A3BE8C"),
	})

	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(lipgloss.Color("#2E3440")).
		Background(lipgloss.Color("#88C0D0"))
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(lipgloss.Color("#88C0D0"))
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(lipgloss.Color("#4C566A"))

	return s
}

// DraculaTheme returns the Dracula style theme.
func DraculaTheme() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   lipgloss.Color("#BD93F9"),
		secondary: lipgloss.Color("#FF79C6"),
		accent:    lipgloss.Color("#50FA7B"),
		keyword:   lipgloss.Color("#F8F8F2"),

		fgBase:       lipgloss.Color("#F8F8F2"),
		fgMoreSubtle: lipgloss.Color("#6272A4"),
		fgSubtle:     lipgloss.Color("#F8F8F2"),
		fgMostSubtle: lipgloss.Color("#44475A"),

		onPrimary: lipgloss.Color("#282A36"),

		bgBase:         lipgloss.Color("#282A36"),
		bgLeastVisible: lipgloss.Color("#1D1F27"),
		bgLessVisible:  lipgloss.Color("#44475A"),
		bgMostVisible:  lipgloss.Color("#6272A4"),

		separator: lipgloss.Color("#44475A"),

		destructive:       lipgloss.Color("#FF5555"),
		error:             lipgloss.Color("#FF5555"),
		warningSubtle:     lipgloss.Color("#F1FA8C"),
		warning:           lipgloss.Color("#FFB86C"),
		denied:            lipgloss.Color("#FFB86C"),
		busy:              lipgloss.Color("#F1FA8C"),
		info:              lipgloss.Color("#8BE9FD"),
		infoMoreSubtle:    lipgloss.Color("#BD93F9"),
		infoMostSubtle:    lipgloss.Color("#6272A4"),
		success:           lipgloss.Color("#50FA7B"),
		successMoreSubtle: lipgloss.Color("#50FA7B"),
		successMostSubtle: lipgloss.Color("#50FA7B"),
	})

	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(lipgloss.Color("#282A36")).
		Background(lipgloss.Color("#BD93F9"))
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(lipgloss.Color("#BD93F9"))
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(lipgloss.Color("#6272A4"))

	return s
}

// CyberpunkTheme returns the Cyberpunk/Neon style theme.
func CyberpunkTheme() Styles {
	s := quickStyle(quickStyleOpts{
		primary:   lipgloss.Color("#00F0FF"),
		secondary: lipgloss.Color("#FF007F"),
		accent:    lipgloss.Color("#FFE600"),
		keyword:   lipgloss.Color("#FFFFFF"),

		fgBase:       lipgloss.Color("#FFFFFF"),
		fgMoreSubtle: lipgloss.Color("#55557F"),
		fgSubtle:     lipgloss.Color("#CCCCCC"),
		fgMostSubtle: lipgloss.Color("#333344"),

		onPrimary: lipgloss.Color("#0D0E15"),

		bgBase:         lipgloss.Color("#0D0E15"),
		bgLeastVisible: lipgloss.Color("#1C1D2B"),
		bgLessVisible:  lipgloss.Color("#33334F"),
		bgMostVisible:  lipgloss.Color("#FF007F"),

		separator: lipgloss.Color("#33334F"),

		destructive:       lipgloss.Color("#FF0000"),
		error:             lipgloss.Color("#FF0033"),
		warningSubtle:     lipgloss.Color("#FFE600"),
		warning:           lipgloss.Color("#FFE600"),
		denied:            lipgloss.Color("#FF5500"),
		busy:              lipgloss.Color("#FFE600"),
		info:              lipgloss.Color("#00F0FF"),
		infoMoreSubtle:    lipgloss.Color("#FF007F"),
		infoMostSubtle:    lipgloss.Color("#33334F"),
		success:           lipgloss.Color("#00FF66"),
		successMoreSubtle: lipgloss.Color("#00FF66"),
		successMostSubtle: lipgloss.Color("#00FF66"),
	})

	s.Editor.PromptBangIconFocused = s.Editor.PromptBangIconFocused.
		Foreground(lipgloss.Color("#0D0E15")).
		Background(lipgloss.Color("#00F0FF"))
	s.Editor.PromptBangDotsFocused = s.Editor.PromptBangDotsFocused.
		Foreground(lipgloss.Color("#00F0FF"))
	s.Editor.PromptBangDotsBlurred = s.Editor.PromptBangDotsBlurred.
		Foreground(lipgloss.Color("#FF007F"))

	return s
}
