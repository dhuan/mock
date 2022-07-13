set -ex

APP_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)

export GH_TOKEN=${GH_KEY}

echo $APP_VERSION

gh release create "$APP_VERSION" -t "$APP_VERSION"

RELEASE_FILES=$(ls ./release_downloads/*.tar.gz)

for RELEASE_FILE in $RELEASE_FILES
do
    gh release upload "$APP_VERSION" "$RELEASE_FILE"
done
