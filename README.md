# simpleShot

___Version___: 0.2.0 <br>
___Author___: Niklas Heer<br>
___License___: GPLv3

A simple tool to take screenshots and upload them to an FTP server written in [Go](https://golang.org/).

This project is uses [semantic versioning](http://semver.org/).

For recent changes look at the [change log](https://github.com/niklas-heer/simpleShot/blob/master/CHANGELOG.md).<br>
This project supports the effort of [keepachangelog.com](http://keepachangelog.com/) and follows its guidelines in the change log. We are also using [chag](https://github.com/mtdowling/chag) to automate the process.


```
$ simpleShot -h
NAME:
   simpleShot - Takes a screenshot, uploads it via FTP and copies the url into your clipboard!

USAGE:
   simpleShot [global options] command [command options] [arguments...]

VERSION:
   0.2.0

AUTHOR:
  Niklas Heer - <niklas.heer@gmail.com>

COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --select, -s                     Select the area for the screenshot.
   --upload, -u                     Reads the credentials under ~/.simpleShot.gcfg and uploads it.
   --quiet, -q                      Don't notify me!
   --debug, -d                      Enable debugging.
   --name-length, --nl "6"          Choose the length for the name generator.
   --name-alphabet, --na "alphanum" Choose the alphabet for the name generator.
   --folder, -f "screenshots"       Choose the directory where the screenshots are saved. (e.g. /home/nh/screenshots)
   --type, -t "jpg"                 Choose the format of the screenshot. (png, jpg...)
   --help, -h                       show help
   --version, -v                    print the version
```


Note:<br>
This program is still in its early development phase and while it's not being released as "1.0.0" flags might be renamed or changed.

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

#### Global (recommended)

One easy installation method is through cURL using the following:

```
curl -s https://raw.githubusercontent.com/niklas-heer/simpleShot/master/install.sh | bash
```

You can customize the install source, directory and profile using the SIMPLESHOT_DIR and SIMPLESHOT_VERSION environment variables. The script will place simpleShot in /usr/local/bin by default.

#### As a script
```
cd ~/Downloads
wget https://raw.github.com/niklas-heer/simpleShot/master/simpleShot
chmod +x simpleShot
mkdir ~/bin
mv simpleShot ~/bin/
```

#### From source
This tool assumes you are working in a standard Go workspace, as described in [http://golang.org/doc/code.html](http://golang.org/doc/code.html).

```
go get github.com/tools/godep
go get github.com/niklas-heer/simpleShot
cd $GOPATH/src/github.com/niklas-heer/simpleShot
godep restore
go build
go install
```

Now you can bind it to an keyboard shortcut in your Preferences. This could look like this:

![](assets/screenshot.png?raw=true)

### 3. Configuration

There is not much to configure but if you want to upload to an FTP server you need to edit your config file at `~/.simpleShot.gcfg`.
It should look like this:
```
[ftp]
url = http://screens.your-server.org/with/folder/
server = your-server.org
path = path/to/files
port = 21
user = youruser
pw = yourpassword
```

Note:<br>
`path` is only needed if you don't want to save the files in the root directory of your web-server!

## Additional stuff

### A convient index.php script
You can take this script and put it in the same directory as your screenshots.
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
