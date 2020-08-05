package chaosutil

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/golang/glog"
  . "github.com/patmizi/aws-chaos-cli/lib/types"
)

func ApplyChaosConfig(client *ec2.EC2, naclIds []NaclAssociationPair, chaosNaclId string) ([]NaclAssociationPair, error) {
  glog.Info("Saving original config & applying new chaos config\n")

  var rollbackData []NaclAssociationPair

  for _, nacl := range naclIds {
    response, err := client.ReplaceNetworkAclAssociation(&ec2.ReplaceNetworkAclAssociationInput{
      AssociationId: aws.String(nacl.NaclAssociationId),
      NetworkAclId:  aws.String(chaosNagolanclId),
    })
    if err != nil {
      glog.Errorf("(Skipping) Unable to replace network acl association for %s: %v\n", nacl.NaclAssociationId, err)
    } else {
      rollbackData = append(rollbackData, NaclAssociationPair{
        NaclAssociationId: *response.NewAssociationId,
        NaclId:            nacl.NaclId,
      })
    }
  }

  return rollbackData, nil
}
