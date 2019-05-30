# csv2svg

This command generates svg table from csv format.

## how to install
```
go get -u github.com/umaumax/csv2svg/...
```

## how to use
```
csv2svg data/sample.csv
# output is data/sample.svg

cat data/sample.csv | csv2svg
# output is /dev/stdout
```

## TODO
* Write test code!
* stdin, file csv, markdown table format file
* output stdout or files
* table title
* header style only bold???  線を太くする??? style???

## NOTE
* Don't surround words by `"` in csv files.

## FYI
* [ajstarks/svgo: Go Language Library for SVG generation]( https://github.com/ajstarks/svgo )
