version = 1.3
# installation
## debian or debian based
download the .deb file
## other linux distros
copy int.py to any directory in your path

i recommend ~/.local/bin because it's usually not touched by package managers

## building the debian package
install git to download the repo ```sudo apt update && sudo apt install git```

clone the repo ```git clone https://github.com/WerIstLuka/int```

change directory to int ```cd int```

use the mkdpkg script to build ```mkdpkg```

enter the version number at the top of this page

