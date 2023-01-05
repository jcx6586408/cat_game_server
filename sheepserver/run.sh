#!/bin/bash
nohup ./bin/leaf > leaf.log 2>&1 &
nohup ./bin/rank > rank.log 2>&1 &
nohup ./bin/center > center.log 2>&1 &