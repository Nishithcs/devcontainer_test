package aws

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func CreateARecord(subdomain string) error {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	r53 := route53.NewFromConfig(cfg)

	hostedZoneID := os.Getenv("HOSTED_ZONE_ID")
	domain := os.Getenv("REVERSE_PROXY_BASE_URL")
	ip := os.Getenv("REVERSE_PROXY_IP")

	recordName := fmt.Sprintf("%s.%s.", subdomain, domain)
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: types.RRTypeA,
						TTL:  aws.Int64(2147483647),
						ResourceRecords: []types.ResourceRecord{
							{Value: aws.String(ip)},
						},
					},
				},
			},
			Comment: aws.String(fmt.Sprintf("Updated by Go app at %v", time.Now())),
		},
	}

	response, err := r53.ChangeResourceRecordSets(ctx, input)
	fmt.Println(response)
	return err
}
