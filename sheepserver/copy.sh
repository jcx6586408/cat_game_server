#!/bin/bash
tempPath="/bin/sheep/"
if [ ! -d "$tempPath" ]; then
mkdir /bin/sheep
fi
cp -r /home/* /bin/sheep/