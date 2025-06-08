# Mark Master Sheet Consolidator

Automates consolidation of student marks from individual Excel files into a master spreadsheet.

**All rights reserved to Vinay Koirala**

## Author

**Vinay Koirala**
- Email: koiralavinay@gmail.com
- Professional: binaya.koirala@iic.edu.np
- LinkedIn: [veenay](https://linkedin.com/in/veenay)
- GitHub: [v-eenay](https://github.com/v-eenay)
- Repository: https://github.com/v-eenay/MarkMasterSheetConsolidator.git

## Quick Start

### Build
```bash
# Windows
build.bat

# Linux/macOS
chmod +x build.sh && ./build.sh
```

### Usage

**Command Line (Always Available):**
```bash
# Windows
mark-master-sheet.exe -dry-run    # Test run
mark-master-sheet.exe             # Process files

# Linux/macOS
./mark-master-sheet -dry-run      # Test run
./mark-master-sheet               # Process files
```

**GUI (If Build Succeeds):**
```bash
# Windows
mark-master-sheet-gui.exe

# Linux/macOS
./mark-master-sheet-gui
```

*Note: GUI requires OpenGL/CGO dependencies. If GUI build fails, use the fully functional CLI version.*

## Configuration

Copy `config.sample.toml` to `config.toml` and edit paths to match your files. The GUI provides an easy interface for configuration.

## License

**All rights reserved to Vinay Koirala**

This software is proprietary and confidential. For licensing inquiries, contact koiralavinay@gmail.com.
