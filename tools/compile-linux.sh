cd ..
NAME=strangelet

mkdir -p binaries
mkdir -p $NAME-$1

cp LICENSE $NAME-$1
cp README.md $NAME-$1
cp LICENSE-THIRD-PARTY $NAME-$1

HASH="$(git rev-parse --short HEAD)"
VERSION="$(go run tools/build-version.go)"
DATE="$(go run tools/build-date.go)"
ADDITIONAL_GO_LINKER_FLAGS="$(go run tools/info-plist.go $VERSION)"

echo "Linux 64"
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=$1 -X main.CommitHash=$HASH -X 'main.CompileDate=$DATE'" -o $NAME-$1/$NAME ./cmd/$NAME
tar -czf $NAME-$1-linux64.tar.gz $NAME-$1
mv $NAME-$1-linux64.tar.gz binaries
echo "Linux 32"
GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.Version=$1 -X main.CommitHash=$HASH -X 'main.CompileDate=$DATE'" -o $NAME-$1/$NAME ./cmd/$NAME
tar -czf $NAME-$1-linux32.tar.gz $NAME-$1
mv $NAME-$1-linux32.tar.gz binaries
echo "Linux arm 32"
GOOS=linux GOARCH=arm go build -ldflags "-s -w -X main.Version=$1 -X main.CommitHash=$HASH -X 'main.CompileDate=$DATE'" -o $NAME-$1/$NAME ./cmd/$NAME
tar -czf $NAME-$1-linux-arm.tar.gz $NAME-$1
mv $NAME-$1-linux-arm.tar.gz binaries
echo "Linux arm 64"
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.Version=$1 -X main.CommitHash=$HASH -X 'main.CompileDate=$DATE'" -o $NAME-$1/$NAME ./cmd/$NAME
tar -czf $NAME-$1-linux-arm64.tar.gz $NAME-$1
mv $NAME-$1-linux-arm64.tar.gz binaries

rm -rf $NAME-$1
