package webpage

import (
	_ "embed"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/rafalb8/VSModUpdater/internal/mod"
)

//go:embed templates/style.css
var cssTemplate string

//go:embed templates/script.js
var jsTemplate string

//go:embed templates/page.html
var htmlTemplate string

// Options for webpage generation
type Options struct {
	Title      string
	ServerName string
}

// ModData represents a mod for template rendering
type ModData struct {
	Name        string
	Version     string
	Authors     string
	Description string
	GameVersion string
	URL         string
}

// PageData represents the full page data for template rendering
type PageData struct {
	Title      string
	ModCount   int
	UpdatedAt  string
	Mods       []ModData
	CSS        template.CSS
	JavaScript template.JS
}

// Generate creates a self-contained HTML page with mod information
func Generate(mods []*mod.Info, opts Options) string {
	// Sort mods by name
	sort.Slice(mods, func(i, j int) bool {
		return strings.ToLower(mods[i].Name) < strings.ToLower(mods[j].Name)
	})

	// Filter out mods with errors and convert to ModData
	modDataList := make([]ModData, 0, len(mods))
	for _, m := range mods {
		if m.Error != nil {
			continue
		}

		authors := strings.Join(m.Authors, ", ")
		if authors == "" {
			authors = "Unknown"
		}

		desc := m.Description
		if len(desc) > 200 {
			desc = desc[:197] + "..."
		}
		if desc == "" {
			desc = "No description available"
		}

		gameVer := "any"
		if gv, ok := m.Dependencies["game"]; ok && gv != "" && gv != "*" {
			gameVer = gv
		}

		modDataList = append(modDataList, ModData{
			Name:        m.Name,
			Version:     string(m.Version),
			Authors:     authors,
			Description: desc,
			GameVersion: gameVer,
			URL:         getModURL(m.ModID),
		})
	}

	// Prepare page data
	pageData := PageData{
		Title:      opts.Title,
		ModCount:   len(modDataList),
		UpdatedAt:  time.Now().Format("January 2, 2006 at 3:04 PM MST"),
		Mods:       modDataList,
		CSS:        template.CSS(cssTemplate),
		JavaScript: template.JS(jsTemplate),
	}

	// Parse and execute template
	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, pageData)
	if err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}

func getModURL(modID string) string {
	if modID == "" {
		return ""
	}
	return fmt.Sprintf("https://mods.vintagestory.at/%s", modID)
}
