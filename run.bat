@echo on
call kill.bat
start /min "leaf" .\bin\leaf.exe ^
./conf/leafserver.json ^
./conf/server.json ^
./conf/room.json ^
./table ^
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN 
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
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN 