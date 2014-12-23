build:
	go bulid

package:
	zip -9 camelfmt-%DATE:/=%.zip camelfmt.exe camelfmt.go vbsfmt.cmd

sweep:
	for %%I in (*~) do del %%I

clean:
	if exist camelfmt.exe del camelfmt.exe
