Overview
--------

This tool removes unnecessary blocks of:

	ExplicitLeft = 0
	ExplicitTop = 0
	ExplicitWidth = 0
	ExplicitHeight = 0

from Delphi DFM files.

The Delphi IDE (Rad Studio) sometimes inserts these when you edit a Form visually.

These Explicit properties serve a real purpose sometimes. But if they are all zero, they are unnecessary.
Thus this program will leave them alone if any of them either does not exist or is non-zero. Only if
they are all zero, they are removed.

This will keep your diffs from being cluttered all the time.


Installation
------------

Install the Go programming language from https://golang.org/

and then run:

	go get -u github.com/gonutz/dfm_clear_explicit_zeros

where -u is for getting the latest online version.

Call

	dfm_clear_explicit_zeros file1.dfm file2.dfm ...

which will clear all the .dfm files from unnecessary Explicit properties.
