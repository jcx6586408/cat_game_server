#!/bin/bash
docker run -itd --restart=unless-stopped -p 5600:5600 --name sheepCenter  sheepserver:center  /bin/sheep/runCenter.sh