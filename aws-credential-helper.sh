#!/bin/bash

# Might be good to use this line and somedefaultname
#export TOTP="$(keyfob otp ${AWS_MFA_NAME:-somedefaultname})"

export TOTP="$(keyfob otp ${AWS_MFA_NAME})"
if [[ -n "${TOTP:-}" ]]
then
  aws-vault exec --mfa-token=${TOTP} -j engineer
else
    echo "No MFA TOTP! 2fa did not find a MFA TOTP."
fi