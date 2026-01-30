package styles

import "testing"

func TestFontSizeConfig_Small(t *testing.T) {
	cfg := GetFontSizeConfig(1)
	if cfg.Level != 1 {
		t.Errorf("expected level 1, got %d", cfg.Level)
	}
	if cfg.VersePadding != 0 {
		t.Errorf("expected VersePadding 0, got %d", cfg.VersePadding)
	}
	if cfg.VerseIndent != 2 {
		t.Errorf("expected VerseIndent 2, got %d", cfg.VerseIndent)
	}
}

func TestFontSizeConfig_Medium(t *testing.T) {
	cfg := GetFontSizeConfig(2)
	if cfg.Level != 2 {
		t.Errorf("expected level 2, got %d", cfg.Level)
	}
	if cfg.NumberWidth != 4 {
		t.Errorf("expected NumberWidth 4, got %d", cfg.NumberWidth)
	}
	if cfg.SectionGap != 1 {
		t.Errorf("expected SectionGap 1, got %d", cfg.SectionGap)
	}
}

func TestFontSizeConfig_Large(t *testing.T) {
	cfg := GetFontSizeConfig(3)
	if cfg.Level != 3 {
		t.Errorf("expected level 3, got %d", cfg.Level)
	}
	if cfg.VersePadding != 1 {
		t.Errorf("expected VersePadding 1, got %d", cfg.VersePadding)
	}
	if cfg.SectionGap != 2 {
		t.Errorf("expected SectionGap 2, got %d", cfg.SectionGap)
	}
}

func TestFontSizeBounds_Low(t *testing.T) {
	cfg := GetFontSizeConfig(0)
	if cfg.Level != 1 {
		t.Errorf("expected clamped to 1, got %d", cfg.Level)
	}
}

func TestFontSizeBounds_High(t *testing.T) {
	cfg := GetFontSizeConfig(4)
	if cfg.Level != 3 {
		t.Errorf("expected clamped to 3, got %d", cfg.Level)
	}
}

func TestFontSizeBounds_Negative(t *testing.T) {
	cfg := GetFontSizeConfig(-1)
	if cfg.Level != 1 {
		t.Errorf("expected clamped to 1, got %d", cfg.Level)
	}
}
