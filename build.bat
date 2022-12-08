SET GOOS=windows
go build -o ./bin/rank.exe .\rank\rankExec.go
go build -o ./bin/leaf.exe .\leaf.\leafserver\src\server\leafserver.go
SET GOOS=linux
go build -o ./bin/bin/rank .\rank\rankExec.go
go build -o ./bin/bin/leaf .\leaf.\leafserver\src\server\leafserver.go
