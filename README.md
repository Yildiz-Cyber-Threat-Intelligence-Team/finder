# Finder Tool

Finder Tool is a tool that extracts data from a compressed file and searches for user-specified text in the files inside.

**Flags**

    -file: Specify the compressed file path - required
    -text: Specify the text to search for or specify multiple texts separated by (,) - required
    -output: Specify the path to the output file to save the results - optional
    -case-sensitive: Specify whether the search should be case-sensitive (default: false) - optional
    -help: Help for using the finder tool

**Usage**

    go run main.go -file myfile.zip -text lorem,ipsum -output results.txt -case-sensitive true
