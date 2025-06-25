
echo build web-site

cd ./web-site
npm run build
if [ $? -ne 0 ]; then
  echo npm run build failed
  exit 1
fi
echo build web-site success

cd ../

echo build server-site

cd ./server-site
export CGO_ENABLED=1
go build -o ../output/cloud-drive main.go
cd ../

echo build success

