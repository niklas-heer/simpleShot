# simpleShot
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