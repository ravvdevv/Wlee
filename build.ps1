Write-Host "Cleaning Go build cache..." -ForegroundColor Cyan
go clean -cache

Write-Host "Syncing dependencies..." -ForegroundColor Cyan
go mod tidy

Write-Host "Building wlee.exe..." -ForegroundColor Cyan
# -ldflags="-H windowsgui" ensures no console window pops up when running the exe
go build -ldflags="-H windowsgui" -o wlee.exe main.go

if ($?) {
    Write-Host "Build Successful! Run .\wlee.exe to start (remember the random delay)." -ForegroundColor Green
}
else {
    Write-Host "Build failed. Please check the errors above." -ForegroundColor Red
}
