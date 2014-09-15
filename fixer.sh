#!/bin/bash
for i in "$@"
do
     mv $i $i.old
     sed 's;users/;user/;g;' $i.old > $i
     rm -f $i.old
done     