cd ../..
go build cmd/node/node.go
rm -rf ./data
cp ./scripts/pressure_test/config.json ./
cp ./scripts/pressure_test/Procfile ./
goreman start
