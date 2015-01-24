#!/usr/bin/env bash
set -e

[ -z "$SIMPLESHOT_DIR" ] && SIMPLESHOT_DIR="/usr/local/bin"
[ -z "$SIMPLESHOT_VERSION" ] && SIMPLESHOT_VERSION="master"

CHAG_SOURCE="https://raw.githubusercontent.com/niklas-heer/simpleShot/$SIMPLESHOT_VERSION/simpleShot"

echo "=> Downloading simpleShot to '$SIMPLESHOT_DIR'"
curl -sS "$CHAG_SOURCE" -o "$SIMPLESHOT_DIR/simpleShot" || {
  echo >&2 "Failed to download '$CHAG_SOURCE'"
  return 1
}

echo "=> Setting executable permissions on $SIMPLESHOT_DIR/simpleShot"
chmod +x "$SIMPLESHOT_DIR/simpleShot" || {
  echo >&2 "Failed setting executable permission on $SIMPLESHOT_DIR/simpleShot"
  return 1
}

echo "simpleShot is ready to use! Enjoy! :)"