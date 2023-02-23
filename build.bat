@echo off


SET GOOS=windows
start go build -o ./bin/home.exe .\rank\rankExec.go

start go build -o ./bin/center.exe .\center\main.go

start go build -o ./bin/rank.exe .\storage\main\rank.go

start go build -o ./bin/leaf.exe .\leaf.\leafserver\src\server\leafserver.go

SET GOOS=linux
start go build -o ./bin/bin/home .\rank\rankExec.go

start go build -o ./bin/bin/center .\center\main.go

start go build -o ./bin/bin/rank .\storage\main\rank.go

start go build -o ./bin/bin/leaf .\leaf.\leafserver\src\server\leafserver.go

