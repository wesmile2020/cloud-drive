
echo build web-site

cd ./web-site
npm run build
cd ../

echo build server-site

cd ./server-site
export CGO_ENABLED=1
go build -o ../output/cloud-drive.exe main.go
cd ../

echo build success

