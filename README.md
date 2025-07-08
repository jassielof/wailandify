# Linux Browser Desktop Manager

> 🚀 **Fix Wayland gesture support and PWA issues for Linux browsers**

Tired of manually editing desktop entries after every browser update? This tool automatically manages desktop entries for Chromium-based browsers on Linux, enabling proper Wayland support and gesture navigation for both main browsers and PWAs.

## ✨ Features

### Current Features
- **🔄 Automatic Desktop Entry Sync** - Copies system desktop entries to user directory (survives updates)
- **🖱️ Wayland Gesture Support** - Enables swipe back/forward navigation for PWAs
- **📱 PWA Support** - Automatically finds and fixes all PWA desktop entries
- **🛡️ Update-Proof** - Your custom entries won't be overridden by browser updates
- **🎯 Smart Flag Injection** - Only adds flags if not already present (no duplicates)
- **🔍 Auto-Discovery** - Finds PWA files using browser-specific prefixes

### Supported Browsers
- **Brave Browser** (`brave-*.desktop` PWAs)
- **Microsoft Edge** (`msedge-*.desktop` PWAs)

### Applied Wayland Flags
- `--enable-features=TouchpadOverscrollHistoryNavigation` - Enables swipe gestures
- `--ozone-platform=wayland` - Proper Wayland platform support

## 🚀 Installation

### Option 1: Download Binary (Recommended)
```bash
# Download from releases (coming soon)
wget https://github.com/yourusername/linux-browser-desktop-manager/releases/latest/download/desktop-manager
chmod +x desktop-manager
sudo mv desktop-manager /usr/local/bin/
```

### Option 2: Build from Source
```bash
# Clone repository
git clone https://github.com/yourusername/linux-browser-desktop-manager.git
cd linux-browser-desktop-manager

# Build
go build -o desktop-manager main.go

# Install
sudo mv desktop-manager /usr/local/bin/
```

### Option 3: Go Install
```bash
go install github.com/yourusername/linux-browser-desktop-manager@latest
```

## 📖 Usage

### Basic Usage
```bash
# Run the desktop manager
desktop-manager
```

### Typical Workflow
1. Install your browsers (Brave, Edge)
2. Create some PWAs
3. Run `desktop-manager` to fix gesture support
4. After browser updates, run `desktop-manager` again to sync changes

### Example Output
```
🚀 Desktop Entry Manager for Linux Browsers
============================================

🔍 Processing Brave...
✅ Updated Brave main desktop file
🔗 Found 3 PWA files for Brave
✅ Updated PWA: brave-youtube.desktop
✅ Updated PWA: brave-gmail.desktop
✅ Updated PWA: brave-discord.desktop

🔍 Processing Edge...
✅ Updated Edge main desktop file
🔗 Found 2 PWA files for Edge
✅ Updated PWA: msedge-spotify.desktop
✅ Updated PWA: msedge-teams.desktop

🎉 Desktop entry management completed!
💡 Tip: Run this script after browser updates to keep entries synchronized
```

## 🧩 Why This Tool Exists

### The Problem
On Linux, Chromium-based browsers have several issues:

1. **Firefox PWAs** - Requires extensions, creates separate processes, no cookie sync
2. **Brave/Edge PWAs** - Good PWA support but missing swipe gestures
3. **Manual Flag Management** - Need to manually add Wayland flags to each desktop entry
4. **Update Overrides** - System desktop entries get reset after browser updates
5. **PWA Inheritance** - PWAs don't inherit browser flags unless browser is opened first

### The Solution
This tool automatically:
- Copies desktop entries to user directory (won't be overridden)
- Adds proper Wayland flags to enable gesture support
- Handles both main browser and all PWA entries
- Can be re-run after updates to stay synchronized

## 🔮 TODO / Future Features

### Browser Support
- [ ] **Google Chrome** support (`google-chrome-*.desktop`)
- [ ] **Firefox** PWA support (if they improve their implementation)
- [ ] **Chromium** support (`chromium-*.desktop`)
- [ ] **Opera** support (`opera-*.desktop`)
- [ ] **Vivaldi** support (`vivaldi-*.desktop`)

### Configuration & Flexibility
- [ ] **Config file support** - Custom flags per browser via YAML/JSON
- [ ] **Command-line flags** - Override default Wayland flags
- [ ] **Selective processing** - Choose which browsers to process
- [ ] **Dry-run mode** - Preview changes without applying them
- [ ] **Backup/restore** - Backup original entries before modification

### Advanced Features
- [ ] **Auto-detection** - Discover installed browsers automatically
- [ ] **Flag validation** - Check if flags are actually supported
- [ ] **Desktop environment detection** - Different flags for GNOME/KDE/etc
- [ ] **Systemd integration** - Run automatically after package updates
- [ ] **GUI version** - Simple graphical interface for non-technical users

### Developer Experience
- [ ] **Plugin system** - Add custom browser support via plugins
- [ ] **Testing suite** - Automated tests for desktop entry parsing
- [ ] **Cross-platform** - Support for other Linux distributions
- [ ] **Logging levels** - Verbose/quiet output modes

### Integration & Automation
- [ ] **Package manager hooks** - Auto-run after browser updates
- [ ] **Desktop notifications** - Notify when entries are updated
- [ ] **System tray integration** - Background monitoring and updates
- [ ] **Web interface** - Manage desktop entries via web UI

## 🛠️ Technical Details

### File Locations
- **System desktop entries**: `/usr/share/applications/`
- **User desktop entries**: `~/.local/share/applications/`
- **PWA detection patterns**: `brave-*.desktop`, `msedge-*.desktop`

### How It Works
1. Scans system applications directory for browser desktop files
2. Copies found files to user applications directory
3. Modifies `Exec=` lines to include Wayland flags
4. Processes both main browser and PWA files
5. Ensures flags aren't duplicated on subsequent runs

## 🤝 Contributing

Contributions are welcome! Please feel free to:
- Report bugs and issues
- Suggest new features
- Submit pull requests
- Improve documentation

### Development Setup
```bash
git clone https://github.com/yourusername/linux-browser-desktop-manager.git
cd linux-browser-desktop-manager
go run main.go
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the frustrations of Linux desktop browser management
- Thanks to the Wayland and browser development communities
- Built for the Linux community, by the Linux community

---

**⭐ Star this repo if it helped you fix your browser gesture issues!**
