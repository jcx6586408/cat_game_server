#!/bin/bash
nohup ./bin/leaf \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
./easygame2021.com_nginx/easygame2021.com_bundle.crt \
./easygame2021.com_nginx/easygame2021.com.key \
> leaf.log 2>&1 &
nohup ./bin/rank \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
./easygame2021.com_nginx/easygame2021.com_bundle.crt \
./easygame2021.com_nginx/easygame2021.com.key \
> rank.log 2>&1 &
nohup ./bin/center \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
./easygame2021.com_nginx/easygame2021.com_bundle.crt \
./easygame2021.com_nginx/easygame2021.com.key \
> center.log 2>&1 &