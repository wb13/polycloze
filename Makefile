.PHONY:	check
check:
	cd srs; go test -tags sqlite_math_functions
