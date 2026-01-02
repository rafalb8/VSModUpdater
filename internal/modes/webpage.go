package modes

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rafalb8/VSModUpdater/internal/config"
	"github.com/rafalb8/VSModUpdater/internal/mod"
	"github.com/rafalb8/VSModUpdater/internal/webpage"
)

func Webpage() {
	fmt.Printf("Loading mods from: %s\n", config.ModPath)
	
	// Load mods
	mods, err := mod.InfoFromPath(config.ModPath)
	if err != nil {
		fmt.Printf("Error loading mods: %v\n", err)
		return
	}
	
	fmt.Printf("Loaded %d mods\n", len(mods))
	
	// Generate HTML
	opts := webpage.Options{
		Title: config.WebpageTitle,
	}
	
	html := webpage.Generate(mods, opts)
	
	// Write to file
	err = os.WriteFile(config.WebpageOutput, []byte(html), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Generated webpage: %s\n", config.WebpageOutput)
	
	// Deploy if requested
	if config.WebpageDeploy != "" {
		deployToCloudflare(config.WebpageOutput, config.WebpageDeploy)
	}
}

func deployToCloudflare(filename, project string) {
	fmt.Printf("\nDeploying to Cloudflare Pages (project: %s)...\n", project)
	
	// Check if wrangler is installed
	_, err := exec.LookPath("wrangler")
	if err != nil {
		fmt.Println("Error: wrangler is not installed")
		fmt.Println("\nTo install wrangler:")
		fmt.Println("  npm install -g wrangler")
		fmt.Println("\nThen login:")
		fmt.Println("  wrangler login")
		return
	}
	
	// Create temporary directory for deployment
	tmpDir, err := os.MkdirTemp("", "vsmodupdater-deploy-*")
	if err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)
	
	// Copy HTML file to temp directory as index.html
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading HTML file: %v\n", err)
		return
	}
	
	indexPath := filepath.Join(tmpDir, "index.html")
	err = os.WriteFile(indexPath, content, 0644)
	if err != nil {
		fmt.Printf("Error writing to temp directory: %v\n", err)
		return
	}
	
	// Try to deploy
	cmd := exec.Command("wrangler", "pages", "deploy", tmpDir, "--project-name", project)
	var stderr strings.Builder
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	
	err = cmd.Run()
	if err != nil {
		// Check if project doesn't exist
		if strings.Contains(stderr.String(), "Project not found") || strings.Contains(stderr.String(), "does not match any of your existing projects") {
			fmt.Printf("\n⚠ Project '%s' does not exist.\n", project)
			fmt.Print("Would you like to create it? (y/N): ")
			
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				response := strings.TrimSpace(strings.ToLower(scanner.Text()))
				if response == "y" || response == "yes" {
					createProject(tmpDir, project)
					return
				}
			}
			fmt.Println("Deployment cancelled.")
			return
		}
		
		fmt.Printf("\nDeployment failed: %v\n", err)
		fmt.Println("\nMake sure you've logged in with:")
		fmt.Println("  wrangler login")
		return
	}
	
	fmt.Printf("\n✓ Successfully deployed to Cloudflare Pages!\n")
	fmt.Printf("  Visit: https://%s.pages.dev\n", project)
}

func createProject(deployDir, project string) {
	fmt.Printf("\nCreating project '%s'...\n", project)
	
	// Create project using wrangler pages project create
	cmd := exec.Command("wrangler", "pages", "project", "create", project, "--production-branch", "main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\nFailed to create project: %v\n", err)
		return
	}
	
	fmt.Println("\n✓ Project created successfully!")
	fmt.Println("\nNow deploying...")
	
	// Deploy to the newly created project
	deployCmd := exec.Command("wrangler", "pages", "deploy", deployDir, "--project-name", project)
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	
	err = deployCmd.Run()
	if err != nil {
		fmt.Printf("\nDeployment failed: %v\n", err)
		return
	}
	
	fmt.Printf("\n✓ Successfully deployed to Cloudflare Pages!\n")
	fmt.Printf("  Visit: https://%s.pages.dev\n", project)
}
