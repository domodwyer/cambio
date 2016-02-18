#Cambio

A simple dynamic DNS service in-a-box.

##How to
Assuming you have an [Amazon AWS](https://aws.amazon.com/) account, a domain pointed at Route53, and a hosted zone configured:

1. Create a restricted IAM account (see below), unused by anything else and download the credentials.
2. Configure your [AWS shared credentials file](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html#the-shared-credentials-file)
3. Run `cambio -zone <hosted-zone-id> -domain <vpn.example.com.>` perodically via cron or the likes.

__Note__: If the record doesn't exist, it will be created automatically.

##Advanced
You can configure cambio to use different AWS profiles, set different DNS time-to-live values (default 5 minutes), or create/update different record types with the following arguements:
```
-profile string
    	AWS Credential profile name (default "default")
-record-type string
    	Record type (default "A")
-region string
    	Region (default "eu-west-1")
-ttl int
    	Time-to-live value (default 300)
```

###Seperate those privilages! (Restricted IAM Account)
Please only grant the bare minimum permissions to this IAM account, there's really no need to grant privileges to spin up thirty d2.8xlarge instances to something that changes a DNS record.

Anywhoo:

1. Create a new IAM user with a descriptive name like `home-dns-updater` - _use something you'll recognise in 6 months!_
2. Create an access key (on the `Security Credentials` tab), and save it for the next step.
3. Add the access key to the [AWS shared credentials file](https://docs.aws.amazon.com/cli/latest/topic/config-vars.html#the-shared-credentials-file) (defaults to _~/.aws/credentials_)
4. Grant restricted permissions to the IAM account:
	1. Under the `Permissions` tab, click `Inline Policies` to expand the container, and click create a new policy
	2. Select `Custom Policy`
	3. Give the policy a descriptive name, like `UpdateHomeDNSRecords`
	4. Paste the policy below, making sure you replace `<zone-id>` with your actual hosted zone ID
	```
	{
	    "Version": "2012-10-17",
	    "Statement": [
	        {
	            "Sid": "Stmt1452005095000",
	            "Effect": "Allow",
	            "Action": [
	                "route53:ChangeResourceRecordSets"
	            ],
	            "Resource": [
	                "arn:aws:route53:::hostedzone/<zone-id>"
	            ]
	        }
	    ]
	}
	```
	5. Click `Apply Policy`


Alternatively if you're on a EC2 instance you can use an [IAM role](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) for authentication (but if you're on EC2 you probably don't need this).