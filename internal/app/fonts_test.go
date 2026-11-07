package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"

	"mygo-music/internal/i18n"
)

// Choosing a font changes text in both windows, survives a restart, and can
// be reset to the system's with the keyboard.
func TestChooseFont(t *testing.T) {
	a := newTestApp()
	a.fontFamilies = []string{"monospace"}
	a.settings.path = filepath.Join(t.TempDir(), "settings.json")
	var settingsWidth float32
	settings := ui.NewTester(func(c *ui.Context) {
		a.settingsView(c)
		settingsWidth, _ = c.MeasureText(0, ui.Span{Text: "General"})
	}, 780, 640)
	main := ui.NewTester(a.view, 1180, 780)
	settings.Frame()
	main.Frame()
	beforeSettings := settingsWidth
	beforeMain := find(t, main, "Home").W

	// Follow System also labels the language choice: use the font row's
	// control at its right edge.
	row := find(t, settings, "Font")
	settings.ClickAt(780-28-16-85, row.Y+row.H/2)
	if err := settings.Click("monospace"); err != nil {
		t.Fatal(err)
	}
	settings.Frame()
	main.Frame()
	if a.settings.Font != "monospace" {
		t.Fatalf("chose %q", a.settings.Font)
	}
	if settingsWidth == beforeSettings || find(t, main, "Home").W == beforeMain {
		t.Fatal("the font did not change text layout in both windows")
	}
	b, err := os.ReadFile(a.settings.path)
	if err != nil {
		t.Fatal(err)
	}
	restored := defaultSettings()
	if err := json.Unmarshal(b, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Font != "monospace" {
		t.Fatalf("restored %q", restored.Font)
	}
	if err := settings.Click("monospace"); err != nil {
		t.Fatal(err)
	}
	settings.Key(0, ui.KeyHome)
	settings.Key(0, ui.KeyEnter)
	settings.Frame()
	main.Frame()
	if a.settings.Font != "" || settingsWidth != beforeSettings || find(t, main, "Home").W != beforeMain {
		t.Fatal("reset did not restore the system font in both windows")
	}
	b, err = os.ReadFile(a.settings.path)
	if err != nil {
		t.Fatal(err)
	}
	// The empty choice is omitted, so a fresh settings value is used on load.
	restored = defaultSettings()
	if err := json.Unmarshal(b, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.Font != "" {
		t.Fatalf("reset kept %q", restored.Font)
	}
}

// A font uninstalled since the last run stays chosen until the listener
// selects another; opening and cancelling its picker must not overwrite it.
func TestMissingFont(t *testing.T) {
	a := newTestApp()
	a.settings.Font = "MyGo Music Missing Font"
	a.fontFamilies = []string{"monospace"}
	tt := ui.NewTester(a.settingsView, 780, 640)
	tt.Frame()
	if err := tt.Click(a.settings.Font); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if a.settings.Font != "MyGo Music Missing Font" {
		t.Fatalf("opening the picker replaced the missing font with %q", a.settings.Font)
	}
	if err := tt.Click(a.settings.Font); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("monospace"); err != nil {
		t.Fatal(err)
	}
	if a.settings.Font != "monospace" {
		t.Fatalf("could not replace the missing font: %q", a.settings.Font)
	}
}

// With MM_FRAMES set, saves the font picker in both languages and appearances.
func TestFontSettingsFrames(t *testing.T) {
	dir := framesDir(t)
	defer i18n.SetLanguage(i18n.English)
	for _, language := range []string{i18n.English, i18n.Chinese} {
		i18n.SetLanguage(language)
		for _, dark := range []bool{false, true} {
			a := newTestApp()
			ch := posts(a)
			tt := ui.NewTester(a.settingsView, 780, 640)
			tt.SetDark(dark)
			tt.Frame()
			a.loadFontFamilies()
			runPost(t, ch)
			settled(tt)
			name := "font-settings-" + language
			if dark {
				name += "-dark"
			}
			saveFrame(t, dir, name+".png", tt.Image())
			row := find(t, tt, T("Font"))
			tt.ClickAt(780-28-16-85, row.Y+row.H/2)
			settled(tt)
			saveFrame(t, dir, name+"-open.png", tt.Image())
		}
	}
}
