@setlocal
@set PROMPT=$G
if not "%1" == "" goto %1

:build
    go build
    goto end
:get
    go get github.com/mattn/go-runewidth
    go get github.com/shiena/ansicolor 
    go get github.com/zetamatta/nyagos/conio
    goto end
:fmt
    go fmt
    goto end
:clean
    if exist cure.exe del cure.exe
    goto end
:end
