#!/bin/sh
mkdir -p run
echo "> Compiling"
go build
mv goelan run/
cd run/
echo "> Running"
echo
./goelan
cd ..
