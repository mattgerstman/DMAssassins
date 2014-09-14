#!/bin/bash
for i in "$@"
do
     mv $i $i.old
     sed 's;response.response;response;g;' $i.old > $i
     rm -f $i.old
done     