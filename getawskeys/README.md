# get_old_aws_access_keys
This tool is designed to find and print AWS IAM Access Keys from an AWS Account. 
For example if your security policy indicates that you should rotate keys older than 5 months then this tool will find any keys older than 5 months and  output them to stdout. To specify the "months older than" period, you simply pass a parameter and a value combination called monthsolderthan to this tool.
For example:

```
./getawskeys -monthsolderthan 5
```
Will find and return keys older than 5 months. 

##Requires the following IAM permissions
The running user of this tool will need the following AWS permissions.
```
iam:ListUsers
iam:ListAccessKeys
```
