# 📘 User Guide - Eggsmaker 2

## 📋 Index

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [First Start](#first-start)
4. [Main Interface](#main-interface)
5. [Workflow](#workflow)
6. [Advanced Options](#advanced-options)
7. [Troubleshooting](#troubleshooting)
8. [FAQ](#faq)

---

## 🎯 Introduction

**Eggsmaker 2** is a graphical application for **Penguins' Eggs** to create remasterable ISOs and manage Live installations and system packages.

### What does Eggsmaker 2 do?

- ✅ Creates ISOs of your current Linux system
- ✅ Clones your desktop and configurations
- ✅ Customizes ISO configuration
- ✅ Copies generated ISOs to custom locations
- ✅ Runs the Eggsmaker Plus ISO creation sequence
- ✅ Starts Live installations with Calamares or Krill
- ✅ Collects package lists and installs or removes native, AUR, Chaotic-AUR, and Flatpak packages
- ✅ Edits Eggs configuration in a graphical form
- ✅ Allows build temporaries to be removed at the end

### Prerequisites

- **Penguins' Eggs** installed on the system
- **sudo** permissions (password requested at startup)
- Python 3.8 or higher
- Internet connection (to install packages or download components)

### Live Installation and Packages

The **Live Installation** section starts Calamares or Krill. Calamares opens as a graphical application; Krill requires a supported graphical terminal and prompts there for the `sudo` password for its text interface.

In **Package Manager**, select **Collect** to create lists in `~/$(hostname)-pak/`, choose a source and packages, then use **Install selection** or **Uninstall selection**. Available sources depend on the distribution and installed tools. Arch supports native packages, AUR through `yay`, Chaotic-AUR, and Flatpak; Debian/Ubuntu supports manually installed APT packages and Flatpak.

**Configure Eggs** opens a graphical editor for `eggs.yaml`; it does not run the interactive text-mode `eggs config` wizard.

---

## 📦 Installation

### Method 1: Automatic Installer (Recommended)

#### Local Installation (User)
```bash
cd /path/to/eggsmaker2
./install.sh
```

This will install Eggsmaker 2 in `~/.local/share/eggsmaker2`

#### Global Installation (System)
```bash
cd /path/to/eggsmaker2
sudo ./install.sh --system
```

This will install Eggsmaker 2 in `/opt/eggsmaker2`

### Method 2: Installation with pip

```bash
cd /path/to/eggsmaker2
pip install -e .
```

### Method 3: Direct Execution (No Installation)

```bash
cd /path/to/eggsmaker2
python3 main.py
```

---

## 🚀 First Start

### 1. Run the Application

After installation:
```bash
eggsmaker2
```

Or from the application menu: search for **"Eggsmaker 2"**

### 2. Initial Authentication

Upon startup, a dialog will appear requesting the **sudo password**:

```
┌─────────────────────────────┐
│ Authentication required     │
│ (sudo)                      │
│                             │
│ Sudo password: [****]       │
│                             │
│ [OK]       [Cancel]         │
└─────────────────────────────┘
```

> 💡 **Important**: This password is required to execute privileged Penguins' Eggs commands.

### 3. Web Browser

The application will automatically open in your default browser at:
```
http://localhost:8080
```

To access it from another device, start the application with `EGGSMAKER_HOST=0.0.0.0 eggsmaker2` and use the address shown in the footer. Restrict access to trusted networks.

---

## 🖥️ Main Interface

### Top Bar

```
┌──────────────────────────────────────────────────────────────┐
│ 🥚 Eggsmaker - Version 26.10.21      Eggsmaker              │
│                                                              │
│                              Penguins' Eggs: X.X.X           │
└──────────────────────────────────────────────────────────────┘
```

### Output Terminal

```
┌──────────────────────────────────────────────────────────────┐
│ Output terminal (command log and operations)               │
│                                                              │
│ > Starting...                                               │
│ > ...                                                        │
└──────────────────────────────────────────────────────────────┘
```

### Main Actions

```
Create ISO  |  Live Installation  |  Configure Eggs  |  Package Manager
```

### Status Bar

```
┌──────────────────────────────────────────────────────────────┐
│ Status: Waiting                                              │
│ ████████████████████░░░░░░░░░ 65%                          │
│                                                              │
│ 65% | Copies: 2 | ISO: 3.5 GB | Copy: 00:02:30 | ...       │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔄 Workflow

1. At the desktop configuration prompt, choose whether to run `sudo eggs tools skel`.
2. If you choose to edit Eggs configuration, the graphical `eggs.yaml` form opens; the text-mode `eggs config` wizard is not launched.
3. Choose whether to create a standard ISO or a system clone (`--clone`).
4. Confirm generation. Eggsmaker runs `sudo eggs kill && eggs tools clean && eggs remaster`, adding `--clone` when selected.
5. Select a destination in the dialog. The application folder picker is reused here, and `rsync --info=progress2` reports progress in the terminal and progress bar.
6. Copy to additional destinations or choose to skip copying.
7. At the end, you can run `eggs kill` to remove build temporaries.

The sequence is linear; each decision appears graphically before the next action runs.

## ⚡ Advanced Options

### 🔧 Edit ISO Configuration

Available through **Configure Eggs**, which opens a graphical form with common fields and the complete YAML.

#### Editable Fields:

| Field | Description | Example |
|-------|-------------|---------|
| **Root password** | Superuser password in the ISO | `toor123` |
| **Snapshot basename** | Distribution name | `my-debian` |
| **Snapshot prefix** | Prefix for versions | `custom-` |
| **User password** | Live user password | `live123` |

> 📝 **Edited file**: `/etc/penguins-eggs.d/eggs.yaml`

### 📊 Real-time Indicators

#### Timers

| Timer | Color | Description |
|-------|-------|-------------|
| **Copy** | Cyan | Current copy time |
| **Generation** | Red | ISO generation time |
| **Total** | Lime Green | Sum of all times |

#### Progress Bar

- **Light Blue**: Active operation
- **Wavy**: ISO generation (indeterminate duration)
- **Red**: ISO copy (shows exact percentage)

#### Additional Indicators

- **ISO Size**: Updates after generation
- **Copies made**: Counter of successful copies
- **Percentage**: During copy (0-100%)

---

## 🔍 Troubleshooting

### ❌ Error: "Sudo password required"

**Cause**: Password was not provided at startup or session expired.

**Solution**: Reload the web page and provide the password again.

---

### ❌ Error: "ISO to copy not found"

**Cause**: ISO was not generated correctly or is not in the expected location.

**Solutions**:
1. Verify that generation completed successfully
2. Check the log terminal for errors during generation
3. Manually verify if the file exists:
   ```bash
   ls -lh /home/eggs/*.iso
   ls -lh /home/eggs/.mnt/*.iso
   ```

---

### ❌ Error during Eggs/Calamares update

**Cause**: Network or permission issues.

**Solutions**:
1. Check your internet connection
2. Try manually:
   ```bash
   cd ~
   git clone https://github.com/pieroproietti/fresh-eggs
   cd fresh-eggs
   sudo ./fresh-eggs.sh
   ```

---

### ⚠️ Application does not open in browser

**Solutions**:
1. Manually open: `http://localhost:8080`
2. Verify port is available:
   ```bash
   netstat -tuln | grep 8080
   ```
3. Use a different port:
   ```bash
   PORT=8081 eggsmaker2
   ```

---

### ⚠️ "Permission denied" when copying ISO

**Cause**: ISO is in a protected directory.

**Solution**: The application will automatically try to use `sudo` to copy.

---

### 🐛 Terminal shows no output

**Solution**: Reload the web page (F5).

---

## ❓ FAQ

### How much space do I need to create an ISO?

Depends on the mode:
- **Without data**: ~2-5 GB
- **With data (full /home)**: Size of your `/home` + 2-5 GB
- **Max compression**: ~70-80% of normal size

> 💡 **Tip**: Use "Clean session" after copying the ISO to free up space.

---

### Can I use Eggsmaker 2 from another device on the network?

For security, the application listens only on `127.0.0.1` by default. To enable LAN access, start it with:
```
EGGSMAKER_HOST=0.0.0.0 eggsmaker2
```

The footer will show the access address. Use this only on a trusted network: the interface can run administrative actions after `sudo` authentication.

---

### What happens if I close the browser during an operation?

The operation **continues running** in the background. You can:
1. Reopen `http://localhost:8080`
2. View progress in the terminal

---

### Can I cancel an operation in progress?

Currently there is no cancel button. Options:
1. Wait for it to finish
2. Close the application completely (Ctrl+C in terminal)
3. Use "Clean session" afterwards to clean partial files

---

### Where are generated ISOs saved?

Typical locations (checked in order):
1. `/home/eggs/` - Main location
2. `/home/eggs/.mnt/` - Alternative location

Size is automatically shown in the status bar.

---

### Can I create multiple ISOs without cleaning session?

Yes, but:
- ⚠️ The new ISO will overwrite the previous one in `/home/eggs/`
- 💡 **Tip**: Copy each ISO before generating the next one

---

### How do I update Eggsmaker 2?

```bash
cd /path/to/eggsmaker2
git pull
./install.sh  # Reinstall
```

Or if you used pip:
```bash
pip install -e . --upgrade
```

---

## 📚 Additional Resources

### Useful Links

- **Penguins' Eggs**: https://penguins-eggs.net/
- **GitHub Repository**: https://github.com/pieroproietti/penguins-eggs
- **Eggs Documentation**: https://penguins-eggs.net/docs/

### Support

For specific Eggsmaker 2 issues, check:
- The log terminal in the application
- System log files
- Penguins' Eggs community

---

## 📝 Summary

1. Click **Create ISO**.
2. Choose whether to copy desktop settings, edit `eggs.yaml`, and create a system clone.
3. Confirm generation and follow progress in the terminal and progress bar.
4. Select a folder with the integrated picker, copy again, or skip copying.
5. Choose whether to run `eggs kill` to remove temporary build files.

## 🎉 Enjoy creating your ISOs!

This guide covers all main functionalities of Eggsmaker 2. For advanced usage, consult the Penguins' Eggs documentation.

**Guide Version**: 1.0
**Compatible with**: Eggsmaker 2 v26.10.21
