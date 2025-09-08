# 🖥️ Shef GUI Desktop Setup

Easy setup for desktop integration with the **Evil Gopher** logo!

## 🚀 Quick Setup

```bash
# Clone and setup
git clone <your-repo-url>
cd shef
./install-desktop.sh
```

## 📋 What This Does

The `install-desktop.sh` script will:

✅ **Build the application** (`shef-gui` binary)  
✅ **Create desktop launcher** with Evil Gopher icon  
✅ **Add to applications menu** (searchable as "Shef GUI" or "Shodan")  
✅ **Copy to desktop** for double-click launching  

## 🎨 Logo & Icon

- **Logo**: `evilgopher.png` - The Evil Gopher mascot
- **Desktop Icon**: Automatically configured in launcher
- **App Window**: Uses the same Evil Gopher branding

## 🖱️ Usage After Setup

**Desktop**: Double-click the "Shef GUI" icon on your desktop  
**App Menu**: Search for "Shef GUI" or "Shodan" in your application launcher  
**Terminal**: Run `./shef-gui` from the project directory  

## 📁 Files Created

- `~/Desktop/shef-gui.desktop` - Desktop launcher
- `~/.local/share/applications/shef-gui.desktop` - App menu entry
- `shef-gui` - Main application binary

## 🔧 Manual Setup (Alternative)

If you prefer manual setup:

```bash
# Build the app
go build -o shef-gui .

# Make desktop file executable
chmod +x shef-gui.desktop

# Copy to desktop
cp shef-gui.desktop ~/Desktop/

# Copy to applications
cp shef-gui.desktop ~/.local/share/applications/
```

## 🎯 Features

- **GTK4 Native GUI** for Shodan searches
- **Evil Gopher Branding** throughout the interface
- **Cross-platform** (Linux, Windows, Mac)
- **Minimal & Fast** - No bloat, just functionality
