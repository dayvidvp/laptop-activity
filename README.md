# 🚀 Laptop Activity Tool

## Keep Your System Awake, Credibly.

This simple yet effective Go-based utility is designed to prevent your laptop from entering sleep mode or locking its screen during short periods of inactivity. Whether you're stepping away for a quick coffee break, walking your dog, or simply need to ensure your system remains active without constant interaction, the Laptop Activity Tool has you covered.

Unlike basic mouse jigglers, this tool simulates a range of credible user activities, making it appear as though you're still actively using your computer. This helps maintain your presence in online meetings, prevents interruptions during long downloads, or keeps your development environment ready.

## ✨ Features

*   **Mouse Movements:** Subtle, random cursor movements across the screen.
*   **Mouse Clicks & Scrolls:** Occasional, non-disruptive left clicks and scroll actions.
*   **Keyboard Presses:** Simulation of safe, non-interfering key presses (F13-F15, Scroll Lock).
*   **Background Operations:** Performs light CPU and memory operations to simulate background tasks.
*   **System Awake:** Explicitly prevents the system from going to sleep or turning off the display (Windows only).
*   **Configurable Intensity:** Adjust the frequency and type of activities to suit your needs (1-5 levels).
*   **Interactive Mode:** Easy-to-use command-line interface for on-the-fly configuration.
*   **Duration Control:** Set a specific runtime or let it run indefinitely.

## 🛠️ Installation & Usage

### Prerequisites

*   Go (version 1.21 or higher)

### Build from Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/YOUR_USERNAME/laptop-activity-tool.git
    cd laptop-activity-tool
    ```
2.  Build the executable:
    ```bash
    go build -o laptop-activity-tool.exe .
    ```

### Running the Tool

#### Interactive Mode (Recommended for first-time users)

```bash
./laptop-activity-tool.exe --interactive
```

This will launch an interactive menu where you can easily configure settings like interval, intensity, and enable/disable specific activities.

#### Command-Line Flags

You can also run the tool directly with command-line flags:

```bash
./laptop-activity-tool.exe --interval 5s --intensity 3 --mouse --keyboard --memory --duration 30m --verbose
```

**Available Flags:**

*   `--interval <duration>`: Interval between activity bursts (e.g., `3s`, `10s`). Default: `3s`.
*   `--duration <duration>`: How long to run the tool (e.g., `1h`, `30m`). Use `0` for infinite. Default: `0`.
*   `--mouse`: Enable mouse movements, clicks, and scrolls. Default: `true`.
*   `--keyboard`: Enable safe keyboard presses. Default: `true`.
*   `--memory`: Enable background memory and CPU operations. Default: `true`.
*   `--intensity <level>`: Activity intensity level (1-5). 1 = very low, 5 = very high. Default: `2`.
*   `--verbose`: Enable verbose logging of activities. Default: `false`.
*   `--interactive`: Run in interactive mode. Default: `false`.
*   `--version`: Show version information.

## ⚠️ Compatibility

This tool is primarily optimized for **Windows** operating systems. Mouse and keyboard simulations, as well as explicit system sleep prevention, are currently only functional on Windows. On other platforms (Linux, macOS), only the background memory and CPU operations will be active.

## 🛑 Stopping the Tool

Press `Ctrl+C` in the terminal to gracefully stop the tool at any time.

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests to improve the tool.

## 📄 License

This project is licensed under the MIT License - see the `LICENSE` file for details. (Note: A `LICENSE` file is not yet present, but will be added.)

## 🌟 Future Enhancements

*   Cross-platform support for mouse/keyboard simulation (e.g., using platform-specific libraries).
*   More advanced, human-like mouse path generation.
*   Simulated browser activity (opening tabs, navigating, scrolling).
*   Simulated application launching/closing.
