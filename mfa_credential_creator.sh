# #This file (~/.aws/aws-mfa.file) should exist and contain environment variables AWS_MFA_ROLE_ARN, AWS_MFA_ARN, aws_access_key_id and aws_secret_access_key
# #If it doesn't then you need to provide them some other way i.e. exporting them as environment variables manually.

source ~/.aws/aws-mfa.file

role_arn="${AWS_MFA_ROLE_ARN:-}"
serial_number="${AWS_MFA_ARN:-}"
aws_access_key_id="${aws_access_key_id:-}"
aws_secret_access_key="${aws_secret_access_key:-}"

echo -n "Enter MFA Code: "
read -s token_code

temporary_credentials="$(aws \
                sts assume-role \
                --role-arn="${role_arn}" \
                --serial-number="${serial_number}" \
                --token-code="${token_code}" \
                --role-session-name="terraform-access" \
)"

unset AWS_PROFILE

export "AWS_ACCESS_KEY_ID=$(echo "${temporary_credentials}" | jq -re '.Credentials.AccessKeyId')"
export "AWS_SECRET_ACCESS_KEY=$(echo "${temporary_credentials}" | jq -re '.Credentials.SecretAccessKey')"
export "AWS_SESSION_TOKEN=$(echo "${temporary_credentials}" | jq -re '.Credentials.SessionToken')"
