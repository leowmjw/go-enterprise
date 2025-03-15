#!/bin/bash
##
## URL: https://medium.com/@betashorts1998/understanding-bash-pipelines-and-set-o-pipefail-ba7e06ffb684
## Alway set pipefail
# Now we’ll enable set -o pipefail and observe the behavior when the middle command fails.
# set -o pipefail
# Watch out for common Bash issues --> https://mywiki.wooledge.org/BashPitfalls#myprogram_2.3E.26-
# Always use strict mode .. unless u like debugging
# URL: http://redsymbol.net/articles/unofficial-bash-strict-mode/#expect-nonzero-exit-status
set -euo pipefail
IFS=$'\n\t'

# Be careful; in strict mode; use if else instead ..
echo "First Command" | unknown_command | echo "Second Command" | grep "Command" | awk "{print \$2}"

# Check there is OrderID passed in and it is int?
echo "Collected Payment for OrderID"

errCode=$?
if [ $errCode -ne 0 ]; then
  echo "If exit status is not 0, this line is printed."
  echo "ERR: $errCode"
  exit $errCode
else
  echo "Payment with ID completed!!!"
fi

