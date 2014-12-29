if "%1" == "clean" (
        for /R . %%I in (*~) do del %%I
) else if "%1" == "install" (
        if not "%2" == "" for /R . %%I in (*.exe) do copy "%%I" "%2"\.
) else (
        for %%I in (ansistrip camelfmt cure hexdump) do (pushd %%I & go build & popd)
)
