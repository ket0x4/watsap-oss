# Watsap builder

## Requirements
- Go 1.20 or later
- Git

## Usage
- clone watsap repository 
- put builder in the `installer` directory
- Run program and fill in the required fields.
- Click "Build" to generate the installer.


## Build
- Ensure you have the `Fyne` toolkit installed.
- `fyne package --release --os windows/linux`

### Building windows from linux
- Install `mingw-w64`
- `CC=x86_64-w64-mingw32-gcc fyne package --release --os windows`