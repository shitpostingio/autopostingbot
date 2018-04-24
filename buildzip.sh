#!/bin/bash

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
