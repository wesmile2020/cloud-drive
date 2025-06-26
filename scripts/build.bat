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
call :build_server_site windows amd64
call :build_server_site windows arm64
call :build_server_site linux amd64
call :build_server_site linux arm64
call :build_server_site darwin amd64
call :build_server_site darwin arm64

cd ../

echo build success
exit /b 0

:build_server_site
set GOOS=%1
set GOARCH=%2
if "%GOOS%" == "" (
  echo GOOS is empty, please input windows or linux or darwin
  exit /b 1
)
if "%GOARCH%" == "" (
  echo GOARCH is empty, please input amd64 or arm64 or 386
  exit /b 1
)

echo build %GOOS% %GOARCH% platform start

set assetName=cloud-drive-%GOOS%-%GOARCH%
if "%GOOS%" == "windows" (
  set assetName=%assetName%.exe
)

go build -o ../out/%assetName% main.go
if %errorlevel% neq 0 (
  echo go build %GOOS% %GOARCH% failed
  exit /b %errorlevel%
)
echo build %GOOS% %GOARCH% platform success
exit /b 0
