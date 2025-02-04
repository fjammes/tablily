#!/bin/bash

set -euxo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

datadir=$DIR/chase-the-devil

go build $DIR

tablily_tab=$($DIR/tablily --instrument bass --input $datadir/chase-the-devil.bass.txt)

sed "s/%TABLILY_TAB%/"$tablily_tab/" $DIR/chase-the-devil.ly.tpl > $DIR/chase-the-devil.ly