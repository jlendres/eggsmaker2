# User Guide: First Installation and Portability

This guide details the files required to port this project to other computers and how to perform the installation.

## 1. Files Required for Porting

To install this project on other computers, you must copy the entire project folder. Ensure it includes the following essential files and folders:

### Source Code and Resources
*   **`main.py`**: Main application file.
*   **`backend.py`**: System logic.
*   **`version.py`**: Application version.
*   **`assets/`**: Folder with images and resources (logo, etc.).

### Installation
*   **`install.sh`**: Automated installation script.
*   **`requirements.txt`**: Python dependencies.
*   **`setup.py`** and **`pyproject.toml`**: Package configuration.
*   **`eggsmaker2.desktop`**: Desktop shortcut for the menu.

---

## 2. Installation Instructions

Once the project folder is copied to the new computer, you have two installation options depending on your needs.

### Option A: Local Installation (Recommended)
Ideal for a single user. Does not affect the global operating system.

1.  Open a terminal inside the project folder.
2.  Run the following command:
    ```bash
    ./install.sh
    ```

**File Location:**
*   The program will be installed in: `~/.local/share/eggsmaker2/`
*   The executable will be in: `~/.local/bin/eggsmaker2`

### Option B: System Installation
Makes the program available to all users on the computer. Requires administrator permissions.

1.  Open a terminal inside the project folder.
2.  Run the following command with `sudo`:
    ```bash
    sudo ./install.sh --system
    ```

**File Location:**
*   The program will be installed in: `/opt/eggsmaker2/`
*   The executable will be in: `/usr/local/bin/eggsmaker2`

---

## 3. Execution

Once installed, you can search for **"Eggsmaker 2"** in your system's application menu or run it from the terminal with the command:

```bash
eggsmaker2
```

---

## 4. Updating

If you need to update the application to a new version:

1.  Download or copy the new version of the project.
2.  Repeat the installation step (Option A or B) you originally used.
    *   The installer will automatically overwrite the old files with the new ones.
    *   You will not lose personal configurations if they are stored outside the program folder.

---

## 5. Uninstallation

To completely remove the application:

1.  Open a terminal in the project folder (if you still have it) or download the `uninstall.sh` script.
2.  Run the uninstaller according to your installation type:

**If you installed locally (Option A):**
```bash
./uninstall.sh
```

**If you installed system-wide (Option B):**
```bash
sudo ./uninstall.sh --system
```


