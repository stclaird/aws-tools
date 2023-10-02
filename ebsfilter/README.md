# EBS Filter

## Summary
This tool will retrieve the EC2 instances in an AWS account. It will print a single line summary for each instance it retrieves. This will look like the following:

```
i-0b3d971517a9a0105, devops-elasticsearch-data, 13.40.66.2, 10.0.1.13, t2.small, running, 21.5
i-000c41dee200753ba, devops-elasticsearch-master, 3.8.211.61, 10.0.1.138, t2.small, running, 25.8
i-0669da019fdc013f4, devops-elasticsearch-data, 13.40.42.33, 10.0.1.59, t2.small, running, 17.2
i-07b8bf6bb48ef4315, devops-elasticsearch-data, 3.9.175.54, 10.0.1.148, t2.small, running, 17.2
```
The first field is the instanceID, Name Tag, Public IP address, Private IP address, Instance Type, Instance state and the sum of the attached volumes in GigaBytes (Rather than GibiBytes)

## Building the tool
This tool is written and tested in Go 1.19. Other versions of go should work as well.
- Build the binary
`go build -o ebsfilter .`
This will compile an output a binary called ebsfilter

## Running the tool

`./ebsfilter`

If all goes well this will output some instance data

```
i-0b3d971517a9a0105, devops-elasticsearch-data, 13.40.66.2, 10.0.1.13, t2.small, running, 21.5
i-000c41dee200753ba, devops-elasticsearch-master, 3.8.211.61, 10.0.1.138, t2.small, running, 25.8
i-0669da019fdc013f4, devops-elasticsearch-data, 13.40.42.33, 10.0.1.59, t2.small, running, 17.2
i-07b8bf6bb48ef4315, devops-elasticsearch-data, 3.9.175.54, 10.0.1.148, t2.small, running, 17.2
```

Alternatively you can add the instancename parameter to filter by the 

`./ebsfilter -instancename devops-elasticsearch-master`

This will produce something like the following:

```
i-000c41dee200753ba, devops-elasticsearch-master, 3.8.211.61, 10.0.1.138, t2.small, running, 25.8
```

## Troubleshooting

If you experience empty results, please check you are running against the correct AWS region.  
This can be fixed quickly by setting the correct region
`AWS_REGION=eu-west-1`

If you receive a message like the following: 
```
NoCredentialProviders: no valid providers in chain. Deprecated.
```
Check you have setup AWS authentication correctly