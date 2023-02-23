@echo off
start /min "rank" .\bin\rank.exe ^
./conf/leafserver.json ^
./conf/server.json ^
./conf/room.json ^
./table ^
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN 
start /min "center" .\bin\center.exe ^
./conf/leafserver.json ^
./conf/server.json ^
./conf/room.json ^
./table ^
":5000" ":5600"
start /min "home" .\bin\home.exe ^
./conf/leafserver.json ^
./conf/server.json ^
./conf/room.json ^
./table ^
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN ^
":5601"
