package chaosutil

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/autoscaling"
  "github.com/patmizi/aws-chaos-cli/lib/datautil"
  "github.com/pkg/errors"
  "strings"
)

func LimitAutoScaling(client *autoscaling.AutoScaling, subnetsToChaos []string) (*autoscaling.Group, error) {
  fmt.Print("Limit autoscaling to the remaining subnets\n")

  autoscalingResponse, err := client.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
    AutoScalingGroupNames: nil,
  })
  if err != nil {
    return nil, errors.Wrap(err, "Unable to describe autoscaling groups")
  }
  autoScalingGroups := autoscalingResponse.AutoScalingGroups

  // Find ASG that needs to be modified assuming only one ASG should be impacted
  var subnetsToKeep []string
  var asgName string
  for _, asg := range autoScalingGroups {
    asgName = *asg.AutoScalingGroupName
    asgSubnets := strings.Split(*asg.VPCZoneIdentifier, ",")

    subnetsToKeep = datautil.ListDiff(asgSubnets, subnetsToChaos)

    correctAsg := len(subnetsToKeep) < len(asgSubnets)
    if correctAsg {
      vpcZoneIdentifier := strings.Join(subnetsToKeep, ",")
      _, err := client.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
        AutoScalingGroupName:             aws.String(asgName),
        VPCZoneIdentifier:                aws.String(vpcZoneIdentifier),
      })
      if err != nil {
        return nil, errors.Wrapf(err, "Unable to update autoscaling group %s", asgName)
      }
      return asg, nil
    }
  }

  return nil, errors.New("Cannot find impacted ASG")
}
