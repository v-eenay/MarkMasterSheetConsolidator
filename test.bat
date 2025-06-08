@echo off
echo Testing Mark Master Sheet Consolidator...

REM Check if executable exists
if not exist "mark-master-sheet.exe" (
    echo Error: mark-master-sheet.exe not found
    echo Please run build.bat first
    pause
    exit /b 1
)

REM Check if config file exists
if not exist "config.toml" (
    echo Warning: config.toml not found, copying from sample
    copy config.sample.toml config.toml
)

echo.
echo === Test 1: Show Version ===
mark-master-sheet.exe -version

echo.
echo === Test 2: Show Statistics ===
mark-master-sheet.exe -stats

echo.
echo === Test 3: Dry Run ===
echo This will scan files without making changes...
mark-master-sheet.exe -dry-run

echo.
echo === Test 4: Validate Configuration ===
echo Checking if all required directories and files exist...

if not exist "StudentFiles" (
    echo Warning: StudentFiles directory not found
) else (
    echo ✓ StudentFiles directory exists
)

if not exist "MasterSheet" (
    echo Warning: MasterSheet directory not found
) else (
    echo ✓ MasterSheet directory exists
)

if not exist "output" (
    echo Creating output directory...
    mkdir output
)

if not exist "logs" (
    echo Creating logs directory...
    mkdir logs
)

if not exist "backups" (
    echo Creating backups directory...
    mkdir backups
)

echo.
echo Testing completed!
echo.
echo To run the actual processing:
echo   mark-master-sheet.exe
echo.
echo To run with custom config:
echo   mark-master-sheet.exe -config custom.toml
echo.
pause
