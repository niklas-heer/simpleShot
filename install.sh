#!/usr/bin/env bash
set -e

[ -z "$SIMPLESHOT_DIR" ] && SIMPLESHOT_DIR="/usr/local/bin"
[ -z "$SIMPLESHOT_VERSION" ] && SIMPLESHOT_VERSION="master"
[ -z "$SIMPLESHOT_CONF_DIR" ] && SIMPLESHOT_CONF_DIR="$HOME"

SIMPLESHOT_SOURCE="https://raw.githubusercontent.com/niklas-heer/simpleShot/$SIMPLESHOT_VERSION/simpleShot"
SIMPLESHOT_CONF_SOURCE="https://raw.githubusercontent.com/niklas-heer/simpleShot/$SIMPLESHOT_VERSION/simpleShot-sample.gcfg"

echo "=> Downloading simpleShot to '$SIMPLESHOT_DIR'"
curl -sS "$SIMPLESHOT_SOURCE" -o "$SIMPLESHOT_DIR/simpleShot" || {
  echo >&2 "Failed to download '$SIMPLESHOT_SOURCE'"
  return 1
}

echo "=> Setting executable permissions on $SIMPLESHOT_DIR/simpleShot"
chmod +x "$SIMPLESHOT_DIR/simpleShot" || {
  echo >&2 "Failed setting executable permission on $SIMPLESHOT_DIR/simpleShot"
  return 1
}

echo "=> Downloading simpleShot config file to '$SIMPLESHOT_CONF_DIR'"
curl -sS "$SIMPLESHOT_CONF_SOURCE" -o "$SIMPLESHOT_CONF_DIR/.simpleShot.gcfg" || {
  echo >&2 "Failed to download '$SIMPLESHOT_CONF_SOURCE'"
  return 1
}

echo "simpleShot is ready to use! Enjoy! :)"