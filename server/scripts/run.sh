#!/bin/sh

set -e

/etc/init.d/nginx start
supervisord -c /etc/supervisord.conf

/bin/bash
