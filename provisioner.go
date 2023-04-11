package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/joho/godotenv"
	"log"
	"os"
	"provisioner/library"
)

var (
	RoleArn = ""
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	RoleArn = os.Getenv("ROLE_ARN")
	if RoleArn == "" {
		log.Fatalf("ROLE_ARN can not be empty, please check .env file")
	}
}

type Provisioner interface {
	Provision(ctx context.Context, tf library.Terraform) error
	Deprovision(ctx context.Context, tf library.Terraform) error
}

type AWSProvisioner struct{}

func (p *AWSProvisioner) Provision(ctx context.Context, tf library.Terraform) error {

	// Create VPC network using Terraform
	err := tf.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize Terraform: %v", err)
	}

	err = tf.Apply(ctx)
	if err != nil {
		return fmt.Errorf("failed to apply Terraform configuration: %v", err)
	}
	return nil
}

func (p *AWSProvisioner) Deprovision(ctx context.Context, tf library.Terraform) error {
	err := tf.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize Terraform: %v", err)
	}
	err = tf.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("failed to destroy Terraform resources: %v", err)
	}

	return nil
}

func main() {
	commands := os.Args[1:]
	ctx := context.Background()

	// Initial credentials loaded from SDK's default credential chain.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	// Create the credentials from AssumeRoleProvider to assume the role referenced by the ARN.
	stsSvc := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, RoleArn)

	cfg.Credentials = aws.NewCredentialsCache(creds)
	setupEnv(ctx, creds, cfg)

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.4.4")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	tf := getTerraformProvider(execPath)

	var provisioner Provisioner = &AWSProvisioner{}
	if len(commands) == 0 {
		log.Fatalf("use one of available commands apply/destroy")
	}
	switch commands[0] {
	case "apply":
		// Provision sandbox infrastructure
		vpcID := provisioner.Provision(ctx, tf)
		if err != nil {
			fmt.Println("Failed to provision sandbox infrastructure:", err)
			return
		}
		fmt.Println("Sandbox VPC ID:", vpcID)
	case "destroy":
		//Deprovision sandbox infrastructure
		err = provisioner.Deprovision(ctx, tf)
		if err != nil {
			fmt.Println("Failed to deprovision sandbox infrastructure:", err)
		}
		fmt.Println("Sandbox resources deprovisioned successfully.")
	default:
		log.Fatalf("unknown command")
	}

}

func setupEnv(ctx context.Context, creds *stscreds.AssumeRoleProvider, cfg aws.Config) {
	//generates a new set of temporary credentials using STS.
	sess, err := creds.Retrieve(ctx)
	if err != nil {
		log.Fatalf("error creating AWS session: %s", err)
	}
	os.Setenv("AWS_REGION", cfg.Region)
	os.Setenv("AWS_ACCESS_KEY_ID", sess.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", sess.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", sess.SessionToken)
}

func getTerraformProvider(execPath string) library.Terraform {
	provider, err := tfexec.NewTerraform("terraform", execPath)
	provider.SetStdout(os.Stdout)
	provider.SetStderr(os.Stdout)
	if err != nil {
		log.Fatalf("error initializing Terraform: %s", err)
	}
	tf := library.New(provider)
	if err != nil {
		log.Fatalf("Failed to create Terraform client: %v", err)
	}

	return tf

}
