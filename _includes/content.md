## Description
simpleShot is a simple tool to take quick screenshots and upload them so that you can share them easily.

---

## Install
Often the dependencies are already installed on your system, but maybe you also need to install them:

#### Archlinux/Antergos/Manjaro
{% highlight console%}
sudo pacman -S xclip imagemagick
{% endhighlight %}

#### Ubuntu/Linux Mint
{% highlight console%}
sudo apt-get install xclip imagemagick
{% endhighlight %}

After that's done installing __simpleShot__ is easy:
{% highlight bash%}
su -c "curl -s https://raw.githubusercontent.com/niklas-heer/simpleShot/master/install.sh | bash"
{% endhighlight %}

---

## Usage
The command line interface:
{% highlight console%}
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
   help, h  Shows a list of commands or help for one command
   
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
{% endhighlight %}

---

## Changelog

### __0.2.0__ - 2015-01-24

#### Added
- The ability to define a `path` on the web-server was added to the `.simpleShot.gcfg` file. This parameter is optional!
- Added debugging messages for the `--upload -u` flag. If debugging of the FTP connection is needed.

### __0.1.0__ - 2015-01-23

#### Added
- The `--debug -d` flag was added to easily display debugging messages. Probably only interesting for developing this tool.

#### Changed
- Renamed the `--format -f` flag to `--type -t`.
- Renamed the `--directory -d` flag to `--folder -f`.

### __0.0.1__ - 2015-01-23
initial release.