set -ex

BUILD_BIN_PATH="./bin/mock"
APP_NAME="mock"
APP_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)
if [[ -z "$APP_VERSION" ]]
then
    APP_VERSION="dev"
fi

TARGETS=(
    "linux,386"
    "linux,amd64"
    "linux,arm"
    "linux,arm64"
    "darwin,amd64"
)

rm -rf ./release_downloads

mkdir ./release_downloads

for TARGET in "${TARGETS[@]}"
do
    GOOS=$(echo $TARGET | cut -d "," -f 1)
    GOARCH=$(echo $TARGET | cut -d "," -f 2)

    TARGET_NAME="${GOOS}-${GOARCH}"

    printf "Generating build for ${TARGET_NAME}\n"

    TARGET_PATH="./release_downloads/$TARGET_NAME"

    mkdir $TARGET_PATH

    cp ./README.md $TARGET_PATH/.
    cp ./LICENSE $TARGET_PATH/.

    TMP_BKP=$(mktemp)
    cp internal/cmd/version.go "$TMP_BKP"
    sed -i "s/__VERSION__/$APP_VERSION/g" internal/cmd/version.go
    sed -i "s/__GOOS__/$GOOS/g" internal/cmd/version.go
    sed -i "s/__GOARCH__/$GOARCH/g" internal/cmd/version.go

    GOOS=$GOOS GOARCH=$GOARCH make
    cp "$BUILD_BIN_PATH" "$TARGET_PATH"/.

    cp "$TMP_BKP" internal/cmd/version.go
done

TARGET_FOLDERS=$(ls ./release_downloads)

for TARGET_FOLDER in ${TARGET_FOLDERS[@]}
do
    tar czvf "./release_downloads/${APP_NAME}_${APP_VERSION}_${TARGET_FOLDER}.tar.gz" -C "./release_downloads/""${TARGET_FOLDER}" .
done
