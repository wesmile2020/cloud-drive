@echo off

echo build web-site

cd ./web-site
call npm run build
if %errorlevel% neq 0 (
  echo npm run build failed
  exit /b %errorlevel%
)
echo build web-site success

cd ../

echo build server-site

cd ./server-site
set CGO_ENABLED=1
go build -o ../output/cloud-drive.exe main.go
cd ../

echo build success
