package cmd

import (
  "strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/autoscaling"
  "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/patmizi/aws-chaos-cli/lib"
	"github.com/google/logger"
	"github.com/spf13/cobra"
)

var (
	region              string
	vpcId               string
	azName              string
	duration            int
	limitAsg            bool
	failoverRds         bool
	failoverElasticache bool
	profile             string
)

type NaclAssociationPair struct {
  NaclAssociationId string
  NaclId string
}

func configureArgs() {
	failAzCmd.Flags().StringVar(&region, "", "", "The AWS region of choice")
	failAzCmd.Flags().StringVar(&vpcId, "", "", "The VPC ID of choice")
	failAzCmd.Flags().StringVar(&azName, "", "", "The name of the availability zone to blackout")
	failAzCmd.Flags().IntVar(&duration, "", 0, "The duration, in seconds, of the blackout")
	failAzCmd.Flags().BoolVar(&limitAsg, "", false, "Remove 'failed' AZ from Auto Scaling Group (ASG)")
	failAzCmd.Flags().BoolVar(&failoverRds, "", false, "Failover RDS if master in the blackout subnet")
	failAzCmd.Flags().BoolVar(&failoverElasticache, "", false, "Failover Elasticache if primary in the blackout subnet")
	failAzCmd.Flags().StringVar(&profile, "", "default", "AWS credential profile to use")
}

func configureLogging() {
	logger.Init("AZFailoverLogger", true, false, nil)
}

func init() {
	configureArgs()
	rootCmd.AddCommand(failAzCmd)
}

var failAzCmd = &cobra.Command{
	Use:   "fail-az",
	Short: "Simulates AZ failure with a Chaos NACL",
	Long:  "Simulate AZ failure: associate subnet(s) with a Chaos NACL that deny ALL Ingress and Egress traffic - blackhole",
	Run:   FailAz,
}

// Simulates AZ failure with a Chaos NACL
func FailAz(cmd *cobra.Command, args []string) {
	failAz(region, vpcId, azName, duration, limitAsg, failoverRds, failoverElasticache, profile)
}

func failAz(region string, vpcId string, azName string, duration int, limitAsg bool, failoverRds bool, failoverElasticache bool, profile string) {
	configureLogging()
	logger.Info("Setting up ec2 client for region %s", region)

	sess := session.Must(
		session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewSharedCredentials("", profile),
		}),
	)

	ec2Client := ec2.New(sess)
	autoscalingClient := autoscaling.New(sess)

	chaosNaclID := createChaosNacl(ec2Client, vpcId)
	subnetsToChaos := getSubnetsToChaos(ec2Client, vpcId, azName)
	naclIds := getNaclsToChaos(ec2Client, subnetsToChaos)

	var originalAsg string
	if limitAsg {
    originalAsg = limitAutoScaling(autoscalingClient, subnetsToChaos)
  }
}

func createChaosNacl(client *ec2.EC2, vpcId string) string {
	logger.Info("Creating a Chaos Network ACL in %s", vpcId)

	chaosNacl, err := client.CreateNetworkAcl(&ec2.CreateNetworkAclInput{
		VpcId: &vpcId,
	})
	if err != nil {
		logger.Error("Failed to create network acl: %v", err)
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
		logger.Error("Failed to create nacl tags: %v", err)
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
		logger.Error("Failed to create nacl entry: %v", err)
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
		logger.Error("Failed to create nacl entry: %v", err)
	}

	return *chaosNacl.NetworkAcl.NetworkAclId
}

func getSubnetsToChaos(client *ec2.EC2, vpcId string, azName string) []string {
	logger.Info("Getting the list of subnets to fail in vpc: %s", vpcId)

	subnetList, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("availability-zone"),
				Values: aws.StringSlice([]string{azName}),
			},
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{azName}),
			},
		},
	})
	if err != nil {
		logger.Error("Failed to get a list of subnets: %v", err)
	}

	var subnetsToChaos []string
	for i := 0; i < len(subnetList.Subnets); i++ {
		subnetsToChaos[i] = *subnetList.Subnets[i].SubnetId
	}

	return subnetsToChaos
}

func getNaclsToChaos(client *ec2.EC2, subnetsToChaos []string) []NaclAssociationPair {
  logger.Info("Getting the list of NACLs to blackhole")

  aclClientResponse, err := client.DescribeNetworkAcls(&ec2.DescribeNetworkAclsInput{
    Filters: []*ec2.Filter {
      {
        Name: aws.String("association.subnet-id"),
        Values: aws.StringSlice(subnetsToChaos),
      },
    },
  })
  if err != nil {
    logger.Error("Failed to get a list of NACLs to chaos: %v", err)
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

  return naclPairs
}

func limitAutoScaling(client *autoscaling.AutoScaling, subnetsToChaos []string) (asg string, err error) {
  logger.Info("Limit autoscaling to the remaining subnets")

  autoscalingResponse, err := client.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
    AutoScalingGroupNames: nil,
  })
  if err != nil {
    logger.Error("Failed to describe autoscaling groups: %v", err)
    return _, err
  }
  autoScalingGroups := autoscalingResponse.AutoScalingGroups

  // Find ADG that needs to be modified assuming only one ASG should be impacted
  correctAsg := false
  var subnetsToKeep []string
  var asgName string
  for _, asg := range autoScalingGroups {
    asgName = *asg.AutoScalingGroupName
    asgSubnets := strings.Split(*asg.VPCZoneIdentifier, ",")

    subnetsToKeep = lib.ListDiff(asgSubnets, subnetsToChaos)

    correctAsg := len(subnetsToKeep) < len(asgSubnets)
    if correctAsg { break }
  }

  if correctAsg {
    vpcZoneIdentifier := strings.Join(subnetsToKeep, ",")
    autoscalingResponse, err := client.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
      AutoScalingGroupName:             aws.String(asgName),
      VPCZoneIdentifier:                aws.String(vpcZoneIdentifier),
    })
    if err != nil {
      logger.Error("Unable to update ASG: %v", err)
    }
  }

}