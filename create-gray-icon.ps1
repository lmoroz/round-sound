# PowerShell script to create a grayscale version of an ICO file
# Uses .NET System.Drawing

Add-Type -AssemblyName System.Drawing

$sourceIco = "app\tray.ico"
$targetIco = "app\tray-gray.ico"

# Load the icon
$icon = New-Object System.Drawing.Icon($sourceIco)
$bitmap = $icon.ToBitmap()

# Create grayscale version
$grayBitmap = New-Object System.Drawing.Bitmap($bitmap.Width, $bitmap.Height)

for ($x = 0; $x -lt $bitmap.Width; $x++) {
    for ($y = 0; $y -lt $bitmap.Height; $y++) {
        $pixel = $bitmap.GetPixel($x, $y)
        
        # Calculate grayscale value (luminosity method)
        $gray = [int]($pixel.R * 0.299 + $pixel.G * 0.587 + $pixel.B * 0.114)
        
        # Create grayscale color with original alpha
        $grayColor = [System.Drawing.Color]::FromArgb($pixel.A, $gray, $gray, $gray)
        
        $grayBitmap.SetPixel($x, $y, $grayColor)
    }
}

# Save as PNG first (ICO conversion is complex)
$tempPng = "app\tray-gray-temp.png"
$grayBitmap.Save($tempPng, [System.Drawing.Imaging.ImageFormat]::Png)

# Convert PNG to ICO using icon handle
$grayIcon = [System.Drawing.Icon]::FromHandle($grayBitmap.GetHicon())
$stream = New-Object System.IO.FileStream($targetIco, [System.IO.FileMode]::Create)
$grayIcon.Save($stream)
$stream.Close()

# Cleanup
$bitmap.Dispose()
$grayBitmap.Dispose()
$icon.Dispose()
$grayIcon.Dispose()

if (Test-Path $tempPng) {
    Remove-Item $tempPng
}

Write-Host "Grayscale icon created: $targetIco" -ForegroundColor Green
