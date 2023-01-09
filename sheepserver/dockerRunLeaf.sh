#!/bin/bash
docker run -itd --restart=unless-stopped -p 5100:5100 --name sheep  sheepserver:v1  /bin/sheep/runLeaf.sh