#!/bin/bash
ROOTDIR=$(pwd)
VERSION=$(grep "VERSION :=" Makefile | awk '{print $3}')
DEST=autoposting-bot-v$VERSION
mkdir $DEST

echo "[!!] version $VERSION"
echo "[+] building autoposting-bot..."
make build

mv autoposting-bot $DEST
cd database/cmd

for i in autoposting-*; do 
	cd $i
	echo "[+] building $i..."
	go build
	mv $i ../../../$DEST
	cd ../
done
cd $ROOTDIR
cd cmd/hash-database
echo "[+] building hash-database..."
go build
mv hash-database ../../$DEST

cd $ROOTDIR
echo "[+] building tar.xz..."
tar -cf - $DEST | xz -9 -c - > $DEST.tar.xz
echo "[+] cleaning..."
rm -r $DEST
echo "[+] done!"
