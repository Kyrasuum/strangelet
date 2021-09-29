cd ..
NAME=strangelet

mkdir -p binaries
mkdir -p NAME-$1

cp LICENSE $NAME-$1
cp README.md $NAME-$1
cp LICENSE-THIRD-PARTY $NAME-$1
cp assets/packaging/$NAME.1 $NAME-$1
cp assets/packaging/$NAME.desktop $NAME-$1
cp assets/$NAME-logo-mark.svg $NAME-$1/$NAME.svg

HASH="$(git rev-parse --short HEAD)"
VERSION="$(go run tools/build-version.go)"
DATE="$(go run tools/build-date.go)"
ADDITIONAL_GO_LINKER_FLAGS="$(go run tools/info-plist.go $VERSION)"

# Mac
echo "OSX 64"
GOOS=darwin GOARCH=amd64 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-osx.tar.gz $NAME-$1
mv $NAME-$1-osx.tar.gz binaries

# Linux
echo "Linux 64"
GOOS=linux GOARCH=amd64 make build
./tools/package-deb.sh $1
mv $NAME-$1-amd64.deb binaries

mv $NAME $NAME-$1
tar -czf $NAME-$1-linux64.tar.gz $NAME-$1
mv $NAME-$1-linux64.tar.gz binaries

echo "Linux 64 fully static"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-linux64-static.tar.gz $NAME-$1
mv $NAME-$1-linux64-static.tar.gz binaries

echo "Linux 32"
GOOS=linux GOARCH=386 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-linux32.tar.gz $NAME-$1
mv $NAME-$1-linux32.tar.gz binaries

echo "Linux ARM 32"
GOOS=linux GOARCH=arm make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-linux-arm.tar.gz $NAME-$1
mv $NAME-$1-linux-arm.tar.gz binaries

echo "Linux ARM 64"
GOOS=linux GOARCH=arm64 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-linux-arm64.tar.gz $NAME-$1
mv $NAME-$1-linux-arm64.tar.gz binaries

# NetBSD
echo "NetBSD 64"
GOOS=netbsd GOARCH=amd64 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-netbsd64.tar.gz $NAME-$1
mv $NAME-$1-netbsd64.tar.gz binaries

echo "NetBSD 32"
GOOS=netbsd GOARCH=386 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-netbsd32.tar.gz $NAME-$1
mv $NAME-$1-netbsd32.tar.gz binaries

# OpenBSD
echo "OpenBSD 64"
GOOS=openbsd GOARCH=amd64 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-openbsd64.tar.gz $NAME-$1
mv $NAME-$1-openbsd64.tar.gz binaries

echo "OpenBSD 32"
GOOS=openbsd GOARCH=386 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-openbsd32.tar.gz $NAME-$1
mv $NAME-$1-openbsd32.tar.gz binaries

# FreeBSD
echo "FreeBSD 64"
GOOS=freebsd GOARCH=amd64 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-freebsd64.tar.gz $NAME-$1
mv $NAME-$1-freebsd64.tar.gz binaries

echo "FreeBSD 32"
GOOS=freebsd GOARCH=386 make build
mv $NAME $NAME-$1
tar -czf $NAME-$1-freebsd32.tar.gz $NAME-$1
mv $NAME-$1-freebsd32.tar.gz binaries

rm $NAME-$1/$NAME

# Windows
echo "Windows 64"
GOOS=windows GOARCH=amd64 make build
mv $NAME.exe $NAME-$1
zip -r -q -T $NAME-$1-win64.zip $NAME-$1
mv $NAME-$1-win64.zip binaries

echo "Windows 32"
GOOS=windows GOARCH=386 make build
mv $NAME.exe $NAME-$1
zip -r -q -T $NAME-$1-win32.zip $NAME-$1
mv $NAME-$1-win32.zip binaries

rm -rf $NAME-$1
