# aws-chaos-cli
CLI tool created with the sole purpose of injecting failure into AWS infrastructure in a convenient manner

## Disclaimer
⚠️USE AT YOUR OWN RISK⚠️

Using this tool may create an unreasonable risk. If you choose to use this tool in your own
activities, you do so at your own risk. None of the authors or contributors, or anyone else connected with the tool(s),
in any way whatsoever, can be responsible for your use of the tool(s) contained in this repository.
Use these tool(s) only if you understand what the code does.

## Usage
```shell script
CLI tool created with the sole purpose of injecting failure into AWS infrastructure in a convenient manner

Usage:
  aws-chaos-cli [command]

Available Commands:
  fail-az     Simulates AZ failure with a Chaos NACL
  help        Help about any command

Flags:
  -h, --help   help for aws-chaos-cli
```

fail-az
```shell script
Simulate AZ failure: associate subnet(s) with a Chaos NACL that deny ALL Ingress and Egress traffic - blackhole

Usage:
  aws-chaos-cli fail-az [flags]

Flags:
      --az-name string   The name of the availability zone to blackout
      --duration int     The duration, in seconds, of the blackout
  -h, --help             help for fail-az
      --limit-asg        Remove 'failed' AZ from Auto Scaling Group (ASG)
      --profile string   AWS credential profile to use (default "default")
      --region string    The AWS region of choice
      --vpc-id string    The VPC ID of choice
```

## Build from source
go version 1.14+ is required
```shell script
# Install dependencies
make deps

# Build for your platform
make build
```

## Support
Please use the github issue tracker