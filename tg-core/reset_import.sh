rm -rf go.mod
rm -rf go.sum
rm -rf vendor
go mod init
go get -u -t ./...
sed -i "___bak" "s/v0.14.1/0.9.3/g" go.mod
rm -rf go.mod___bak
go mod download
go mod vendor
