```go
package main

import (
	"fmt"
	"github.com/dome9/dome9-sdk-go/dome9"
	"github.com/dome9/dome9-sdk-go/services/awp"
	awponboarding "github.com/dome9/dome9-sdk-go/services/awp"
	awpawsonboarding "github.com/dome9/dome9-sdk-go/services/awp/aws_onboarding"
	"github.com/dome9/dome9-sdk-go/services/cloudaccounts"
	"github.com/dome9/dome9-sdk-go/services/cloudaccounts/aws"
)

func main() {
	// Pass accessID, secretKey, rawUrl, or set environment variables
	config, _ := dome9.NewConfig("Access-ID", "Secret-Key", "https://api.dome9.com/v2/")
	srv := awpawsonboarding.New(config)
	dome9AwsService := aws.New(config)
	// Get awp aws onboarding data
	appAwsOnboardingDataResponse, _, err := srv.GetOnboardingData()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Get AWP AWS Onboarding Data response type: %T\n Content: %+v", appAwsOnboardingDataResponse, appAwsOnboardingDataResponse)

	// Get cloud account ID
	getCloudAccountQueryParams := cloudaccounts.QueryParameters{ID: "AWS External ID or Cloudguard Account ID"}
	cloudAccountresp, _, err := dome9AwsService.Get(&getCloudAccountQueryParams)
	if err != nil {
		panic(err)
	}
	cloudAccountId := cloudAccountresp.ID

	fmt.Printf("Cloud Account ID: %s\n", cloudAccountId)

	scanMode := "ScanMode"
	
	// Define the request
	awpAwsOnboardingRequest := awpawsonboarding.CreateAWPOnboardingRequestAws{
		CrossAccountRoleName:       "Cross-Account-Role-Name",
		CrossAccountRoleExternalId: "Cross-Account-Role-External-Id",
		ScanMode:                   scanMode, // can be "inAccount", "saas", inAccountHub, inAccountSub
		IsTerraform:                true,
		// following settings available for inAccount and saas scan modes
		AgentlessAccountSettings: &awp_onboarding.AgentlessAccountSettings{
			DisabledRegions:              []string{"eu-west-3"},
			ScanMachineIntervalInHours:   4,
			MaxConcurrenceScansPerRegion: 1,
			SkipFunctionAppsScan:         true,
			InAccountScannerVPC:          "ManagedByAWP",
			ScanAWSLicensedImages:        false,
			CustomTags:                   map[string]string{"test": "testValue"},
		},
	}

	// Create AWP AWS Onboarding

	options := awponboarding.CreateOptions{
		ShouldCreatePolicy: "true",
	}

	awpAwsOnboardingResponse, err := srv.CreateAWPOnboarding(cloudAccountId, awpAwsOnboardingRequest, options)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Onboarding AWP AWS Account with for type: %T\n Content: %+v", awpAwsOnboardingResponse, awpAwsOnboardingResponse)

	getAwpOnboardingDataResponse, _, err := srv.GetAWPOnboarding(cloudAccountId)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Get AWP Onbording account configuration response type: %T\n Content: %+v", getAwpOnboardingDataResponse, getAwpOnboardingDataResponse)

	// Define the new settings
	newSettings := awp_onboarding.AgentlessAccountSettings{
		DisabledRegions:              []string{"us-east-1", "us-west-1", "ap-northeast-1", "ap-southeast-2"},
		ScanMachineIntervalInHours:   10,
		MaxConcurrenceScansPerRegion: 6,
		SkipFunctionAppsScan:         false,
		InAccountScannerVPC:          "ManagedByAWP",
		ScanAWSLicensedImages:        false,
		CustomTags:                   map[string]string{"newTag": "newValue"},
	}

	// Update AWP AWS Onboarding settings
	updateResponse, err := srv.UpdateAWPSettings(cloudAccountId, scanMode, newSettings)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Update AWP AWS Onboarding settings response type: %T\n Content: %+v", updateResponse, updateResponse)

	// Get the updated AWP Onboarding data
	updatedAwpOnboardingDataResponse, _, err := srv.GetAWPOnboarding(cloudAccountId)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Get updated AWP Onboarding account configuration response type: %T\n Content: %+v", updatedAwpOnboardingDataResponse, *updatedAwpOnboardingDataResponse.AgentlessAccountSettings)

	// Define the delete options
	deleteOptions := awponboarding.DeleteOptions{
		ForceDelete: "true",
	}

	// Delete AWP AWS Onboarding
	deleteResponse, err := srv.DeleteAWPOnboarding(cloudAccountId, deleteOptions)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Delete AWP AWS Onboarding response type: %T\n Content: %+v", deleteResponse, deleteResponse)

}
```