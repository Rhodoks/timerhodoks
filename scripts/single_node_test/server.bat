rmdir /s /q .\data
mkdir data
go build ../../cmd/node/node.go
.\node.exe --id=1