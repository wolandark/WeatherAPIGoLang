Write-Host "Fixing import statements to match go.mod..." -ForegroundColor Green

if (-not (Test-Path "go.mod")) {
    Write-Host "Error: go.mod not found. Make sure you're in the project root directory." -ForegroundColor Red
    exit 1
}

$goModContent = Get-Content "go.mod" -Raw
if ($goModContent -match "module\s+(\S+)") {
    $moduleName = $matches[1]
    Write-Host "Found module name: $moduleName" -ForegroundColor Cyan
} else {
    Write-Host "Error: Could not find module name in go.mod" -ForegroundColor Red
    exit 1
}

$goFiles = Get-ChildItem -Recurse -Filter "*.go"

foreach ($file in $goFiles) {
    Write-Host "Updating imports in: $($file.FullName)" -ForegroundColor Yellow
    
    $content = Get-Content $file.FullName -Raw
    $updatedContent = $content -replace 'weather-api/', "$moduleName/"
    
    $updatedContent | Out-File -FilePath $file.FullName -Encoding UTF8 -NoNewline
}

Write-Host "`nImport statements updated successfully!" -ForegroundColor Green
Write-Host "You can now run:" -ForegroundColor Green
Write-Host "go run .\cmd\server\main.go" -ForegroundColor White
