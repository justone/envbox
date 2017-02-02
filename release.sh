#!/bin/bash

set -ex

if [[ ! $(type -P gox) ]]; then
    echo "Error: gox not found."
    echo "To fix: run 'go get github.com/mitchellh/gox', and/or add \$GOPATH/bin to \$PATH"
    exit 1
fi

if [[ ! $(type -P github-release) ]]; then
    echo "Error: github-release not found."
    exit 1
fi

VER=$1

if [[ -z $VER ]]; then
    echo "Need to specify version."
    exit 1
fi

PRE_ARG=
if [[ $VER =~ pre ]]; then
    PRE_ARG="--pre-release"
fi

git tag $VER

echo "Building $VER"
echo

gox -ldflags "-X main.version=$VER" -osarch="darwin/amd64 linux/amd64 windows/amd64 linux/arm"

echo "* " > desc
echo "" >> desc

echo "$ sha1sum envbox_*" >> desc
sha1sum envbox_* >> desc
echo "$ sha256sum envbox_*" >> desc
sha256sum envbox_* >> desc
echo "$ md5sum envbox_*" >> desc
md5sum envbox_* >> desc

vi desc

git push --tags

sleep 2

github-release release $PRE_ARG --user justone --repo envbox --tag $VER --name $VER --description desc
github-release upload --user justone --repo envbox --tag $VER --name envbox_darwin_amd64 --file envbox_darwin_amd64
github-release upload --user justone --repo envbox --tag $VER --name envbox_linux_amd64 --file envbox_linux_amd64
github-release upload --user justone --repo envbox --tag $VER --name envbox_linux_arm --file envbox_linux_arm
github-release upload --user justone --repo envbox --tag $VER --name envbox_windows_amd64.exe --file envbox_windows_amd64.exe
