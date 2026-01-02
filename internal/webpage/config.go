package webpage

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/mod"
)

// Config represents the webpage configuration
type Config struct {
	Title       string               `json:"title,omitempty"`
	ProjectName string               `json:"projectName,omitempty"`
	Categories  map[string]Category  `json:"categories"`
	Mods        map[string]ModConfig `json:"mods"`
}

// Category represents a category with name and color
type Category struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ModConfig represents custom configuration for a mod
type ModConfig struct {
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// LoadConfig loads the webpage configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty config if file doesn't exist
			return &Config{
				Categories: make(map[string]Category),
				Mods:       make(map[string]ModConfig),
			}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Initialize maps if nil
	if config.Categories == nil {
		config.Categories = make(map[string]Category)
	}
	if config.Mods == nil {
		config.Mods = make(map[string]ModConfig)
	}

	return &config, nil
}

// SaveConfig saves the webpage configuration to a JSON file
func (c *Config) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GetDescription returns the custom description for a mod, or the default if not set
func (c *Config) GetDescription(modInfo *mod.Info) string {
	if cfg, ok := c.Mods[modInfo.ModID]; ok && cfg.Description != "" {
		return cfg.Description
	}
	return modInfo.Description
}

// GetCategory returns the category for a mod, or empty string if not set
func (c *Config) GetCategory(modID string) string {
	if cfg, ok := c.Mods[modID]; ok {
		return cfg.Category
	}
	return ""
}

// SetModConfig sets the configuration for a specific mod
func (c *Config) SetModConfig(modID string, description, category string) {
	if c.Mods == nil {
		c.Mods = make(map[string]ModConfig)
	}

	cfg := c.Mods[modID]
	if description != "" {
		cfg.Description = description
	}
	if category != "" {
		cfg.Category = category
	}
	c.Mods[modID] = cfg
}

// RemoveModConfig removes custom configuration for a mod
func (c *Config) RemoveModConfig(modID string) {
	delete(c.Mods, modID)
}

// AddCategory adds a new category
func (c *Config) AddCategory(name, color string) error {
	if c.Categories == nil {
		c.Categories = make(map[string]Category)
	}

	// Check if category already exists
	if _, exists := c.Categories[name]; exists {
		return fmt.Errorf("category '%s' already exists", name)
	}

	c.Categories[name] = Category{
		Name:  name,
		Color: color,
	}
	return nil
}

// EditCategory edits an existing category
func (c *Config) EditCategory(oldName, newName, color string) error {
	cat, exists := c.Categories[oldName]
	if !exists {
		return fmt.Errorf("category '%s' does not exist", oldName)
	}

	// If name is changing, update all mods using this category
	if oldName != newName {
		for modID, modCfg := range c.Mods {
			if modCfg.Category == oldName {
				modCfg.Category = newName
				c.Mods[modID] = modCfg
			}
		}
		delete(c.Categories, oldName)
	}

	cat.Name = newName
	if color != "" {
		cat.Color = color
	}
	c.Categories[newName] = cat
	return nil
}

// DeleteCategory deletes a category
func (c *Config) DeleteCategory(name string) error {
	if _, exists := c.Categories[name]; !exists {
		return fmt.Errorf("category '%s' does not exist", name)
	}

	// Check if any mods use this category
	usedBy := []string{}
	for modID, modCfg := range c.Mods {
		if modCfg.Category == name {
			usedBy = append(usedBy, modID)
		}
	}

	if len(usedBy) > 0 {
		return fmt.Errorf("category '%s' is used by %d mod(s): %s", name, len(usedBy), strings.Join(usedBy, ", "))
	}

	delete(c.Categories, name)
	return nil
}

// ListCategories returns a sorted list of category names
func (c *Config) ListCategories() []string {
	names := make([]string, 0, len(c.Categories))
	for name := range c.Categories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetCategoryColor returns the color for a category, or empty string if not found
func (c *Config) GetCategoryColor(name string) string {
	if cat, ok := c.Categories[name]; ok {
		return cat.Color
	}
	return ""
}

// GetCategoryName returns the display name for a category, or the key if not found
func (c *Config) GetCategoryName(key string) string {
	if cat, ok := c.Categories[key]; ok {
		return cat.Name
	}
	return key
}
