package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	svc := ec2.New(session.New())
	var maxResults int64
	maxResults = 500
	var nextToken string
	shouldDelete := true

	for {
		descSGInput := &ec2.DescribeSecurityGroupsInput{
			MaxResults: aws.Int64(maxResults),
		}

		if nextToken != "" {
			descSGInput.NextToken = aws.String(nextToken)
		}

		descSGOut, err := svc.DescribeSecurityGroups(descSGInput)
		if err != nil {
			log.Println("Unable to get security groups. Err:", err)
			os.Exit(1)
		}

		for _, sg := range descSGOut.SecurityGroups {
			niInput := &ec2.DescribeNetworkInterfacesInput{
				Filters: []*ec2.Filter{&ec2.Filter{
					Name:   aws.String("group-id"),
					Values: []*string{sg.GroupId},
				}},
			}

			niOut, err := svc.DescribeNetworkInterfaces(niInput)
			if err != nil {
				log.Println("Unable to get network interfaces for the security group. Err:", err)
				os.Exit(1)
			}

			if len(niOut.NetworkInterfaces) == 0 && shouldDelete {
				log.Println("Deleting Security Group: ", *sg.GroupId, *sg.GroupName)
				svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
					GroupId: sg.GroupId,
				})
			}
		}

		if descSGOut.NextToken != nil {
			nextToken = *descSGOut.NextToken
		} else {
			break
		}
	}
}
