package main

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/rdegges/go-ipify"
)

type appConfig struct {
	Profile string
	Region  string
	Domain  string
	Zone    string
	TTL     int64
	Type    string
}

var Config = &appConfig{}

func init() {
	flag.StringVar(&Config.Profile, "profile", "default", "AWS Credential profile name")
	flag.StringVar(&Config.Region, "region", "eu-west-1", "Region")
	flag.StringVar(&Config.Type, "record-type", "A", "Record type")
	flag.Int64Var(&Config.TTL, "ttl", 300, "Time-to-live value")

	// Required
	flag.StringVar(&Config.Zone, "zone", "", "Hosted zone ID")
	flag.StringVar(&Config.Domain, "domain", "vpn.example.com", "FQDN to update")

	flag.Parse()
	if Config.Zone == "" || Config.Domain == "" {
		flag.PrintDefaults()
		log.Fatal("Error: Arguments \"-zone\" and \"-domain\" are required")
	}
}

func main() {
	ip, err := ipify.GetIp()
	if err != nil {
		log.Fatal("Couldn't get external IP address: ", err)
	}

	out, err := UpdateIp(ip)
	if err != nil {
		log.Fatal("Couldn't update DNS record: ", err)
	}

	log.Print(out)
}

func UpdateIp(newIp string) (*route53.ChangeResourceRecordSetsOutput, error) {
	svc := route53.New(session.New(), &aws.Config{
		Region:      &Config.Region,
		Credentials: credentials.NewSharedCredentials("", Config.Profile),
	})

	request := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{{
				Action: aws.String(route53.ChangeActionUpsert),
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name: &Config.Domain,
					Type: &Config.Type,
					TTL:  &Config.TTL,
					ResourceRecords: []*route53.ResourceRecord{
						{Value: &newIp},
					},
				}},
			},
			Comment: aws.String("Changed by cambio"),
		},
		HostedZoneId: &Config.Zone,
	}

	return svc.ChangeResourceRecordSets(request)
}
