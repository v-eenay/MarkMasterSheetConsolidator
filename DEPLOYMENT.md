# Deployment Guide

**Mark Master Sheet Consolidator**
**All rights reserved to Vinay Koirala**

This guide provides step-by-step instructions for deploying and using the Mark Master Sheet Consolidator application.

## Author

**Vinay Koirala**
- Personal Email: koiralavinay@gmail.com
- Professional Email: binaya.koirala@iic.edu.np
- LinkedIn: [veenay](https://linkedin.com/in/veenay)
- GitHub: [v-eenay](https://github.com/v-eenay)

## Prerequisites

### System Requirements
- **Operating System**: Windows 10/11, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **Memory**: Minimum 4GB RAM (8GB recommended for large datasets)
- **Storage**: At least 1GB free space for application and logs
- **Go Runtime**: Go 1.21 or higher (for building from source)

### Required Files Structure
Ensure your directory structure matches:
```
mark-master-sheet/
├── StudentFiles/           # Contains student Excel files
│   ├── L2C1/
│   │   ├── G1/
│   │   ├── G2/
│   │   └── ...
│   ├── L2C2/
│   └── ...
├── MasterSheet/           # Contains master Excel file
│   └── CS5054NT 2024-25 SEM2 Result.xlsx
├── output/                # Will contain updated master sheets
├── logs/                  # Will contain application logs
└── backups/              # Will contain backup files
```

## Installation

### Option 1: Using Pre-built Binary (Recommended)

1. **Download the binary** for your operating system
2. **Extract** to your desired directory
3. **Copy configuration** file:
   ```bash
   cp config.sample.toml config.toml
   ```
4. **Edit configuration** as needed (see Configuration section)

### Option 2: Building from Source

#### Windows
```cmd
# Clone or download source code
# Navigate to project directory
build.bat
```

#### Linux/macOS
```bash
# Clone or download source code
# Navigate to project directory
chmod +x build.sh
./build.sh
```

## Configuration

### Basic Configuration

Edit `config.toml` to match your environment:

```toml
[paths]
student_files_folder = "./StudentFiles"
master_sheet_path = "./MasterSheet/CS5054NT 2024-25 SEM2 Result.xlsx"
output_folder = "./output"
log_folder = "./logs"
backup_folder = "./backups"
```

### Excel Settings

Verify these match your Excel file structure:

```toml
[excel_settings]
student_worksheet_name = "Grading Sheet"
master_worksheet_name = "001"
student_id_cell = "B2"
mark_cells = ["C6", "C7", "C8", "C9", "C10", "C11", "C12", "C13", "C15", "C16", "C17", "C18", "C19", "C20"]
master_columns = ["I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"]
```

### Performance Tuning

Adjust based on your system capabilities:

```toml
[processing]
max_concurrent_files = 10    # Reduce if system is slow
timeout_seconds = 300        # Increase for large datasets
retry_attempts = 3           # Increase for unreliable storage
```

## Usage

### Basic Operations

#### 1. Test Configuration
```bash
# Windows
mark-master-sheet.exe -stats

# Linux/macOS
./mark-master-sheet -stats
```

#### 2. Dry Run (Recommended First)
```bash
# Windows
mark-master-sheet.exe -dry-run

# Linux/macOS
./mark-master-sheet -dry-run
```

#### 3. Production Run
```bash
# Windows
mark-master-sheet.exe

# Linux/macOS
./mark-master-sheet
```

### Advanced Usage

#### Custom Configuration
```bash
mark-master-sheet.exe -config custom.toml
```

#### Automated Execution
Create a batch file for regular execution:

**Windows (run.bat):**
```cmd
@echo off
cd /d "C:\path\to\mark-master-sheet"
mark-master-sheet.exe
if %errorlevel% neq 0 (
    echo Processing failed with error %errorlevel%
    pause
)
```

**Linux/macOS (run.sh):**
```bash
#!/bin/bash
cd "/path/to/mark-master-sheet"
./mark-master-sheet
if [ $? -ne 0 ]; then
    echo "Processing failed with error $?"
    exit 1
fi
```

## Monitoring and Troubleshooting

### Log Files

Logs are stored in the configured log folder:
- **Daily logs**: `mark-master-sheet-YYYY-MM-DD.log`
- **Automatic rotation**: Based on size and age settings
- **Log levels**: DEBUG, INFO, WARN, ERROR

### Common Issues

#### 1. "Master sheet not found"
- Verify the file path in `config.toml`
- Ensure the file exists and is accessible
- Check file permissions

#### 2. "Student ID not found"
- Review student ID format in both files
- Check for case sensitivity issues
- Verify the student ID cell location

#### 3. "Permission denied"
- Close Excel files if open in other applications
- Run as administrator if necessary
- Check file and folder permissions

#### 4. "Worksheet not found"
- Verify worksheet names in configuration
- Check if student files have the expected structure
- Ensure worksheet names match exactly

### Performance Optimization

#### For Large Datasets
1. **Increase memory allocation**:
   ```bash
   # Set Go memory limit
   export GOMEMLIMIT=4GiB
   ```

2. **Adjust concurrency**:
   ```toml
   max_concurrent_files = 5  # Reduce for slower systems
   ```

3. **Enable file skipping**:
   ```toml
   skip_invalid_files = true
   ```

#### For Network Storage
1. **Increase timeout**:
   ```toml
   timeout_seconds = 600
   ```

2. **Increase retry attempts**:
   ```toml
   retry_attempts = 5
   ```

## Backup and Recovery

### Automatic Backups
- Enabled by default before any changes
- Stored in configured backup folder
- Timestamped for easy identification

### Manual Backup
```bash
# Create manual backup before processing
cp "MasterSheet/CS5054NT 2024-25 SEM2 Result.xlsx" "backups/manual_backup_$(date +%Y%m%d_%H%M%S).xlsx"
```

### Recovery
1. **Locate backup file** in backup folder
2. **Copy back to master location**:
   ```bash
   cp "backups/backup_file.xlsx" "MasterSheet/CS5054NT 2024-25 SEM2 Result.xlsx"
   ```

## Scheduling

### Windows Task Scheduler
1. Open Task Scheduler
2. Create Basic Task
3. Set trigger (daily, weekly, etc.)
4. Set action to run `mark-master-sheet.exe`
5. Configure conditions and settings

### Linux/macOS Cron
```bash
# Edit crontab
crontab -e

# Add entry (example: daily at 2 AM)
0 2 * * * cd /path/to/mark-master-sheet && ./mark-master-sheet >> logs/cron.log 2>&1
```

## Security Considerations

1. **File Permissions**: Ensure only authorized users can access Excel files
2. **Backup Security**: Store backups in secure location
3. **Log Security**: Logs may contain student information
4. **Network Access**: If using network storage, ensure secure connections

## Support and Maintenance

### Regular Maintenance
1. **Monitor log files** for errors and warnings
2. **Clean old backups** periodically
3. **Update configuration** as needed
4. **Test with dry-run** before major changes

### Getting Help
1. **Check logs** for detailed error messages
2. **Run with debug logging**:
   ```toml
   [logging]
   level = "DEBUG"
   ```
3. **Use dry-run mode** to test changes safely
4. **Review configuration** for common issues
