#! /bin/bash

rm -R admin captain superadmin
ln -s $ASSPATH/webapp/dist/superadmin superadmin
ln -s $ASSPATH/webapp/dist/admin admin
ln -s $ASSPATH/webapp/dist/captain captain