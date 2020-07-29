package chaosutil

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/pkg/errors"
)

func DeleteChaosNacl(client *ec2.EC2, chaosNaclId string) error {
  fmt.Print("Deleting the Chaos NACL")

  _, err := client.DeleteNetworkAcl(&ec2.DeleteNetworkAclInput{
    NetworkAclId: aws.String(chaosNaclId),
  })

  if err != nil {
    return errors.Wrap(err, "Could not delete network ACL")
  }

  return nil
}
