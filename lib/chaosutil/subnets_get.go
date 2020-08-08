package chaosutil

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/pkg/errors"
)

func GetSubnetsToChaos(client *ec2.EC2, vpcId string, azName string) ([]string, error) {
  fmt.Printf("Getting the list of subnets to fail in vpc: %s\n", vpcId)

  subnetList, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
    Filters: []*ec2.Filter{
      {
        Name:   aws.String("availability-zone"),
        Values: aws.StringSlice([]string{azName}),
      },
      {
        Name:   aws.String("vpc-id"),
        Values: aws.StringSlice([]string{vpcId}),
      },
    },
  })
  if err != nil {
    return make([]string, 0), errors.Wrap(err, "Failed to get a list of subnets")
  }

  var subnetsToChaos = make([]string, len(subnetList.Subnets))
  for i := 0; i < len(subnetList.Subnets); i++ {
    subnetsToChaos[i] = *subnetList.Subnets[i].SubnetId
  }

  return subnetsToChaos, nil
}
