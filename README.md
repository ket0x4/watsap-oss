# watsap legacy

This is a first attempt at creating a watsap with python. works fine but getting binary output is a bit of a pain. So I decided to use golang for this project. final binary output is much better than python and binary size reduced from 13mb to 1.8mb. that's a huge difference imo. Thats not all, memory usage reduced from 60mb to only 4mb plus golang is much faster than python. I can get linux and windows binaries from the same code. I am happy with the result. I will be using golang for this project from now on.

## How to use
1. Clone the repo
2. `python -m venv venv`
3. `source venv/bin/activate`
4. `pip install -r requirements.txt`
5. Edit the `watsap.py` file and add your own TOKEN and CHAT_ID
6. `python watsap.py` # Or use pyinstaller to create a binary