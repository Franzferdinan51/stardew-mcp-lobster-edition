package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rivo/tview"
)

// Configuration
var (
	app           *tview.Application
	pages         *tview.Pages
	logView       *tview.TextView
	stardewPath   string
	openclawEnabled bool
	remoteEnabled   bool
	autoStart      bool
)

func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()
	logView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	showWelcome()

	if err := app.SetRoot(pages, true).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func log(msg string) {
	fmt.Fprintf(logView, "%s\n", msg)
	app.Draw()
}

func logSuccess(msg string) {
	fmt.Fprintf(logView, "[green]✓ %s[white]\n", msg)
	app.Draw()
}

func logError(msg string) {
	fmt.Fprintf(logView, "[red]✗ %s[white]\n", msg)
	app.Draw()
}

func logInfo(msg string) {
	fmt.Fprintf(logView, "[yellow]ℹ %s[white]\n", msg)
	app.Draw()
}

// ============================================================================
// Welcome Screen
// ============================================================================

func showWelcome() {
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(`[yellow]╔══════════════════════════════════════════════╗
║      🦞 Stardew MCP Installer 🦞               ║
║         Lobster Edition                       ║
╚══════════════════════════════════════════════╝[white]`)

	desc := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(`This installer will set up everything you need:
 • Build Go MCP Server
 • Build C# Stardew Valley Mod
 • Install mod to your game folder
 • Configure OpenClaw & Remote options`)

	btnInstall := tview.NewButton("[Install Everything]").SetSelectedFunc(func() {
		showPathDetection()
	})

	btnExit := tview.NewButton("[Exit]").SetSelectedFunc(func() {
		app.Stop()
	})

	buttonBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(btnInstall, 0, 1, true).
		AddItem(btnExit, 0, 1, false)

	menu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 3, false).
		AddItem(header, 8, 0, false).
		AddItem(desc, 6, 0, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonBox, 3, 0, false).
		AddItem(tview.NewBox(), 0, 3, false)

	pages.AddPage("welcome", menu, true, true)
	pages.SwitchToPage("welcome")
}

// ============================================================================
// Path Detection Screen
// ============================================================================

func showPathDetection() {
	stardewPath = detectStardewValley()

	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(`[yellow]Stardew Valley Location[white]

Enter the path where Stardew Valley is installed`)

	pathLabel := tview.NewTextView().
		SetText(fmt.Sprintf("Auto-detected: [green]%s[white]", stardewPath)).
		SetTextAlign(tview.AlignCenter)

	inputField := tview.NewInputField().
		SetLabel("Path: ").
		SetText(stardewPath).
		SetFieldWidth(50)

	inputField.SetChangedFunc(func(text string) {
		stardewPath = text
		pathLabel.SetText(fmt.Sprintf("Path: [green]%s[white]", text))
		app.Draw()
	})

	btnDetect := tview.NewButton("[Auto-Detect]").SetSelectedFunc(func() {
		stardewPath = detectStardewValley()
		inputField.SetText(stardewPath)
		app.Draw()
	})

	btnNext := tview.NewButton("[Next >]").SetSelectedFunc(func() {
		if stardewPath == "" || !pathExists(stardewPath) {
			showErrorModal("Please enter a valid Stardew Valley path")
			return
		}
		showOptions()
	})

	btnBack := tview.NewButton("[< Back]").SetSelectedFunc(func() {
		showWelcome()
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(header, 5, 0, false).
		AddItem(pathLabel, 2, 0, false).
		AddItem(inputField, 3, 0, false).
		AddItem(btnDetect, 1, 0, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(btnBack, 1, 0, false).
		AddItem(btnNext, 1, 0, false).
		AddItem(tview.NewBox(), 0, 2, false)

	pages.AddPage("path", flex, true, true)
	pages.SwitchToPage("path")
}

// ============================================================================
// Options Screen
// ============================================================================

func showOptions() {
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(`[yellow]Additional Options[white]`)

	var openclawCheck, remoteCheck, autoCheck *tview.CheckBox

	openclawCheck = tview.NewCheckBox().SetLabel("Enable OpenClaw Gateway").SetChecked(false)
	remoteCheck = tview.NewCheckBox().SetLabel("Enable Remote Server Mode").SetChecked(false)
	autoCheck = tview.NewCheckBox().SetLabel("Auto-start agent on connect").SetChecked(true)

	openclawCheck.SetChangedFunc(func(checked bool) {
		openclawEnabled = checked
	})

	remoteCheck.SetChangedFunc(func(checked bool) {
		remoteEnabled = checked
	})

	autoCheck.SetChangedFunc(func(checked bool) {
		autoStart = checked
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(header, 3, 0, false).
		AddItem(openclawCheck, 1, 0, false).
		AddItem(remoteCheck, 1, 0, false).
		AddItem(autoCheck, 1, 0, false).
		AddItem(tview.NewBox(), 0, 1, false)

	btnInstall := tview.NewButton("[Install Now]").SetSelectedFunc(func() {
		showInstallProgress()
	})

	btnBack := tview.NewButton("[< Back]").SetSelectedFunc(func() {
		showPathDetection()
	})

	flex.AddItem(btnBack, 1, 0, false).
		AddItem(btnInstall, 1, 0, false).
		AddItem(tview.NewBox(), 0, 2, false)

	pages.AddPage("options", flex, true, true)
	pages.SwitchToPage("options")
}

// ============================================================================
// Installation Progress Screen
// ============================================================================

func showInstallProgress() {
	logView.Clear()

	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(`[yellow]Installing... Please Wait[white]`)

	logBox := tview.NewFrame(logView).
		SetBorders(tview.BorderDouble, " ", " ", " ", " ", " ", " ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(header, 3, 0, false).
		AddItem(logBox, 0, 8, false).
		AddItem(tview.NewBox(), 0, 1, false)

	pages.AddPage("install", flex, true, true)
	pages.SwitchToPage("install")

	go runInstallation()
}

func runInstallation() {
	time.Sleep(500 * time.Millisecond)

	// Step 1: Check Go
	logInfo("Checking Go installation...")
	if !commandExists("go") {
		logError("Go not found!")
		showErrorModal("Go is not installed. Please install Go 1.23+ from https://go.dev/dl/")
		return
	}
	logSuccess("Go found!")

	// Step 2: Check .NET
	logInfo("Checking .NET SDK...")
	if !commandExists("dotnet") {
		logError(".NET SDK not found!")
		showErrorModal(".NET SDK not found. Please install .NET 6.0 from https://dotnet.microsoft.com/download")
		return
	}
	logSuccess(".NET found!")

	// Step 3: Build Go server
	logInfo("Building Go MCP Server...")
	if err := buildGoServer(); err != nil {
		logError("Failed to build Go server")
		showErrorModal(fmt.Sprintf("Failed to build Go server: %v", err))
		return
	}
	logSuccess("Go MCP Server built!")

	// Step 4: Build C# Mod
	logInfo("Building C# Stardew Mod...")
	if err := buildCSharpMod(); err != nil {
		logError("Failed to build C# mod")
		showErrorModal(fmt.Sprintf("Failed to build C# mod: %v", err))
		return
	}
	logSuccess("C# Mod built!")

	// Step 5: Install Mod
	logInfo("Installing mod to Stardew Valley...")
	if err := installMod(); err != nil {
		logError("Failed to install mod")
		showErrorModal(fmt.Sprintf("Failed to install mod: %v", err))
		return
	}
	logSuccess("Mod installed!")

	// Step 6: Create config
	logInfo("Creating configuration...")
	if err := createConfig(); err != nil {
		logError("Failed to create config")
	} else {
		logSuccess("Configuration created!")
	}

	log("")
	logSuccess("🎉 Installation Complete! 🎉")
	log("")

	showSuccess()
}

func showSuccess() {
	app.QueueUpdate(func() {
		options := ""
		if openclawEnabled {
			options += "\n • OpenClaw Gateway Enabled"
		}
		if remoteEnabled {
			options += "\n • Remote Server Enabled"
		}
		if options == "" {
			options = "\n • Default Configuration"
		}

		header := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(`[green]╔══════════════════════════════════════════════╗
║           Installation Complete!           ║
╚══════════════════════════════════════════════╝[white]`)

		text := fmt.Sprintf(`[yellow]Location:[white] %s

[yellow]Next Steps:[white]
1. Start Stardew Valley through SMAPI
2. Load your save file
3. Run: cd setup && run.bat

[yellow]Enabled Options:[white]%s`, stardewPath, options)

		desc := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)

		btnExit := tview.NewButton("[Exit]").SetSelectedFunc(func() {
			app.Stop()
		})

		flex := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 2, false).
			AddItem(header, 6, 0, false).
			AddItem(desc, 0, 4, false).
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(btnExit, 1, 0, false).
			AddItem(tview.NewBox(), 0, 2, false)

		pages.AddPage("success", flex, true, true)
		pages.SwitchToPage("success")
	})
}

func showErrorModal(msg string) {
	app.QueueUpdate(func() {
		modal := tview.NewModal().
			SetText(msg).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.SwitchToPage("install")
			})

		pages.AddPage("error", modal, true, false)
		app.Draw()
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

func detectStardewValley() string {
	paths := []string{}

	switch runtime.GOOS {
	case "windows":
		paths = []string{
			`C:\Program Files\Stardew Valley`,
			`C:\Program Files (x86)\Stardew Valley`,
			filepath.Join(os.Getenv("LocalAppData"), "StardewValley"),
			`D:\Games\Stardew Valley`,
		}
	case "darwin":
		paths = []string{
			"/Applications/Stardew Valley.app/Contents/MacOS",
			filepath.Join(os.Getenv("HOME"), "Applications/Stardew Valley.app/Contents/MacOS"),
		}
	case "linux":
		paths = []string{
			filepath.Join(os.Getenv("HOME"), ".local/share/Steam/steamapps/common/Stardew Valley"),
			filepath.Join(os.Getenv("HOME"), ".steam/steamapps/common/Stardew Valley"),
			"/opt/stardew-valley",
		}
	}

	for _, p := range paths {
		if pathExists(p) {
			return p
		}
	}

	return ""
}

func pathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func buildGoServer() error {
	dir := getCurrentDir()
	cmd := exec.Command("go", "build", "-o", "stardew-mcp")
	cmd.Dir = filepath.Join(dir, "..", "mcp-server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildCSharpMod() error {
	dir := getCurrentDir()
	cmd := exec.Command("dotnet", "build", "-c", "Release")
	cmd.Dir = filepath.Join(dir, "..", "mod", "StardewMCP")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installMod() error {
	modsDir := filepath.Join(stardewPath, "Mods", "StardewMCP")
	if err := os.MkdirAll(modsDir, 0755); err != nil {
		return err
	}

	srcDir := filepath.Join(getCurrentDir(), "..", "mod", "StardewMCP", "bin", "Release", "net6.0")
	return copyDir(srcDir, modsDir)
}

func createConfig() error {
	config := fmt.Sprintf(`server:
  game_url: "ws://localhost:8765/game"
  auto_start: %v
  log_level: "info"

remote:
  host: "0.0.0.0"
  port: 8765

openclaw:
  gateway_url: "ws://127.0.0.1:18789"
  token: ""
  agent_name: "stardew-farmer"
`, autoStart)

	configPath := filepath.Join(getCurrentDir(), "..", "mcp-server", "config.yaml")
	return os.WriteFile(configPath, []byte(config), 0644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}
		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
