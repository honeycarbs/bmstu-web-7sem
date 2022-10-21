#!/bin/sh

set -e

# Allow nginx to stay in the foreground
# so that Docker can track the process properly
nginx -g 'daemon off;'