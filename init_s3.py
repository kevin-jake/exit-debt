import boto3
import os
s3_client = boto3.client(
    "s3",
    endpoint_url=os.getenv("S3_ENDPOINT"),
    aws_access_key_id=os.getenv("S3_ACCESS_KEY_ID"),
    aws_secret_access_key=os.getenv("S3_SECRET_ACCESS_KEY")
)

s3_client.create_bucket(Bucket=str(os.getenv("S3_BUCKET_NAME")))

print(f"Bucket {os.getenv('S3_BUCKET_NAME')} created successfully")