#!/bin/bash
docker run -itd --restart=unless-stopped -p 5200:5200 --name sheepRank  sheepserver:v1  /bin/sheep/runRank.sh