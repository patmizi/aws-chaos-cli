package chaosutil

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  . "github.com/patmizi/aws-chaos-cli/lib/types"
  "github.com/pkg/errors"
)

func GetNaclsToChaos(client *ec2.EC2, subnetsToChaos []string) ([]NaclAssociationPair, error) {
  fmt.Print("Getting the list of NACLs to blackhole\n")

  aclClientResponse, err := client.DescribeNetworkAcls(&ec2.DescribeNetworkAclsInput{
    Filters: []*ec2.Filter {
      {
        Name: aws.String("association.subnet-id"),
        Values: aws.StringSlice(subnetsToChaos),
      },
    },
  })
  if err != nil {
    return make([]NaclAssociationPair, 0), errors.Wrap(err, "Failed to get a list of NACLs to chaos")
  }
  networkAcls := aclClientResponse.NetworkAcls

  var naclPairs []NaclAssociationPair
  for _, nacl := range networkAcls {
    for _, naclAssociation := range nacl.Associations {
      for _, el := range subnetsToChaos {
        if *naclAssociation.SubnetId == el {
          naclPairs = append(naclPairs, NaclAssociationPair{
            NaclAssociationId: *naclAssociation.NetworkAclAssociationId,
            NaclId:            *naclAssociation.NetworkAclId,
          })
        }
      }
    }
  }

  return naclPairs, nil
}