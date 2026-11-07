package app

import (
	"slices"

	"github.com/egoist/mygo/ui"
)

// loadFontFamilies asks for installed fonts once, without holding up the
// settings window while the system lists them.
func (a *app) loadFontFamilies() {
	if a.fontFamilies != nil || a.fontsLoading {
		return
	}
	a.fontsLoading = true
	go func() {
		families := ui.FontFamilies()
		names := make([]string, len(families))
		for i, f := range families {
			names[i] = f.Name
		}
		a.post(func() {
			a.fontFamilies = names
			a.fontsLoading = false
		})
	}()
}

// fontChoice keeps a family by its name, including one no longer installed:
// text falls back to the system's until the listener chooses another.
func (a *app) fontChoice(c *ui.Context) {
	follow := T("Follow System")
	names := append([]string{follow}, a.fontFamilies...)
	chosen := a.settings.Font
	if chosen == "" {
		chosen = follow
	} else if !slices.Contains(a.fontFamilies, chosen) {
		names = append(names, chosen)
	}
	if ui.Select(c.Key("font"), &chosen, names).Label(T("Font")).Width(170).Changed() {
		if chosen == follow {
			chosen = ""
		}
		a.settings.Font = chosen
		a.changed()
	}
}
