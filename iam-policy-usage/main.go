package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"ministryofjustice/cloud-platform-automation/utils"

	"github.com/aws/aws-sdk-go-v2/config"
)

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func main() {
	ctx := context.Background()

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Could not load AWS config:", err)
		return
	}

	// Create IAM client
	iamClient := utils.IAMClient(cfg)

	// List policies
	policies, err := utils.ListAllPolicies(iamClient, ctx)
	if err != nil {
		fmt.Println("Error listing policies:", err)
		return
	}

	// Get account ID
	accountID, err := utils.GetAccountID(ctx, cfg)
	if err != nil {
		fmt.Println("Error getting account ID:", err)
		return
	}

	// Get policy last used information
	results, err := utils.GetPolicyLastUsed(iamClient, ctx, policies, accountID)
	if err != nil {
		fmt.Println("Error getting policy last used:", err)
		return
	}

	// byte slice to hold CSV data
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	writer.Write([]string{"PolicyName", "PolicyArn", "LastUsed", "IsAttachedable", "UsedByEntity", "Flag", "Owner", "SlackChannel", "SourceCode", "Application", "GitHubTeam", "BusinessUnit", "Namespace", "InfrastructureSupport"})

	//Write policy usage data
	for _, policy := range results {
		// Convert tags to a string
		var owner string
		var slackchannel string
		var sourcecode string
		var application string
		var githubteam string
		var businessunit string
		var namespace string
		var infrastructuresupport string

		for _, tag := range policy.Tags {
			// equals the value of the key that equals owner
			if derefString(tag.Key) == "owner" {
				owner = derefString(tag.Value)
			}
			if derefString(tag.Key) == "slack-channel" {
				slackchannel = derefString(tag.Value)
			}
			if derefString(tag.Key) == "source-code" {
				sourcecode = derefString(tag.Value)
			}
			if derefString(tag.Key) == "application" {
				application = derefString(tag.Value)
			}
			if derefString(tag.Key) == "GithubTeam" {
				githubteam = derefString(tag.Value)
			}
			if derefString(tag.Key) == "business-unit" {
				businessunit = derefString(tag.Value)
			}
			if derefString(tag.Key) == "namespace" {
				namespace = derefString(tag.Value)
			}
			if derefString(tag.Key) == "infrastructure-support" {
				infrastructuresupport = derefString(tag.Value)
			}
		}

		policyName := derefString(policy.PolicyName)
		policyArn := derefString(policy.PolicyArn)
		lastUsed := ""
		if !policy.LastUsed.IsZero() {
			lastUsed = policy.LastUsed.Format("2006-01-02")
		}
		writer.Write([]string{
			policyName,
			policyArn,
			lastUsed,
			policy.IsAttachedable,
			policy.UsedBy,
			policy.Flag,
			owner,
			slackchannel,
			sourcecode,
			application,
			githubteam,
			businessunit,
			namespace,
			infrastructuresupport,
		})
	}

	writer.Flush()

	s3Client, err := utils.S3Client(ctx, "eu-west-2")
	if err != nil {
		fmt.Println("Error creating S3 client:", err)
		return
	}

	s3Bucket := "cloud-platform-hoodaw-reports"

	err = utils.UpdateS3FileDetails(s3Bucket, "policy_usage.csv", buf.Bytes(), s3Client)
	if err != nil {
		fmt.Println("Error updating S3 file details:", err)
		return
	}
}
