package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"

	// Add math import
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

const (
	version = "1.0.0"
	banner  = `
╔══════════════════════════════════════════════════════════════╗
║                    🚀 LAPTOP ACTIVITY TOOL                   ║
║                        Keep Your System Awake                ║
║                           Version %s                        ║
╚══════════════════════════════════════════════════════════════╝
`
)

// Windows API constants
const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	MOUSEEVENTF_MOVE = 0x0001
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP = 0x0004
	MOUSEEVENTF_WHEEL = 0x0800
	KEYEVENTF_KEYUP = 0x0002
	VK_F13 = 0x7C
	VK_F14 = 0x7D
	VK_F15 = 0x7E
	VK_SCROLL = 0x91
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procSendInput = user32.NewProc("SendInput")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procSetThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
)

// Windows structures
type POINT struct {
	X, Y int32
}

type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
	Ki   KEYBDINPUT
}

type ActivityConfig struct {
	Interval       time.Duration
	MouseEnabled   bool
	KeyboardEnabled bool
	MemoryEnabled  bool
	Verbose        bool
	Duration       time.Duration
	IntensityLevel int
}

type ActivityStats struct {
	StartTime      time.Time
	MouseMoves     int
	KeyPresses     int
	MemoryOps      int
	TotalActions   int
	LastActivity   string
}

type ActivityTool struct {
	config *ActivityConfig
	stats  *ActivityStats
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	// Colors
	green  = color.New(color.FgGreen, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
	blue   = color.New(color.FgBlue, color.Bold)
	cyan   = color.New(color.FgCyan, color.Bold)

	// Random words for memory operations
	randomWords = []string{
		"laptop", "computer", "system", "active", "awake", "running", "process",
		"background", "task", "automation", "script", "program", "application",
		"memory", "cpu", "performance", "monitor", "status", "update", "refresh",
		"simulate", "generate", "execute", "maintain", "preserve", "prevent",
		"sleep", "idle", "standby", "hibernate", "power", "energy", "efficient",
	}
)

func main() {
	// Parse command line flags
	var (
		interval    = flag.Duration("interval", 3*time.Second, "Interval between activities")
		duration    = flag.Duration("duration", 0, "How long to run (0 = infinite)")
		mouse       = flag.Bool("mouse", true, "Enable mouse movements")
		keyboard    = flag.Bool("keyboard", true, "Enable keyboard simulation")
		memory      = flag.Bool("memory", true, "Enable memory operations")
		verbose     = flag.Bool("verbose", false, "Enable verbose logging")
		intensity   = flag.Int("intensity", 2, "Activity intensity level (1-5)")
		interactive = flag.Bool("interactive", false, "Run in interactive mode")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Laptop Activity Tool v%s\n", version)
		fmt.Printf("Built with Go %s for %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	// Print banner
	fmt.Printf(banner, version)

	config := &ActivityConfig{
		Interval:        *interval,
		MouseEnabled:    *mouse,
		KeyboardEnabled: *keyboard,
		MemoryEnabled:   *memory,
		Verbose:         *verbose,
		Duration:        *duration,
		IntensityLevel:  *intensity,
	}

	if *interactive {
		runInteractiveMode(config)
	} else {
		runActivity(config)
	}
}

func runInteractiveMode(config *ActivityConfig) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printMenu()
		fmt.Print("Kies een optie: ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			configureSettings(config, scanner)
		case "2":
			runActivity(config)
		case "3":
			showStatus(config)
		case "4":
			showHelp()
		case "5", "q", "quit", "exit":
			cyan.Println("\n👋 Tot ziens!")
			return
		default:
			red.Println("❌ Ongeldige keuze. Probeer opnieuw.")
		}
	}
}

func printMenu() {
	fmt.Println()
	blue.Println("=== HOOFDMENU ===")
	fmt.Println("1. 🔧 Instellingen aanpassen")
	fmt.Println("2. 🚀 Start activiteit")
	fmt.Println("3. 📊 Status bekijken")
	fmt.Println("4. ❓ Help")
	fmt.Println("5. 🚪 Afsluiten")
	fmt.Println()
}

func configureSettings(config *ActivityConfig, scanner *bufio.Scanner) {
	fmt.Println()
	yellow.Println("=== INSTELLINGEN ===")

	// Interval
	fmt.Printf("Huidige interval: %v\n", config.Interval)
	fmt.Print("Nieuw interval in seconden (enter = behouden): ")
	if scanner.Scan() {
		if input := strings.TrimSpace(scanner.Text()); input != "" {
			if seconds, err := strconv.Atoi(input); err == nil && seconds > 0 {
				config.Interval = time.Duration(seconds) * time.Second
				green.Printf("✅ Interval ingesteld op %v\n", config.Interval)
			} else {
				red.Println("❌ Ongeldige waarde")
			}
		}
	}

	// Intensity
	fmt.Printf("Huidige intensiteit: %d (1-5)\n", config.IntensityLevel)
	fmt.Print("Nieuwe intensiteit (enter = behouden): ")
	if scanner.Scan() {
		if input := strings.TrimSpace(scanner.Text()); input != "" {
			if intensity, err := strconv.Atoi(input); err == nil && intensity >= 1 && intensity <= 5 {
				config.IntensityLevel = intensity
				green.Printf("✅ Intensiteit ingesteld op %d\n", config.IntensityLevel)
			} else {
				red.Println("❌ Intensiteit moet tussen 1 en 5 zijn")
			}
		}
	}

	// Duration
	fmt.Printf("Huidige duur: %v (0 = oneindig)\n", config.Duration)
	fmt.Print("Nieuwe duur in minuten (0 = oneindig, enter = behouden): ")
	if scanner.Scan() {
		if input := strings.TrimSpace(scanner.Text()); input != "" {
			if minutes, err := strconv.Atoi(input); err == nil && minutes >= 0 {
				if minutes == 0 {
					config.Duration = 0
				} else {
					config.Duration = time.Duration(minutes) * time.Minute
				}
				green.Printf("✅ Duur ingesteld op %v\n", config.Duration)
			} else {
				red.Println("❌ Ongeldige waarde")
			}
		}
	}

	// Toggles
	fmt.Printf("Muis bewegingen: %v | Toetsenbord: %v | Geheugen: %v | Verbose: %v\n",
		config.MouseEnabled, config.KeyboardEnabled, config.MemoryEnabled, config.Verbose)
	fmt.Print("Wijzig instellingen? (m=muis, k=toetsenbord, g=geheugen, v=verbose, enter=klaar): ")
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		for _, char := range input {
			switch char {
			case 'm':
				config.MouseEnabled = !config.MouseEnabled
				fmt.Printf("Muis: %v\n", config.MouseEnabled)
			case 'k':
				config.KeyboardEnabled = !config.KeyboardEnabled
				fmt.Printf("Toetsenbord: %v\n", config.KeyboardEnabled)
			case 'g':
				config.MemoryEnabled = !config.MemoryEnabled
				fmt.Printf("Geheugen: %v\n", config.MemoryEnabled)
			case 'v':
				config.Verbose = !config.Verbose
				fmt.Printf("Verbose: %v\n", config.Verbose)
			}
		}
	}

	green.Println("✅ Instellingen bijgewerkt!")
}

func showStatus(config *ActivityConfig) {
	fmt.Println()
	cyan.Println("=== HUIDIGE STATUS ===")
	fmt.Printf("🕐 Interval: %v\n", config.Interval)
	fmt.Printf("⚡ Intensiteit: %d/5\n", config.IntensityLevel)
	fmt.Printf("⏱️  Duur: %v\n", config.Duration)
	fmt.Printf("🖱️  Muis: %v\n", formatBool(config.MouseEnabled))
	fmt.Printf("⌨️  Toetsenbord: %v\n", formatBool(config.KeyboardEnabled))
	fmt.Printf("🧠 Geheugen: %v\n", formatBool(config.MemoryEnabled))
	fmt.Printf("📝 Verbose: %v\n", formatBool(config.Verbose))
	fmt.Printf("💻 Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func showHelp() {
	fmt.Println()
	cyan.Println("=== HELP ===")
	fmt.Println("Deze tool houdt je laptop wakker door verschillende activiteiten uit te voeren:")
	fmt.Println()
	fmt.Println("🖱️  Muis bewegingen - Beweegt de cursor subtiel")
	fmt.Println("⌨️  Toetsenbord - Simuleert veilige toetsaanslagen (F13-F15)")
	fmt.Println("🧠 Geheugen operaties - Voert achtergrond berekeningen uit")
	fmt.Println("💤 Systeem awake - Voorkomt automatische slaapstand")
	fmt.Println()
	fmt.Println("Intensiteit niveaus:")
	fmt.Println("1 - Zeer laag (minimale activiteit)")
	fmt.Println("2 - Laag (standaard)")
	fmt.Println("3 - Gemiddeld")
	fmt.Println("4 - Hoog")
	fmt.Println("5 - Zeer hoog (maximale activiteit)")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println("- Gebruik Ctrl+C om de activiteit te stoppen")
	fmt.Println("- Lagere intervallen = meer activiteit")
	fmt.Println("- Verbose mode toont alle acties")
	fmt.Println("- Duur 0 = oneindig uitvoeren")
	fmt.Println("- F13-F15 toetsen zijn veilig en doen niets")
}

func runActivity(config *ActivityConfig) {
	fmt.Println()
	green.Println("🚀 Laptop Activity Tool wordt gestart...")

	// Check if we're on Windows
	if runtime.GOOS != "windows" {
		yellow.Println("⚠️  Deze versie is geoptimaliseerd voor Windows")
		yellow.Println("💡 Alleen geheugen operaties zullen werken op andere platforms")
	}

	tool := &ActivityTool{
		config: config,
		stats: &ActivityStats{
			StartTime: time.Now(),
		},
	}

	tool.ctx, tool.cancel = context.WithCancel(context.Background())

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Setup duration timeout if specified
	if config.Duration > 0 {
		tool.ctx, tool.cancel = context.WithTimeout(tool.ctx, config.Duration)
	}

	// Prevent system sleep
	tool.preventSystemSleep()

	// Start activity in goroutine
	go tool.runActivityLoop()

	// Show initial status
	tool.printStatus()

	// Wait for completion or signal
	select {
	case <-sigChan:
		yellow.Println("\n⏹️  Stop signaal ontvangen...")
	case <-tool.ctx.Done():
		if config.Duration > 0 {
			green.Println("\n⏰ Geplande duur bereikt...")
		}
	}

	tool.cancel()
	tool.allowSystemSleep()
	tool.printFinalStats()
	green.Println("👋 Activiteit gestopt. Tot ziens!")
}

func (a *ActivityTool) runActivityLoop() {
	ticker := time.NewTicker(a.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.performActivity()
		}
	}
}

func (a *ActivityTool) performActivity() {
	activities := a.getEnabledActivities()
	if len(activities) == 0 {
		return
	}

	// Perform multiple activities based on intensity
	actionsCount := a.config.IntensityLevel
	for i := 0; i < actionsCount; i++ {
		activity := activities[rand.Intn(len(activities))]
		activity()

		// Small delay between actions
		time.Sleep(time.Millisecond * 100)
	}

	a.stats.TotalActions++
}

func (a *ActivityTool) getEnabledActivities() []func() {
	var activities []func()

	if a.config.MouseEnabled && runtime.GOOS == "windows" {
		activities = append(activities,
			a.simulateMouseMove,
			a.simulateMouseClick,
			a.simulateMouseScroll,
		)
	}

	if a.config.KeyboardEnabled && runtime.GOOS == "windows" {
		activities = append(activities,
			a.simulateKeyPress,
		)
	}

	if a.config.MemoryEnabled {
		activities = append(activities,
			a.performMemoryOperation,
			a.performCPUOperation,
		)
	}

	return activities
}

// Windows API functions
func (a *ActivityTool) simulateMouseMove() {
	if runtime.GOOS != "windows" {
		return
	}

	// Get screen dimensions
	width, _, _ := procGetSystemMetrics.Call(0)  // SM_CXSCREEN
	height, _, _ := procGetSystemMetrics.Call(1) // SM_CYSCREEN

	// Get current position
	var pt POINT
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))

	// Move more noticeably (100-300 pixels in any direction)
	moveDistance := 100 + rand.Intn(200) // 100-300 pixels
	angle := rand.Float64() * 2 * math.Pi // Random direction

	deltaX := int32(float64(moveDistance) * math.Cos(angle))
	deltaY := int32(float64(moveDistance) * math.Sin(angle))

	newX := int32(pt.X) + deltaX
	newY := int32(pt.Y) + deltaY

	// Keep within screen bounds
	if newX < 0 { newX = 0 }
	if newY < 0 { newY = 0 }
	if newX >= int32(width) { newX = int32(width) - 1 }
	if newY >= int32(height) { newY = int32(height) - 1 }

	procSetCursorPos.Call(uintptr(newX), uintptr(newY))

	a.stats.MouseMoves++
	a.stats.LastActivity = fmt.Sprintf("Muis verplaatst naar (%d, %d)", newX, newY)

	if a.config.Verbose {
		fmt.Printf("🖱️  %s\n", a.stats.LastActivity)
	}
}

func (a *ActivityTool) simulateMouseClick() {
	if runtime.GOOS != "windows" {
		return
	}

	// Simulate a very brief left click
	input := []INPUT{
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DwFlags: MOUSEEVENTF_LEFTDOWN,
			},
		},
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DwFlags: MOUSEEVENTF_LEFTUP,
			},
		},
	}

	procSendInput.Call(
		uintptr(len(input)),
		uintptr(unsafe.Pointer(&input[0])),
		uintptr(unsafe.Sizeof(input[0])),
	)

	a.stats.MouseMoves++
	a.stats.LastActivity = "Muis klik uitgevoerd"

	if a.config.Verbose {
		fmt.Println("🖱️  Muis klik uitgevoerd")
	}
}

func (a *ActivityTool) simulateMouseScroll() {
	if runtime.GOOS != "windows" {
		return
	}

	direction := rand.Intn(2) // 0 = up, 1 = down
	var wheelData uint32
	if direction == 0 {
		wheelData = 120  // Positive for up
	} else {
		wheelData = 0xFFFFFF88  // Two's complement of -120 for down
	}

	input := INPUT{
		Type: INPUT_MOUSE,
		Mi: MOUSEINPUT{
			DwFlags:   MOUSEEVENTF_WHEEL,
			MouseData: wheelData,
		},
	}

	procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&input)),
		uintptr(unsafe.Sizeof(input)),
	)

	a.stats.MouseMoves++
	dirStr := "up"
	if direction == 1 {
		dirStr = "down"
	}
	a.stats.LastActivity = "Muis scroll " + dirStr

	if a.config.Verbose {
		fmt.Printf("🖱️  Muis scroll %s\n", dirStr)
	}
}

func (a *ActivityTool) simulateKeyPress() {
	if runtime.GOOS != "windows" {
		return
	}

	// Use safe function keys that don't interfere with system
	safeKeys := []uint16{VK_F13, VK_F14, VK_F15, VK_SCROLL}
	keyNames := []string{"F13", "F14", "F15", "ScrollLock"}

	keyIndex := rand.Intn(len(safeKeys))
	vkCode := safeKeys[keyIndex]
	keyName := keyNames[keyIndex]

	// Key down then key up
	input := []INPUT{
		{
			Type: INPUT_KEYBOARD,
			Ki: KEYBDINPUT{
				WVk: vkCode,
			},
		},
		{
			Type: INPUT_KEYBOARD,
			Ki: KEYBDINPUT{
				WVk:     vkCode,
				DwFlags: KEYEVENTF_KEYUP,
			},
		},
	}

	procSendInput.Call(
		uintptr(len(input)),
		uintptr(unsafe.Pointer(&input[0])),
		uintptr(unsafe.Sizeof(input[0])),
	)

	a.stats.KeyPresses++
	a.stats.LastActivity = "Toets: " + keyName

	if a.config.Verbose {
		fmt.Printf("⌨️  Toets ingedrukt: %s\n", keyName)
	}
}

func (a *ActivityTool) performMemoryOperation() {
	// Create and manipulate some data in memory
	size := rand.Intn(1000) + 100
	data := make([]float64, size)

	for i := range data {
		data[i] = rand.Float64() * 1000
	}

	// Perform some operations
	sum := 0.0
	for _, v := range data {
		sum += math.Sin(v) * math.Cos(v)
	}

	a.stats.MemoryOps++
	a.stats.LastActivity = fmt.Sprintf("Geheugen operatie (%d items)", len(data))

	if a.config.Verbose {
		fmt.Printf("🧠 Geheugen operatie uitgevoerd (%d items, sum: %.2f)\n", len(data), sum)
	}
}

func (a *ActivityTool) performCPUOperation() {
	// Perform some CPU-intensive but brief calculations
	iterations := a.config.IntensityLevel * 10000
	result := 0.0

	for i := 0; i < iterations; i++ {
		result += math.Sqrt(float64(i)) * math.Log(float64(i+1))
	}

	a.stats.MemoryOps++
	a.stats.LastActivity = fmt.Sprintf("CPU operatie (%d iteraties)", iterations)

	if a.config.Verbose {
		fmt.Printf("⚡ CPU operatie uitgevoerd (%d iteraties, result: %.2f)\n", iterations, result)
	}
}

func (a *ActivityTool) preventSystemSleep() {
	if runtime.GOOS == "windows" {
		// ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED
		procSetThreadExecutionState.Call(0x80000000 | 0x00000001 | 0x00000002)
		if a.config.Verbose {
			fmt.Println("💤 Systeem slaapstand uitgeschakeld")
		}
	}
}

func (a *ActivityTool) allowSystemSleep() {
	if runtime.GOOS == "windows" {
		// ES_CONTINUOUS
		procSetThreadExecutionState.Call(0x80000000)
		if a.config.Verbose {
			fmt.Println("💤 Systeem slaapstand weer ingeschakeld")
		}
	}
}

func (a *ActivityTool) printStatus() {
	fmt.Println()
	green.Println("✅ Activiteit gestart!")
	fmt.Printf("⏱️  Interval: %v\n", a.config.Interval)
	fmt.Printf("⚡ Intensiteit: %d/5\n", a.config.IntensityLevel)
	if a.config.Duration > 0 {
		fmt.Printf("⏰ Duur: %v\n", a.config.Duration)
	} else {
		fmt.Println("⏰ Duur: Oneindig")
	}
	fmt.Printf("🕐 Gestart: %s\n", a.stats.StartTime.Format("15:04:05"))
	fmt.Println()
	yellow.Println("Druk Ctrl+C om te stoppen...")

	// Show progress if duration is set
	if a.config.Duration > 0 {
		go a.showProgress()
	} else {
		go a.showLiveStats()
	}
}

func (a *ActivityTool) showProgress() {
	if a.config.Duration <= 0 {
		return
	}

	bar := progressbar.NewOptions(int(a.config.Duration.Seconds()),
		progressbar.OptionSetDescription("⏳ Voortgang"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "▐",
			BarEnd:        "▌",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("sec"),
	)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			bar.Finish()
			return
		case <-ticker.C:
			elapsed := time.Since(a.stats.StartTime)
			bar.Set(int(elapsed.Seconds()))
		}
	}
}

func (a *ActivityTool) showLiveStats() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			elapsed := time.Since(a.stats.StartTime)
			fmt.Printf("\r📊 Actief: %v | Acties: %d | Laatste: %s",
				elapsed.Truncate(time.Second),
				a.stats.TotalActions,
				a.stats.LastActivity)
		}
	}
}

func (a *ActivityTool) printFinalStats() {
	elapsed := time.Since(a.stats.StartTime)

	fmt.Println()
	cyan.Println("=== EINDSTATISTIEKEN ===")
	fmt.Printf("⏱️  Totale tijd: %v\n", elapsed.Truncate(time.Second))
	fmt.Printf("🎯 Totale acties: %d\n", a.stats.TotalActions)
	fmt.Printf("🖱️  Muis bewegingen: %d\n", a.stats.MouseMoves)
	fmt.Printf("⌨️  Toetsaanslagen: %d\n", a.stats.KeyPresses)
	fmt.Printf("🧠 Geheugen operaties: %d\n", a.stats.MemoryOps)

	if elapsed.Seconds() > 0 {
		rate := float64(a.stats.TotalActions) / elapsed.Seconds()
		fmt.Printf("📈 Gemiddelde: %.2f acties/sec\n", rate)
	}
}

func formatBool(b bool) string {
	if b {
		return green.Sprint("AAN")
	}
	return red.Sprint("UIT")
}

// go.mod inhoud:
/*
module laptop-activity-tool

go 1.21

require (
    github.com/fatih/color v1.15.0
    github.com/schollz/progressbar/v3 v3.13.1
)

require (
    github.com/mattn/go-colorable v0.1.13 // indirect
    github.com/mattn/go-isatty v0.0.17 // indirect
    github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
    github.com/rivo/uniseg v0.4.4 // indirect
    golang.org/x/sys v0.6.0 // indirect
    golang.org/x/term v0.6.0 // indirect
)
*/