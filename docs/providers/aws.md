# AWS Route 53 Provider

## Requirements

* AWS Access Key
* Setup IAM policy

## Get AWS Access Key

!!! warning "Warning"

    It's recommended to create dedicated user

Select your user in the [AWS console](https://console.aws.amazon.com/iam/home?#security_credential
) and click on the **Access Keys** tab.
Create a new access key and copy the Access Key ID and Secret Access Key.

## Setup IAM policy

Associate the IAM policy with the user.

```json title="IAM POLICY"
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "route53:GetChange",
                "route53:GetHostedZone",
                "route53:ListResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/*",
                "arn:aws:route53:::change/*",
                "arn:aws:route53:::trafficpolicy/*"
            ]
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZones",
                "route53:GetHostedZoneCount",
                "route53:ListHostedZonesByName",
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": "*"
        }
    ]
}
```
