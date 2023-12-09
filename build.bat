@echo off
go build --ldflags="-s -w" --trimpath -gcflags=all="-l -B -wb=false" -o gos.exe .\main.go