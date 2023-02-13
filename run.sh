#!/bin/bash
nohup ./bin/leaf \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
./ssl/Nginx/1_kampfiregames.cn_bundle.crt \
./ssl/Nginx/2_kampfiregames.cn.key \
> leaf.log 2>&1 &
nohup ./bin/rank \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
./ssl/Nginx/1_kampfiregames.cn_bundle.crt \
./ssl/Nginx/2_kampfiregames.cn.key \
> rank.log 2>&1 &
nohup ./bin/center \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
":5600" \
./ssl/Nginx/1_kampfiregames.cn_bundle.crt \
./ssl/Nginx/2_kampfiregames.cn.key \
> center.log 2>&1 &
nohup ./bin/center \
./conf/leafserver.json \
./conf/server.json \
./conf/room.json \
./table \
./IP2LOCATION-LITE-DB3.IPV6.BIN/IP2LOCATION-LITE-DB3.IPV6.BIN \
":5601" \
./ssl/Nginx/1_kampfiregames.cn_bundle.crt \
./ssl/Nginx/2_kampfiregames.cn.key \
> center.log 2>&1 &