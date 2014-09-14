#!/bin/bash
for i in "$@"
do
     mv $i $i.old
     sed 's;navView;NavView;g;' $i.old > $i
     rm -f $i.old
done     