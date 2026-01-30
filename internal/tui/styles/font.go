package styles

// FontSizeConfig defines layout density for different "font size" levels.
// Terminal apps can't change actual font size, so we adjust spacing and padding.
type FontSizeConfig struct {
	Level       int // 1 (small/compact), 2 (medium/default), 3 (large/spacious)
	VersePadding int // blank lines between verses (0, 0, 1)
	VerseIndent  int // text indent in spaces (2, 3, 4)
	NumberWidth  int // verse number display width in chars (3, 4, 5)
	SectionGap   int // blank lines before/after section titles (0, 1, 2)
}

// GetFontSizeConfig returns layout config for the given level.
// Out-of-range values are clamped to [1, 3].
func GetFontSizeConfig(level int) *FontSizeConfig {
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}
	switch level {
	case 1:
		return &FontSizeConfig{Level: 1, VersePadding: 0, VerseIndent: 2, NumberWidth: 3, SectionGap: 0}
	case 3:
		return &FontSizeConfig{Level: 3, VersePadding: 1, VerseIndent: 4, NumberWidth: 5, SectionGap: 2}
	default: // 2
		return &FontSizeConfig{Level: 2, VersePadding: 0, VerseIndent: 3, NumberWidth: 4, SectionGap: 1}
	}
}
