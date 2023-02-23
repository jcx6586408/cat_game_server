@echo off
start /min "center"  .\bin\center.exe ^
./conf/leafserver.json ^
./conf/server.json ^
./conf/room.json ^
./table ^
":5000" ":5600"

