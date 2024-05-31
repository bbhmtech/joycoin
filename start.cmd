@echo off
cd /d %~dp0
set CGO_ENABLED=1
set PATH=%PATH%;D:\msys64\mingw64\bin
go build -gcflags="all=-N -l" -o main.exe .\cmd\main
main.exe
pause
