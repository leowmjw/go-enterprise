#!/bin/bash
##
# Always use strict mode .. unless u like debugging
# URL: http://redsymbol.net/articles/unofficial-bash-strict-mode/#expect-nonzero-exit-status
set -euo pipefail
IFS=$'\n\t'

echo "First Command" | unknown_command | echo "Second Command" | grep "Command" | awk "{print \$2}"

# Open file of csv of OrderID to process ..
# Add loop here to call the important Payment script

if [ $? -ne 0 ]; then
  echo "If exit status is not 0, this line is printed."
  echo "ERR: $?"
else
  echo "Payment with ID completed!!!"
fi

