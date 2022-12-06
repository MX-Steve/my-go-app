cd "%~dp0"
call make_version.bat ./ version.h

cd "%~dp0"
windres.exe -i app.rc -o app.syso

call go build -ldflags "-H windowsgui" -o app.exe
