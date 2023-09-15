GCC=x86_64-w64-mingw32-gcc
GOOS=windows GOARCH=amd64 go get golang.org/x/sys/windows 
CC=$GCC CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build  -o tugboat.exe .


./bin/f2s -raw -compress -input ./plugin/comp_health/module.dll -output ./plugin/comp_health/module.hexbin
