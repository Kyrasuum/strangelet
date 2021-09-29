NAME=strangelet
fpm -s dir -t deb -p $NAME-$1-amd64.deb --name $NAME --license mit --version $1 --deb-recommends xclip --description "A modern and intuitive terminal-based text editor" ./$NAME=/usr/bin/$NAME ./assets/packaging/$NAME.1=/usr/share/man/man1/$NAME.1
