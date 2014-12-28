if "%1" == "clean" (
        for /R . %%I in (*~) do del %%I
) else (
        for %%I in (ansistrip camelfmt cure hexdump) do (pushd %%I & go build & popd)
)
