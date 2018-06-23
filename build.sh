#!/bin/sh

pushd $(dirname "${0}") > /dev/null
HACKED_BASE=$(pwd -L)
FOLDER_NAME=inkyblackness-hacked

echo "Cleaning output directories..."
rm -rf _build
rm main-res.syso

mkdir -p $HACKED_BASE/_build/linux/$FOLDER_NAME
mkdir -p $HACKED_BASE/_build/win/$FOLDER_NAME


echo "Determining version"


MAJOR=`date +%Y`
MINOR=`date +%m`
PATCH=`date +%d`
VERSION=$(git describe exact-match --abbrev=0)
if [ $? -ne 0 ]; then
   echo "Not a tagged build, defaulting to revision for version"
   REV=$(git rev-parse --short HEAD)
   VERSION="rev$REV"
else
   VERSION_RAW=$(echo "$VERSION" | cut -d'v' -f 2)
   MAJOR=$(echo "$VERSION_RAW" | cut -d'.' -f 1)
   MINOR=$(echo "$VERSION_RAW" | cut -d'.' -f 2)
   PATCH=$(echo "$VERSION_RAW" | cut -d'.' -f 3)
fi
echo "Determined version: $VERSION"


echo "Preparing build resources"
mkdir -p $HACKED_BASE/_build/win_temp
cp $HACKED_BASE/_resources/build/win/* $HACKED_BASE/_build/win_temp
sed -i "s/§MAJOR/$MAJOR/g" $HACKED_BASE/_build/win_temp/hacked.exe.manifest
sed -i "s/§MINOR/$MINOR/g" $HACKED_BASE/_build/win_temp/hacked.exe.manifest
sed -i "s/§PATCH/$PATCH/g" $HACKED_BASE/_build/win_temp/hacked.exe.manifest
x86_64-w64-mingw32-windres -o main-res.syso $HACKED_BASE/_build/win_temp/hacked.rc

echo "Building executables..."
go build -ldflags "-X main.version=$VERSION" -a -o $HACKED_BASE/_build/linux/$FOLDER_NAME/hacked .
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -ldflags "-X main.version=$VERSION -H=windowsgui" -a -o $HACKED_BASE/_build/win/$FOLDER_NAME/hacked.exe .


echo "Copying distribution resources..."

for os in "linux" "win"
do
   packageDir=$HACKED_BASE/_build/$os/$FOLDER_NAME

   cp $HACKED_BASE/LICENSE $packageDir
   cp -R $HACKED_BASE/_resources/dist/* $packageDir
done

MINGW_BASE=/usr/x86_64-w64-mingw32/bin
for lib in "libgcc_s_seh-1.dll" "libstdc++-6.dll" "libwinpthread-1.dll"
do
   cp $MINGW_BASE/$lib $HACKED_BASE/_build/win/$FOLDER_NAME
done


echo "Creating packages..."

cd $HACKED_BASE/_build/linux
tar -cvzf $HACKED_BASE/_build/$FOLDER_NAME-$VERSION.linux64.tgz ./$FOLDER_NAME

cd $HACKED_BASE/_build/win
zip -r $HACKED_BASE/_build/$FOLDER_NAME-$VERSION.win64.zip .

popd > /dev/null