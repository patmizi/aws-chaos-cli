package chaosutil

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/pkg/errors"
)

func CreateChaosNacl(client *ec2.EC2, vpcId string) (string, error) {
  fmt.Printf("Creating a Chaos Network ACL in %s\n", vpcId)

  chaosNacl, err := client.CreateNetworkAcl(&ec2.CreateNetworkAclInput{
    VpcId: &vpcId,
  })
  if err != nil {
    return "", errors.Wrap(err, "Failed to create network acl")
  }

  // Tag the network ACL with failover-testing-chaos-nacl for identification
  _, err = client.CreateTags(&ec2.CreateTagsInput{
    Resources: []*string{
      chaosNacl.NetworkAcl.NetworkAclId,
    },
    Tags: []*ec2.Tag{
      {
        Key:   aws.String("Name"),
        Value: aws.String("failover-testing-chaos-nacl"),
      },
    },
  })
  if err != nil {
    return "", errors.Wrap(err, "Failed to create nacl tags")
  }

  // Create egress entry
  _, err = client.CreateNetworkAclEntry(&ec2.CreateNetworkAclEntryInput{
    CidrBlock: aws.String("0.0.0.0/0"),
    Egress:    aws.Bool(true),
    PortRange: &ec2.PortRange{
      From: aws.Int64(0),
      To:   aws.Int64(65535),
    },
    NetworkAclId: chaosNacl.NetworkAcl.NetworkAclId,
    Protocol:     aws.String("-1"),
    RuleAction:   aws.String("deny"),
    RuleNumber:   aws.Int64(100),
  })
  if err != nil {
    return "", errors.Wrap(err, "Failed to create Egress nacl entry")
  }

  // Create ingress entry
  _, err = client.CreateNetworkAclEntry(&ec2.CreateNetworkAclEntryInput{
    CidrBlock: aws.String("0.0.0.0/0"),
    Egress:    aws.Bool(false),
    PortRange: &ec2.PortRange{
      From: aws.Int64(0),
      To:   aws.Int64(65535),
    },
    NetworkAclId: chaosNacl.NetworkAcl.NetworkAclId,
    Protocol:     aws.String("-1"),
    RuleAction:   aws.String("deny"),
    RuleNumber:   aws.Int64(101),
  })
  if err != nil {
    return "", errors.Wrap(err, "Failed to create Ingress nacl entry")
  }

  return *chaosNacl.NetworkAcl.NetworkAclId, nil
}