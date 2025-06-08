# Mark Master Sheet Consolidator - Project Summary

**All rights reserved to Vinay Koirala**

## Author Information

**Vinay Koirala**
- Personal Email: koiralavinay@gmail.com
- Professional Email: binaya.koirala@iic.edu.np
- LinkedIn: [veenay](https://linkedin.com/in/veenay)
- GitHub: [v-eenay](https://github.com/v-eenay)
- Repository: https://github.com/v-eenay/MarkMasterSheetConsolidator.git

## Overview

A production-ready Go application that automates the consolidation of student marks from individual Excel files into a master spreadsheet. The application is designed for educational institutions to efficiently process large numbers of student grade files.

## ✅ Completed Features

### Core Functionality
- ✅ **Recursive File Discovery**: Scans nested folders (unlimited depth) for Excel files
- ✅ **Multi-format Support**: Handles both .xlsx and .xls files
- ✅ **Data Extraction**: Extracts Student ID and marks from specific cells
- ✅ **Data Validation**: Validates student IDs (alphanumeric) and numeric marks (0-100 range)
- ✅ **Master Sheet Updates**: Updates master sheet with configurable cell mapping
- ✅ **Concurrent Processing**: Processes multiple files simultaneously with rate limiting

### Production Features
- ✅ **Comprehensive Error Handling**: Custom error types for different failure scenarios
- ✅ **Backup System**: Creates timestamped backups before modifications
- ✅ **Structured Logging**: Multi-level logging with file rotation
- ✅ **Configuration Management**: TOML-based configuration with validation
- ✅ **Progress Tracking**: Real-time progress bars and status updates
- ✅ **Graceful Shutdown**: Handles SIGINT/SIGTERM signals properly
- ✅ **Retry Logic**: Configurable retry attempts for failed operations
- ✅ **Timeout Handling**: Prevents hanging operations

### Command Line Interface
- ✅ **Dry-run Mode**: Test processing without making changes
- ✅ **Statistics Mode**: Show processing statistics
- ✅ **Custom Configuration**: Support for custom config files
- ✅ **Version Information**: Display application version

### Testing and Quality
- ✅ **Unit Tests**: Comprehensive test coverage for models and configuration
- ✅ **Build Scripts**: Automated build scripts for Windows and Unix
- ✅ **Test Scripts**: Automated testing scripts
- ✅ **Documentation**: Comprehensive README and deployment guides

## 📁 Project Structure

```
mark-master-sheet/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go          # Configuration management
│   │   └── config_test.go     # Configuration tests
│   ├── excel/
│   │   ├── reader.go          # Excel file reading operations
│   │   └── writer.go          # Excel file writing operations
│   ├── processor/
│   │   └── processor.go       # Main processing logic
│   └── logger/
│       └── logger.go          # Logging utilities
├── pkg/
│   └── models/
│       ├── student.go         # Student data structures
│       └── student_test.go    # Model tests
├── config.toml                # Main configuration file
├── config.sample.toml         # Sample configuration
├── go.mod                     # Go module definition
├── build.bat / build.sh       # Build scripts
├── test.bat / test.sh         # Test scripts
├── README.md                  # Main documentation
├── DEPLOYMENT.md              # Deployment guide
└── .gitignore                 # Git ignore rules
```

## 🔧 Technical Implementation

### Dependencies
- **github.com/xuri/excelize/v2**: Excel file operations
- **github.com/BurntSushi/toml**: Configuration file parsing
- **github.com/sirupsen/logrus**: Structured logging
- **github.com/schollz/progressbar/v3**: Progress indicators
- **gopkg.in/natefinch/lumberjack.v2**: Log file rotation

### Architecture
- **Modular Design**: Clear separation of concerns
- **Error Handling**: Custom error types with context
- **Concurrency**: Goroutines with semaphore-based rate limiting
- **Configuration**: Centralized TOML-based configuration
- **Logging**: Structured logging with multiple outputs

### Performance Features
- **Concurrent Processing**: Configurable number of simultaneous operations
- **Memory Efficient**: Streams files without loading everything into memory
- **Rate Limiting**: Prevents system overload
- **Progress Tracking**: Real-time progress indicators
- **Timeout Handling**: Configurable timeouts

## 📊 Data Processing

### Input Format (Student Files)
- **Worksheet**: "Grading Sheet"
- **Student ID**: Cell B2
- **Marks**: Cells C6, C7, C8, C9, C10, C11, C12, C13, C15, C16, C17, C18, C19, C20

### Output Format (Master Sheet)
- **Worksheet**: "001"
- **Student IDs**: Column B
- **Mark Columns**: I, J, K, L, M, N, O, P, Q, R, S, T, U, V

### Data Validation
- **Student ID**: Alphanumeric characters only, non-empty
- **Marks**: Numeric values between 0-100
- **Empty Cells**: Handled gracefully (stored as -1)

## 🚀 Usage Examples

### Basic Usage
```bash
# Windows
mark-master-sheet.exe

# Linux/macOS
./mark-master-sheet
```

### Advanced Usage
```bash
# Dry run (no changes)
mark-master-sheet.exe -dry-run

# Show statistics
mark-master-sheet.exe -stats

# Custom configuration
mark-master-sheet.exe -config custom.toml

# Show version
mark-master-sheet.exe -version
```

## 📈 Performance Characteristics

### Scalability
- **Concurrent Files**: Up to 10 simultaneous file operations (configurable)
- **File Size**: Efficiently handles large Excel files
- **Dataset Size**: Tested with hundreds of student files
- **Memory Usage**: Optimized for minimal memory footprint

### Error Recovery
- **Retry Logic**: Automatic retry for transient failures
- **Skip Invalid**: Option to skip corrupted files and continue
- **Backup System**: Automatic backups before any changes
- **Graceful Degradation**: Continues processing despite individual file failures

## 🔒 Security and Reliability

### Data Protection
- **Automatic Backups**: Timestamped backups before modifications
- **Validation**: Comprehensive input validation
- **Error Logging**: Detailed error tracking without exposing sensitive data
- **File Permissions**: Respects existing file permissions

### Reliability Features
- **Graceful Shutdown**: Handles interruption signals
- **Timeout Protection**: Prevents hanging operations
- **Resource Management**: Proper cleanup of file handles
- **Error Recovery**: Comprehensive error handling and recovery

## 📝 Configuration Options

### Paths Configuration
- Student files folder (recursive scanning)
- Master sheet file path
- Output directory for updated files
- Log directory for application logs
- Backup directory for safety copies

### Processing Configuration
- Maximum concurrent files
- Backup enabled/disabled
- Skip invalid files option
- Timeout settings
- Retry attempts

### Logging Configuration
- Log level (DEBUG, INFO, WARN, ERROR)
- Console output enabled/disabled
- File output enabled/disabled
- Log rotation settings

## 🎯 Key Benefits

1. **Automation**: Eliminates manual copy-paste operations
2. **Accuracy**: Reduces human errors in data consolidation
3. **Speed**: Processes hundreds of files in minutes
4. **Safety**: Automatic backups and validation
5. **Scalability**: Handles large datasets efficiently
6. **Reliability**: Comprehensive error handling and recovery
7. **Monitoring**: Detailed logging and progress tracking
8. **Flexibility**: Configurable for different file structures

## 🔄 Next Steps for Deployment

1. **Install Go** (if building from source)
2. **Copy project files** to target directory
3. **Configure** `config.toml` for your environment
4. **Test** with dry-run mode
5. **Run** actual processing
6. **Monitor** logs for any issues
7. **Schedule** for regular execution (optional)

## 📞 Support

- **Documentation**: Comprehensive README and deployment guides
- **Logging**: Detailed logs for troubleshooting
- **Testing**: Dry-run mode for safe testing
- **Configuration**: Sample configuration with documentation
- **Error Messages**: Clear, actionable error messages

This application is ready for production use and provides a robust solution for automated mark consolidation in educational environments.
