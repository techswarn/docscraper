#!/bin/sh

set -e

env >> /etc/environment

# execute CMD
echo "$@"
exec "$@"