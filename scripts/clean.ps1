Write-Host "Cleaning up watsap installation"

# require admin rights
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    $arguments = "& '" + $myinvocation.mycommand.definition + "'"
    Start-Process powershell -Verb runAs -ArgumentList $arguments
    Break
  }

# set title and window size
$host.ui.RawUI.WindowTitle = "watsap cleaner"
$host.ui.RawUI.WindowSize = New-Object System.Management.Automation.Host.Size(100, 30)


$waDir = "$env:APPDATA\watsap"
$shellAutoStartDir = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup\"

if (-NOT (Test-Path -Path $waDir)) {
    Write-Host "watsap is not installed"
    exit
}

$waBinNames = @("*efender*, *atsap*")

foreach ($bin in $waBinNames) {
  Stop-Process -Name $bin -Force -ErrorAction SilentlyContinue -Confirm:$false
  Remove-Item -Path $shellAutoStartDir\$bin -Force -ErrorAction SilentlyContinue -Recurse -Confirm:$false
}

Remove-Item -Path $waDir -Force -ErrorAction SilentlyContinue -Recurse -Confirm:$false
Remove-MpPreference -ExclusionPath $shellAutoStartDir -ErrorAction SilentlyContinue
Remove-MpPreference -ExclusionPath $waDir -ErrorAction SilentlyContinue
Write-Host "watsap installation cleaned up"
Pause