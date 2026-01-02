package modes

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
	"github.com/rafalb8/VSModUpdater/internal/webpage"
)

// ConfigManagement handles configuration management for webpage
func ConfigManagement() error {
	cfg, err := webpage.LoadConfig(config.WebConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Handle category operations
	if config.ConfigAddCat != "" {
		return handleAddCategory(cfg, config.ConfigAddCat)
	}
	if config.ConfigEditCat != "" {
		return handleEditCategory(cfg, config.ConfigEditCat)
	}
	if config.ConfigDelCat != "" {
		return handleDeleteCategory(cfg, config.ConfigDelCat)
	}

	// Handle single mod edit
	if config.ConfigModID != "" {
		return handleEditMod(cfg, config.ConfigModID)
	}

	// Interactive mode for all mods
	return handleInteractiveConfig(cfg)
}

func handleAddCategory(cfg *webpage.Config, input string) error {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format. Use: -add-category='name:color'")
	}

	name := strings.TrimSpace(parts[0])
	color := strings.TrimSpace(parts[1])

	if err := cfg.AddCategory(name, color); err != nil {
		return err
	}

	if err := cfg.SaveConfig(config.WebConfigFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Added category '%s' with color %s\n", name, color)
	return nil
}

func handleEditCategory(cfg *webpage.Config, input string) error {
	parts := strings.SplitN(input, ":", 3)
	if len(parts) < 2 {
		return fmt.Errorf("invalid format. Use: -edit-category='oldname:newname:color' (color is optional)")
	}

	oldName := strings.TrimSpace(parts[0])
	newName := strings.TrimSpace(parts[1])
	color := ""
	if len(parts) == 3 {
		color = strings.TrimSpace(parts[2])
	}

	if err := cfg.EditCategory(oldName, newName, color); err != nil {
		return err
	}

	if err := cfg.SaveConfig(config.WebConfigFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Updated category '%s' to '%s'\n", oldName, newName)
	return nil
}

func handleDeleteCategory(cfg *webpage.Config, name string) error {
	if err := cfg.DeleteCategory(name); err != nil {
		return err
	}

	if err := cfg.SaveConfig(config.WebConfigFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Deleted category '%s'\n", name)
	return nil
}

func handleEditMod(cfg *webpage.Config, modID string) error {
	// Load mods to get current info
	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		return fmt.Errorf("failed to load mods: %w", err)
	}

	// Find the mod
	var modInfo *mod.Info
	for _, m := range mods {
		if m.ModID == modID {
			modInfo = m
			break
		}
	}

	if modInfo == nil {
		return fmt.Errorf("mod '%s' not found in %s", modID, config.ModPath)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("\n=== Editing Mod: %s ===\n", modInfo.Name)
	fmt.Printf("Current Description: %s\n", cfg.GetDescription(modInfo))
	fmt.Printf("Current Category: %s\n\n", cfg.GetCategory(modID))

	// Get new description
	fmt.Print("New description (leave empty to keep current, or type 'r' to reset to default): ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	if description == "r" {
		// Clear custom description
		if modCfg, ok := cfg.Mods[modID]; ok {
			modCfg.Description = ""
			cfg.Mods[modID] = modCfg
		}
		description = ""
	}

	// Get category
	cats := cfg.ListCategories()
	fmt.Print("\nCategories: ")
	for i, cat := range cats {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d=%s", i+1, cat)
	}
	fmt.Print("\nCategory (enter number, name, or leave empty to keep current, or type 'n' to remove): ")
	scanner.Scan()
	categoryInput := strings.TrimSpace(scanner.Text())

	category := ""
	if categoryInput == "n" {
		// Clear category
		if modCfg, ok := cfg.Mods[modID]; ok {
			modCfg.Category = ""
			cfg.Mods[modID] = modCfg
		}
	} else if categoryInput != "" {
		// Try to parse as number
		var catNum int
		if _, err := fmt.Sscanf(categoryInput, "%d", &catNum); err == nil {
			if catNum > 0 && catNum <= len(cats) {
				category = cats[catNum-1]
			} else {
				return fmt.Errorf("invalid category number")
			}
		} else {
			// Use as category name
			category = categoryInput
			// Verify it exists
			if _, ok := cfg.Categories[category]; !ok {
				return fmt.Errorf("category '%s' does not exist", category)
			}
		}
	}

	// Update config
	if description != "" || category != "" {
		cfg.SetModConfig(modID, description, category)
	}

	if err := cfg.SaveConfig(config.WebConfigFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\nUpdated configuration for %s\n", modInfo.Name)
	return nil
}

func handleInteractiveConfig(cfg *webpage.Config) error {
	// Load all installed mods
	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		return fmt.Errorf("failed to load mods: %w", err)
	}

	// Sort mods by name for easier navigation
	sort.Slice(mods, func(i, j int) bool {
		return strings.ToLower(mods[i].Name) < strings.ToLower(mods[j].Name)
	})

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== Interactive Mod Configuration ===")
	fmt.Println("For each mod, you can set a custom description and category.")
	fmt.Println("Press Enter to skip, type 'quit' to exit.\n")

	// Get categories list
	cats := cfg.ListCategories()

	for i, modInfo := range mods {
		fmt.Printf("\n[%d/%d] %s\n", i+1, len(mods), modInfo.Name)
		fmt.Printf("ModID: %s\n", modInfo.ModID)
		fmt.Printf("Default Description: %s\n", modInfo.Description)

		currentDesc := cfg.GetDescription(modInfo)
		currentCat := cfg.GetCategory(modInfo.ModID)

		if currentDesc != modInfo.Description {
			fmt.Printf("Current Custom Description: %s\n", currentDesc)
		}
		if currentCat != "" {
			fmt.Printf("Current Category: %s\n", currentCat)
		}

		// Get description
		fmt.Print("\nCustom description (or 'r' to reset to default): ")
		scanner.Scan()
		description := strings.TrimSpace(scanner.Text())

		if description == "quit" {
			break
		}

		if description == "r" {
			// Clear custom description
			if modCfg, ok := cfg.Mods[modInfo.ModID]; ok {
				modCfg.Description = ""
				cfg.Mods[modInfo.ModID] = modCfg
			}
			description = ""
		}

		// Get category
		fmt.Print("\nCategories: ")
		for i, cat := range cats {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d=%s", i+1, cat)
		}
		fmt.Print("\nCategory (number/name, or 'n' to remove): ")
		scanner.Scan()
		categoryInput := strings.TrimSpace(scanner.Text())

		if categoryInput == "quit" {
			break
		}

		category := ""
		if categoryInput == "n" {
			// Clear category
			if modCfg, ok := cfg.Mods[modInfo.ModID]; ok {
				modCfg.Category = ""
				cfg.Mods[modInfo.ModID] = modCfg
			}
		} else if categoryInput != "" {
			// Try to parse as number
			var catNum int
			if _, err := fmt.Sscanf(categoryInput, "%d", &catNum); err == nil {
				if catNum > 0 && catNum <= len(cats) {
					category = cats[catNum-1]
				} else {
					fmt.Println("Invalid category number, skipping...")
					continue
				}
			} else {
				// Use as category name
				category = categoryInput
				// Verify it exists
				if _, ok := cfg.Categories[category]; !ok {
					fmt.Printf("Warning: Category '%s' does not exist. Skipping...\n", category)
					continue
				}
			}
		}

		// Update config
		if description != "" || category != "" {
			cfg.SetModConfig(modInfo.ModID, description, category)
		}

		// Save after each mod to avoid losing progress
		if err := cfg.SaveConfig(config.WebConfigFile); err != nil {
			fmt.Printf("Warning: Failed to save config: %v\n", err)
		}
	}

	fmt.Println("\n=== Configuration Complete ===")
	fmt.Printf("Config saved to: %s\n", config.WebConfigFile)
	return nil
}
