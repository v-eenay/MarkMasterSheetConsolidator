# Mark Master Sheet Consolidator

A production-ready Go application that automates the consolidation of student marks from individual Excel files into a master spreadsheet.

## Features

- **Recursive File Discovery**: Scans nested folders to find all Excel files (.xlsx, .xls)
- **Concurrent Processing**: Processes multiple files simultaneously with configurable rate limiting
- **Data Validation**: Validates student IDs and numeric marks with comprehensive error handling
- **Backup System**: Creates timestamped backups before making changes
- **Comprehensive Logging**: Structured logging with rotation and multiple output formats
- **Progress Tracking**: Real-time progress indicators and processing statistics
- **Configuration Management**: TOML-based configuration with validation
- **Graceful Shutdown**: Handles interruption signals properly
- **Dry-Run Mode**: Test processing without making actual changes

## Quick Start

### Prerequisites

- Go 1.21 or higher (for building from source)
- Excel files with the expected structure

### Installation

#### Option 1: Build from Source

**Windows:**
```cmd
# Run the build script
build.bat
```

**Linux/macOS:**
```bash
# Make script executable and run
chmod +x build.sh
./build.sh
```

#### Option 2: Manual Build
```bash
# Download dependencies
go mod tidy

# Build the application
go build -o mark-master-sheet cmd/main.go

# Or for Windows
go build -o mark-master-sheet.exe cmd/main.go
```

### First Run

1. **Copy configuration file:**
   ```bash
   cp config.sample.toml config.toml
   ```

2. **Edit configuration** to match your file paths

3. **Test with dry run:**
   ```bash
   # Windows
   mark-master-sheet.exe -dry-run

   # Linux/macOS
   ./mark-master-sheet -dry-run
   ```

4. **Run actual processing:**
   ```bash
   # Windows
   mark-master-sheet.exe

   # Linux/macOS
   ./mark-master-sheet
   ```

## Configuration

The application uses a `config.toml` file for configuration. Copy the provided example and modify as needed:

```toml
[paths]
student_files_folder = "./StudentFiles"
master_sheet_path = "./MasterSheet/CS5054NT 2024-25 SEM2 Result.xlsx"
output_folder = "./output"
log_folder = "./logs"
backup_folder = "./backups"

[excel_settings]
student_worksheet_name = "Grading Sheet"
master_worksheet_name = "001"
student_id_cell = "B2"
mark_cells = ["C6", "C7", "C8", "C9", "C10", "C11", "C12", "C13", "C15", "C16", "C17", "C18", "C19", "C20"]
master_columns = ["I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V"]

[processing]
max_concurrent_files = 10
backup_enabled = true
skip_invalid_files = true
timeout_seconds = 300
retry_attempts = 3

[logging]
level = "INFO"
console_output = true
file_output = true
max_file_size_mb = 100
max_backup_files = 5
max_age_days = 30
```

## Usage

### Basic Usage

```bash
# Run with default configuration
./mark-master-sheet

# Use custom configuration file
./mark-master-sheet -config /path/to/config.toml

# Dry run (no changes made)
./mark-master-sheet -dry-run

# Show processing statistics
./mark-master-sheet -stats

# Show version
./mark-master-sheet -version
```

### Command Line Options

- `-config string`: Path to configuration file (default: "config.toml")
- `-dry-run`: Run in dry-run mode (no actual changes)
- `-stats`: Show processing statistics and exit
- `-version`: Show version information

## Excel File Structure

### Student Files

Each student Excel file should have:
- A worksheet named "Grading Sheet" (configurable)
- Student ID in cell B2
- Marks in cells: C6, C7, C8, C9, C10, C11, C12, C13, C15, C16, C17, C18, C19, C20

### Master Sheet

The master Excel file should have:
- A worksheet named "001" (configurable)
- Student IDs in column B
- Target columns for marks: I, J, K, L, M, N, O, P, Q, R, S, T, U, V

### Mark Mapping

```
Student File → Master Sheet
C6  → Column I     C11 → Column N     C16 → Column R
C7  → Column J     C12 → Column O     C17 → Column S
C8  → Column K     C13 → Column P     C18 → Column T
C9  → Column L     C15 → Column Q     C19 → Column U
C10 → Column M                        C20 → Column V
```

## Directory Structure

```
mark-master-sheet/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go          # Configuration management
│   ├── excel/
│   │   ├── reader.go          # Excel file reading
│   │   └── writer.go          # Excel file writing
│   ├── processor/
│   │   └── processor.go       # Main processing logic
│   └── logger/
│       └── logger.go          # Logging utilities
├── pkg/
│   └── models/
│       └── student.go         # Data structures
├── config.toml                # Configuration file
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Logging

The application provides comprehensive logging:

- **Console Output**: Real-time progress and status updates
- **File Output**: Detailed logs with rotation (daily files, size limits)
- **Log Levels**: INFO, WARN, ERROR with configurable levels
- **Structured Logging**: JSON-like format with contextual fields

Log files are stored in the configured log folder with automatic rotation.

## Error Handling

The application handles various error scenarios:

- **File Access Errors**: Corrupted or locked Excel files
- **Validation Errors**: Invalid student IDs or marks
- **Missing Data**: Empty cells or missing worksheets
- **Network Issues**: File system access problems
- **Memory Issues**: Large file processing with appropriate limits

## Performance

- **Concurrent Processing**: Configurable number of simultaneous file operations
- **Memory Efficient**: Streams large files without loading everything into memory
- **Progress Tracking**: Real-time progress indicators
- **Timeout Handling**: Configurable timeouts to prevent hanging

## Troubleshooting

### Common Issues

1. **"Master sheet not found"**
   - Verify the master sheet path in config.toml
   - Ensure the file exists and is accessible

2. **"Student ID not found"**
   - Check if student IDs match between files and master sheet
   - Review case sensitivity and formatting

3. **"Permission denied"**
   - Ensure the application has read/write permissions
   - Close Excel files if they're open in other applications

4. **"Worksheet not found"**
   - Verify worksheet names in configuration
   - Check if student files have the expected worksheet structure

### Debug Mode

Enable debug logging by setting the log level to "DEBUG" in config.toml:

```toml
[logging]
level = "DEBUG"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Project Files

### Core Application Files
- `cmd/main.go` - Application entry point with CLI handling
- `internal/config/config.go` - Configuration management and validation
- `internal/excel/reader.go` - Excel file reading and data extraction
- `internal/excel/writer.go` - Excel file writing and master sheet updates
- `internal/processor/processor.go` - Main processing logic and concurrency
- `internal/logger/logger.go` - Structured logging with rotation
- `pkg/models/student.go` - Data structures and validation

### Configuration and Scripts
- `config.toml` - Main configuration file
- `config.sample.toml` - Sample configuration with documentation
- `build.bat` / `build.sh` - Build scripts for Windows/Unix
- `test.bat` / `test.sh` - Testing scripts

### Documentation
- `README.md` - This file
- `DEPLOYMENT.md` - Comprehensive deployment guide
- `.gitignore` - Git ignore rules

## Support

For support and questions:
- Check the logs in the configured log folder
- Review the configuration file for correct paths and settings
- Use dry-run mode to test changes safely
- See `DEPLOYMENT.md` for detailed troubleshooting guide
#   M a r k M a s t e r S h e e t C o n s o l i d a t o r  
 