#!/bin/bash

ROOTDIR=$(pwd)
mkdir autoposting-release
make build

mv autoposting-bot autoposting-release
cd database/cmd


for i in autoposting-*; do 
	cd $i
	go build
	mv $i ../../../autoposting-release
	cd ../
done

cd $ROOTDIR
cd fingerprinting/cmd/hash-database
go build
mv hash-database ../../../autoposting-release
