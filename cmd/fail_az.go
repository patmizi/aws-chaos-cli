package cmd

import (
  "fmt"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/autoscaling"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/golang/glog"
  . "github.com/patmizi/aws-chaos-cli/lib/chaosutil"
  "github.com/spf13/cobra"
  "os"
  "time"
)

type failAzOptions struct {
  region              string
  vpcId               string
  azName              string
  duration            int
  limitAsg            bool
  failoverRds         bool
  failoverElasticache bool
  profile             string
}


func failAzCmd() *cobra.Command {
  o := &failAzOptions{}

  cmd := &cobra.Command{
    Use:   "fail-az",
    Short: "Simulates AZ failure with a Chaos NACL",
    Long:  "Simulate AZ failure: associate subnet(s) with a Chaos NACL that deny ALL Ingress and Egress traffic - blackhole",
    Run:   func(cmd *cobra.Command, args []string) {
      err := failAz(o.region, o.vpcId, o.azName, o.duration, o.limitAsg, o.failoverRds, o.failoverElasticache, o.profile)
      if err != nil {
        glog.Fatal("Fail-AZ Failed: %s\n", err)
        os.Exit(1)
      }
    },
  }

  cmd.Flags().StringVar(&o.region, "region", "", "The AWS region of choice")
  cmd.Flags().StringVar(&o.vpcId, "vpc-id", "", "The VPC ID of choice")
  cmd.Flags().StringVar(&o.azName, "az-name", "", "The name of the availability zone to blackout")
  cmd.Flags().IntVar(&o.duration, "duration", 0, "The duration, in seconds, of the blackout")
  cmd.Flags().BoolVar(&o.limitAsg, "limit-asg", false, "Remove 'failed' AZ from Auto Scaling Group (ASG)")
  cmd.Flags().BoolVar(&o.failoverRds, "rds", false, "Failover RDS if master in the blackout subnet")
  cmd.Flags().BoolVar(&o.failoverElasticache, "elasticache", false, "Failover Elasticache if primary in the blackout subnet")
  cmd.Flags().StringVar(&o.profile, "profile", "default", "AWS credential profile to use")

  cmd.MarkFlagRequired("region")
  cmd.MarkFlagRequired("vpc-id")
  cmd.MarkFlagRequired("az-name")

  return cmd
}


func failAz(region string, vpcId string, azName string, duration int, limitAsg bool, failoverRds bool, failoverElasticache bool, profile string) error {
	print(region)
	fmt.Printf("Setting up ec2 client for region %s", region)

	sess := session.Must(
		session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewSharedCredentials("", profile),
		}),
	)

	ec2Client := ec2.New(sess)
	autoscalingClient := autoscaling.New(sess)

	chaosNaclID, err := CreateChaosNacl(ec2Client, vpcId)
	if err != nil {
	  return err
  }

	subnetsToChaos, err := GetSubnetsToChaos(ec2Client, vpcId, azName)
	if err != nil {
	  return err
  }

	naclIds, err := GetNaclsToChaos(ec2Client, subnetsToChaos)
	if err != nil {
	  return err
  }

	var originalAsg *autoscaling.Group
	if limitAsg {
    originalAsg, err = LimitAutoScaling(autoscalingClient, subnetsToChaos)
    if err != nil {
      return err
    }
  }

  rollbackData, err:= ApplyChaosConfig(ec2Client, naclIds, chaosNaclID)
  if err != nil {
    return err
  }

  if failoverRds {
  }

  if failoverElasticache {

  }

  if duration > 0 {
    time.Sleep(time.Duration(duration) * time.Second)
  } else {
    _, err := fmt.Scanln()
    if err != nil {
      return err
    }
  }

  err = Rollback(ec2Client, rollbackData, autoscalingClient, originalAsg)
  if err != nil {
    return err
  }

  err = DeleteChaosNacl(ec2Client, chaosNaclID)
  if err != nil {
    return err
  }

  return nil
}