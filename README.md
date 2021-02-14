# s3 record filter by time

If the specified time is exceeded, the recrod will be ignored.

## Deploy

### env

```
export BUCKET=""
```

### command

```
make build
make package
make deploy
```

### test

```
export SOURCE_BUCKET=$(aws cloudformation describe-stacks --stack-name s3-record-filter --query 'Stacks[].Outputs[?OutputKey == `TestBucketARN`].OutputValue' --output text)
aws s3 cp <file_name> s3://$SOURCE_FILE
```
