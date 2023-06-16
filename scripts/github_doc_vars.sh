URL_RELEASES='https://api.github.com/repos/dhuan/mock/releases/latest'
TMP=$(mktemp)
TMP_DIR=$(mktemp -d)
VAR_DOWNLOAD_LINK_LINUX=$(curl -s "$URL_RELEASES" | grep browser_download_url | cut -d '"' -f 4 | grep 386)
VAR_DOWNLOAD_LINK_MACOS=$(curl -s "$URL_RELEASES" | grep browser_download_url | cut -d '"' -f 4 | grep darwin)

cat /dev/stdin > $TMP

cat $TMP | tar zx -C "$TMP_DIR"

cd "$TMP_DIR"

find . -type f | xargs sed -i 's|VAR_DOWNLOAD_LINK_LINUX|'"$VAR_DOWNLOAD_LINK_LINUX"'|g'
find . -type f | xargs sed -i 's|VAR_DOWNLOAD_LINK_MACOS|'"$VAR_DOWNLOAD_LINK_MACOS"'|g'

tar cz ./*
