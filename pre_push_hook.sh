#!/usr/bin/env bash
app="csv2svg"
version="v0.0.1"

if ! type >/dev/null 2>&1 'trucolor'; then
  echo 'FYI: [MarkGriffiths/trucolor: 24bit color tools for the command line]( https://github.com/MarkGriffiths/trucolor )'
  echo "npm install --global trucolor"
  exit 1
fi

go test -cover | tee coverage.log
echo "[LOG] generate go test coverage.log"

coverage=$(cat coverage.log | grep coverage | grep -o -E '[0-9]+.[0-9]')
hsb_h=$(echo "scale=0;coverage" | bc)
rgb_hex_color=$(trucolor "hsb:$hsb_h,100,100")

echo "[LOG] generate coverage shields.io url"
shields_io_url="https://img.shields.io/badge/coverage-$coverage%25-$rgb_hex_color.svg"
echo "$shields_io_url"
wget -O coverage.svg "$shields_io_url"

echo "[LOG] generate version shields.io url"
shields_io_url="https://img.shields.io/badge/$app-$version-orange.svg"
echo "$shields_io_url"
wget -O version.svg "$shields_io_url"
