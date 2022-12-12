@echo off


SET GOOS=windows
start go build -o ./bin/rank.exe .\rank\rankExec.go

start go build -o ./bin/leaf.exe .\leaf.\leafserver\src\server\leafserver.go

SET GOOS=linux
start go build -o ./bin/bin/rank .\rank\rankExec.go

start go build -o ./bin/bin/leaf .\leaf.\leafserver\src\server\leafserver.go

