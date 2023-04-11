# Running Locally 
Go Version : 1.19.5
# Dependencies 
Get all dependencies with 
```go mod download```

# Setting up AWS credentials

Use any of below sources to set up AWS environment.
1. Environment Variables
2. Shared Configuration and Shared Credentials files.

i.e.
~/.aws/credentials
```
[default]
aws_access_key_id=TEST_ID
aws_secret_access_key=TEST_KEY
```

# Environment variables 
Add ROLE_ARN which has proper permissions. 

Here's simple wildcard permissions for EC2 resource you can use. 

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:*"
            ],
            "Resource": "*"
        }
    ]
}
```

# Build command 
``` go build provisioner.go ```

# Run command
Apply Terraform script in terraform/main.tf.
```
 ./provisioner apply
```
Deprovision changes made in last step.
```
 ./provisioner destroy
```
