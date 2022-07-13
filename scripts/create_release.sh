set -ex

MOCK_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)

export GH_TOKEN=${GH_KEY}

echo $MOCK_VERSION

gh release create "$MOCK_VERSION" -t "$MOCK_VERSION"

RELEASE_FILES=$(ls ./release_downloads/*.zip)

for RELEASE_FILE in $RELEASE_FILES
do
    gh release upload "$MOCK_VERSION" "$RELEASE_FILE"
done
