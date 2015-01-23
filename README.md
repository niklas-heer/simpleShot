# simpleShot

___Version___: 0.0.1 <br>
___Author___: Niklas Heer<br>
___License___: GPLv3

A simple tool to take screenshots and upload them to an FTP server.

```
$ simpleShot -h
NAME:
   simpleShot - Takes a screenshot, uploads it via FTP and copies the url into your clipboard!

USAGE:
   simpleShot [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
  Niklas Heer - <niklas.heer@gmail.com>

COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --select, -s                      Select the area for the screenshot.
   --upload, -u                      Reads the credentials under ~/.simpleShot.gcfg and uploads it.
   --quiet, -q                       Don't notify me!
   --name-length, --nl "6"           Choose the length for the name generator.
   --name-alphabet, --na "alphanum"  Choose the alphabet for the name generator.
   --directory, -d "screenshots"     Choose the direcotry where the screenshots are saved. (e.g. /home/nh/screenshots)
   --format, -f "jpg"                Choose the format of the screenshot. (png, jpg...)
   --help, -h                        show help
   --version, -v                     print the version
```


## Installation


### 1. Getting the requirements
In order to use the programm you need the following programms install on your GNU/Linux box:

* imagemagick
* xclip

##### Archlinux/Antergos/Manjaro
`sudo pacman -S xclip imagemagick`

##### Ubuntu/Linux Mint
`sudo apt-get install xclip imagemagick`

### 2. Installing simpleShot

#### Global
```
cd ~/Downloads
wget https://raw.github.com/niklas-heer/simpleShot/master/simpleShot
chmod +x simpleShot
sudo mv simpleShot /usr/local/bin/
```

#### As a script
```
cd ~/Downloads
wget https://raw.github.com/niklas-heer/simpleShot/master/simpleShot
chmod +x simpleShot
mkdir ~/bin
mv simpleShot ~/bin/
```

Now you can bind it to an keyboard shortcut in your Preferences. This could look like this:

![](assets/screenshot.png?raw=true)

### 3. Configuration

There is not much to configure but if you want to upload to an FTP server you need to edit your config file at `~/.simpleShot.gcfg`.
It should look like this:
```
[ftp]
url = http://screens.your-server.org/
server = your-server.org
port = 21
user = youruser
pw = yourpassword
```

## Additional stuff

### A convient index.php script
You can take this script and put it in the same directory as your screenshot.
```
<?php

$fileSize = 0;
$filesArray = glob('*.{png,jpg}', GLOB_BRACE);

foreach ($filesArray as $file) {
  $fileSize = $fileSize + (new SplFileInfo($file))->getSize();
}

echo count($filesArray) . ' Bilder verbrauchen etwa ' . round(($fileSize / 1024) / 1024, 2) . ' MB Speicherplatz.';
```

which will look like this:

![](assets/screenshot2.png?raw=true)
