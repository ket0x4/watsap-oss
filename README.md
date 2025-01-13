# watsap

Remote system administration tool for Windows/Linux

### why watsap?
It's a meme comes from the whatsapp gold and 2 memes. First started as a joke to troll friends & learn python/go but then I decided to make it a real project. Still meme and comes with absolutely no warranty.

### how to use?
1. Clone the repository
2. Create .env file with the following content:
``` bash
export TG_BOT_TOKEN="TOKEN"
export TG_CHAT_ID="CHATID"
export RSHELL_IP="IP" #Optional
export RSHELL_PORT="PORT" #Optional
```

3. Run `build.sh` to build the client
4. Run the client on the target machine

## to-do:
- [x] Add a way to run commands on the remote machine
- [ ] Autoupdate the client
- [x] Fix userid generation
- [x] Modularize the code
- [x] Add keylogger
- [x] Add file scraper
- [x] Add screenshot
- [ ] Add webcam
- [ ] Add microphone
- [ ] Add remote desktop
- [x] Add remote shell
- [x] Better build script

