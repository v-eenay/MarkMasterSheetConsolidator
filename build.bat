@echo off
echo Building Mark Master Sheet Consolidator...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

REM Download dependencies
echo Downloading dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo Error: Failed to download dependencies
    pause
    exit /b 1
)

REM Run tests
echo Running tests...
go test ./...
if %errorlevel% neq 0 (
    echo Warning: Some tests failed
)

REM Build the application
echo Building application...
go build -o mark-master-sheet.exe cmd/main.go
if %errorlevel% neq 0 (
    echo Error: Build failed
    pause
    exit /b 1
)

echo Build completed successfully!
echo Executable: mark-master-sheet.exe
echo.
echo Usage:
echo   mark-master-sheet.exe                    # Run with default config
echo   mark-master-sheet.exe -dry-run           # Test run without changes
echo   mark-master-sheet.exe -stats             # Show statistics
echo   mark-master-sheet.exe -config custom.toml # Use custom config
echo.
pause
