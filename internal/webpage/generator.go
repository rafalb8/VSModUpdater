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
	Name          string
	Version       string
	Authors       string
	Description   string
	GameVersion   string
	URL           string
	Category      string
	CategoryColor string
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
	return GenerateWithConfig(mods, opts, nil)
}

// GenerateWithConfig creates a self-contained HTML page with mod information using a config file
func GenerateWithConfig(mods []*mod.Info, opts Options, cfg *Config) string {
	// First pass: collect mod data with categories
	modDataList := make([]ModData, 0, len(mods))
	for _, m := range mods {
		if m.Error != nil {
			continue
		}
		if m.Error != nil {
			continue
		}

		authors := strings.Join(m.Authors, ", ")
		if authors == "" {
			authors = "Unknown"
		}

		// Use custom description if available
		desc := m.Description
		if cfg != nil {
			desc = cfg.GetDescription(m)
		}
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

		// Get category info
		category := ""
		categoryColor := ""
		if cfg != nil {
			categoryKey := cfg.GetCategory(m.ModID)
			if categoryKey != "" {
				category = cfg.GetCategoryName(categoryKey)
				categoryColor = cfg.GetCategoryColor(categoryKey)
			}
		}

		modDataList = append(modDataList, ModData{
			Name:          m.Name,
			Version:       string(m.Version),
			Authors:       authors,
			Description:   desc,
			GameVersion:   gameVer,
			URL:           getModURL(m.ModID),
			Category:      category,
			CategoryColor: categoryColor,
		})
	}

	// Sort by category first, then by name
	sort.Slice(modDataList, func(i, j int) bool {
		// If categories are different, sort by category
		if modDataList[i].Category != modDataList[j].Category {
			// Empty categories go to the end
			if modDataList[i].Category == "" {
				return false
			}
			if modDataList[j].Category == "" {
				return true
			}
			return strings.ToLower(modDataList[i].Category) < strings.ToLower(modDataList[j].Category)
		}
		// If categories are the same, sort by name
		return strings.ToLower(modDataList[i].Name) < strings.ToLower(modDataList[j].Name)
	})

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
