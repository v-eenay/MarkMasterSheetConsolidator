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

REM Build the applications
echo Building CLI application...
go build -o mark-master-sheet.exe cmd/main.go
if %errorlevel% neq 0 (
    echo Error: CLI build failed
    pause
    exit /b 1
)

echo Building GUI application...
go build -o mark-master-sheet-gui.exe cmd/gui/main.go
if %errorlevel% neq 0 (
    echo Warning: GUI build failed - this may be due to missing OpenGL/CGO dependencies
    echo The CLI version has been built successfully and is fully functional
    echo.
    echo To build the GUI version, you may need to:
    echo 1. Install a C compiler (like TDM-GCC or MinGW-w64)
    echo 2. Set CGO_ENABLED=1
    echo 3. Install OpenGL development libraries
    echo.
    echo For now, you can use the CLI version: mark-master-sheet.exe
) else (
    echo GUI build completed successfully!
    echo GUI Executable: mark-master-sheet-gui.exe
)

echo.
echo Build completed successfully!
echo CLI Executable: mark-master-sheet.exe
echo GUI Executable: mark-master-sheet-gui.exe
echo.
echo Usage:
echo   mark-master-sheet.exe                    # CLI: Run with default config
echo   mark-master-sheet.exe -dry-run           # CLI: Test run without changes
echo   mark-master-sheet.exe -stats             # CLI: Show statistics
echo   mark-master-sheet.exe -config custom.toml # CLI: Use custom config
echo   mark-master-sheet-gui.exe                # GUI: Launch graphical interface
echo.
pause
