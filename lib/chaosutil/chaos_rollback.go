package chaosutil

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/autoscaling"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/golang/glog"
  . "github.com/patmizi/aws-chaos-cli/lib/types"
  "github.com/pkg/errors"
)

func Rollback(ec2Client *ec2.EC2, rollbackData []NaclAssociationPair, asClient *autoscaling.AutoScaling, originalAsg *autoscaling.Group) error {
  glog.Info("Rolling back Network ACL to original configuration")

  for _, nacl := range rollbackData {
    _, err := ec2Client.ReplaceNetworkAclAssociation(&ec2.ReplaceNetworkAclAssociationInput{
      AssociationId: aws.String(nacl.NaclAssociationId),
      NetworkAclId:  aws.String(nacl.NaclId),
    })
    if err != nil {
      glog.Error("(Skipping) Could not replace network acl association: %v\n", err)
    } else {
      glog.Info("Rolled back: %s\n", nacl.NaclId)
    }
  }

  if originalAsg != nil {
    glog.Info("Rolling back AutoScalingGroup to original configuration\n")
    _, err := asClient.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
      AutoScalingGroupName:             originalAsg.AutoScalingGroupName,
      VPCZoneIdentifier:                originalAsg.VPCZoneIdentifier,
    })
    if err != nil {
      return errors.Wrap(err, "Could not update autoscaling group")
    }

    glog.Info("Reverted back autoscaling group changes\n")
  }

  return nil
}