#!/bin/bash

if [ -z $COVERALLS_TOKEN ]
then
    exit 1
fi

echo "mode: count" > fullcov.out
dirs=$(find ./* -maxdepth 10 -type d )
dirs=". $dirs"
for dir in $dirs;
do
        if ls "$dir"/*.go &> /dev/null;
        then
            go test -coverprofile=profile.out -covermode=count "$dir" -tags nolibnfc
            if [ -f profile.out ]
            then
                cat profile.out | grep -v "^mode: count" >> fullcov.out
            fi
        fi
done
$HOME/gopath/bin/goveralls -coverprofile=fullcov.out -service=travis-ci -repotoken $COVERALLS_TOKEN
rm -rf ./profile.out
rm -rf ./fullcov.out
