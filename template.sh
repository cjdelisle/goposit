#!/bin/bash
ls ./*.goh | sed 's/\.goh$//' | while read x; do
    gcc -E -P -x c ${x}.goh > ${x}_gen.go && \
    sed -i '' 's@__COMMENT__@//@g' ${x}_gen.go && \
    sed -i '' 's@__NEWLINE__@\
@g' ${x}_gen.go;
done
