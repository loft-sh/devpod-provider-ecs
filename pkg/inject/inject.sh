INSTALL_DIR="{{ .InstallDir }}"
INSTALL_FILENAME="{{ .InstallFilename }}"

INSTALL_PATH="$INSTALL_DIR/$INSTALL_FILENAME"

command_exists() {
  command -v "$@" >/dev/null 2>&1
}

is_arm() {
  case "$(uname -a)" in
  *arm* ) true;;
  *arm64* ) true;;
  *aarch* ) true;;
  *aarch64* ) true;;
  * ) false;;
  esac
}

download() {
  DOWNLOAD_URL="{{ .DownloadAmd }}"
  if is_arm; then
    DOWNLOAD_URL="{{ .DownloadArm }}"
  fi
  iteration=1
  max_iteration=3

  while :; do
    if [ "$iteration" -gt "$max_iteration" ]; then
      >&2 echo "error: failed to download devpod"
      exit 1
    fi

    cmd_status=""
    if command_exists curl; then
        curl -fsSL $DOWNLOAD_URL -o $INSTALL_PATH.$$ && break
        cmd_status=$?
    elif command_exists wget; then
        wget -q $DOWNLOAD_URL -O $INSTALL_PATH.$$ && break
        cmd_status=$?
    else
        echo "error: no download tool found, please install curl or wget"
        exit 127
    fi
    >&2 echo "error: failed to download devpod"
    >&2 echo "       command returned: ${cmd_status}"
    >&2 echo "Trying again in 10 seconds..."
    iteration=$((iteration+1))
    sleep 10
  done

  mv $INSTALL_PATH.$$ $INSTALL_PATH
  chmod +x $INSTALL_PATH
}

if [ "$($INSTALL_PATH --version 2>/dev/null || echo 'false')" = "false" ]; then
  mkdir -p $INSTALL_DIR || true
  rm -f $INSTALL_PATH 2>/dev/null || true

  download
fi

# Execute command
{{ .Command }}