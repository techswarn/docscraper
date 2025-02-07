#!/bin/bash

set -e

figlet Welcome to Ubuntu playground
CMD cron && tail -f /var/log/cron.log