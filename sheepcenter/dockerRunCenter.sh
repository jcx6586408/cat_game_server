#!/bin/bash
docker run -itd --restart=unless-stopped -p 5600:5600 --name sheepCenter  sheepserver:v1  /bin/sheep/runCenter.sh