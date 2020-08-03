package integration

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/convox/console/pkg/crypt"
	"github.com/convox/console/pkg/helpers"
	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
)

const (
	templateURL = "https://console.aws.amazon.com/cloudformation/home#/stacks/create/review"
)

type AWS struct {
	id    string
	oid   string
	token string
}

var (
	regionDisallowed = regexp.MustCompile("[^a-z0-9-]")
	STS              = sts.New(session.New(), nil)
)

// Authorizer

func (i *AWS) Authorize() (string, error) {
	if i.token != "" {
		return i.refreshAuthorization()
	}

	res, err := iam.New(session.New(), nil).GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", errors.WithStack(err)
	}

	if res.User == nil || res.User.Arn == nil {
		return "", errors.WithStack(fmt.Errorf("could not fetch account id"))
	}

	parts := strings.Split(*res.User.Arn, ":")

	if len(parts) < 5 {
		return "", errors.WithStack(fmt.Errorf("could not fetch account id"))
	}

	aid := parts[4]

	u, err := url.Parse(templateURL)
	if err != nil {
		return "", errors.WithStack(err)
	}

	qv := url.Values{}
	qv.Set("stackName", fmt.Sprintf("convox-%s", i.oid))
	qv.Set("templateURL", "https://convox.s3.amazonaws.com/aws/integration.yml")
	qv.Set("param_AccountId", aid)
	qv.Set("param_ExternalId", i.externalID())
	qv.Set("param_IntegrationUrl", fmt.Sprintf("https://%s/organizations/%s/integrations/%s", settings.ExternalHost, i.oid, i.id))
	qv.Set("region", "us-east-1")

	aurl := u.String() + "?" + qv.Encode()

	return aurl, nil
}

func (i *AWS) Exchange(code string, reauthorize bool) (string, map[string]string, error) {
	return "", nil, errors.WithStack(fmt.Errorf("unimplemented"))
}

func (i *AWS) Revoke() error {
	if i.token == "" {
		return nil
	}

	cf, err := i.cloudformation("us-east-1")
	if err != nil {
		return errors.WithStack(err)
	}

	cf.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(fmt.Sprintf("convox-%s", i.oid)),
	})

	return nil
}

// Integration

func (i *AWS) Name() string {
	return "Amazon Web Services"
}

func (i *AWS) Slug() string {
	return "aws"
}

func (i *AWS) Status() (string, error) {
	if i.token == "" {
		return "pending", nil
	} else {
		if _, err := i.Title(nil); err != nil {
			return "disconnected", nil
		} else {
			return "connected", nil
		}
	}
}

func (i *AWS) Title(attrs map[string]string) (string, error) {
	if i.token == "" {
		return "", nil
	}

	parts := strings.Split(i.token, ":")

	if len(parts) < 5 {
		return "", errors.WithStack(fmt.Errorf("invalid token"))
	}

	title := parts[4]

	ii, err := i.iam()
	if err != nil {
		return "", errors.WithStack(err)
	}

	if res, err := ii.ListAccountAliases(&iam.ListAccountAliasesInput{}); err == nil {
		if res.AccountAliases != nil && len(res.AccountAliases) > 0 {
			title = *res.AccountAliases[0]
		}
	}

	return title, nil
}

// Runtime

func (i *AWS) Credentials() (map[string]string, error) {
	cs, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	creds := map[string]string{
		"AWS_ACCESS_KEY_ID":     helpers.DefaultString(cs.AccessKeyId, ""),
		"AWS_SECRET_ACCESS_KEY": helpers.DefaultString(cs.SecretAccessKey, ""),
		"AWS_SESSION_TOKEN":     helpers.DefaultString(cs.SessionToken, ""),
	}

	return creds, nil
}

func (i *AWS) ParameterList() ([]string, error) {
	ps, err := terraformInputs("https://raw.githubusercontent.com/convox/convox/master/terraform/system/aws/variables.tf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ps = helpers.SliceRemove(ps, "name")
	ps = helpers.SliceRemove(ps, "region")

	return ps, nil
}

func (i *AWS) RegionList() ([]string, error) {
	e, err := i.ec2()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err := e.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rs := []string{}

	for _, r := range res.Regions {
		if n := r.RegionName; n != nil {
			rs = append(rs, *n)
		}
	}

	sort.Strings(rs)

	return rs, nil
}

func (i *AWS) config() (*aws.Config, error) {
	cs, err := i.credentials()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &aws.Config{
		Credentials: credentials.NewStaticCredentials(*cs.AccessKeyId, *cs.SecretAccessKey, *cs.SessionToken),
	}

	return c, nil
}

func (i *AWS) credentials() (*sts.Credentials, error) {
	res, err := STS.AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(3600),
		ExternalId:      aws.String(i.externalID()),
		RoleArn:         aws.String(i.token),
		RoleSessionName: aws.String("convox-console"),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res.Credentials, nil
}

func (i *AWS) cloudformation(region string) (*cloudformation.CloudFormation, error) {
	c, err := i.config()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Region = aws.String(region)

	return cloudformation.New(session.New(), c), nil
}

func (i *AWS) describeStacks(region, name string) ([]*cloudformation.Stack, error) {
	cf, err := i.cloudformation(region)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req := &cloudformation.DescribeStacksInput{}

	if name != "" {
		req.StackName = aws.String(name)
	}

	res, err := cf.DescribeStacks(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res.Stacks, nil
}

func (i *AWS) ec2() (*ec2.EC2, error) {
	c, err := i.config()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ec2.New(session.New(), c), nil
}

func (i *AWS) externalID() string {
	return base64.StdEncoding.EncodeToString([]byte(crypt.OneWay(i.oid)))
}

func (i *AWS) iam() (*iam.IAM, error) {
	c, err := i.config()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iam.New(session.New(), c), nil
}

func (i *AWS) refreshAuthorization() (string, error) {
	cf, err := i.cloudformation("us-east-1")
	if err != nil {
		return "", errors.WithStack(err)
	}

	_, err = cf.UpdateStack(&cloudformation.UpdateStackInput{
		Capabilities: []*string{aws.String("CAPABILITY_IAM")},
		Parameters: []*cloudformation.Parameter{
			{ParameterKey: aws.String("AccountId"), UsePreviousValue: aws.Bool(true)},
			{ParameterKey: aws.String("ExternalId"), UsePreviousValue: aws.Bool(true)},
			{ParameterKey: aws.String("IntegrationUrl"), UsePreviousValue: aws.Bool(true)},
		},
		StackName:   aws.String(fmt.Sprintf("convox-%s", i.oid)),
		TemplateURL: aws.String("https://convox.s3.amazonaws.com/aws/integration.yml"),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "No updates are to be performed"):
			return "", nil
		case strings.Contains(err.Error(), "UPDATE_IN_PROGRESS"):
			return "", nil
		default:
			return "", errors.WithStack(err)
		}
	}

	return "", nil
}

// func (i *AWS) RackInstall(region string, name string, version string, params map[string]string) error {
// 	cs, err := i.credentials()
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	args := []string{"rack", "install", "aws", name, "--version", version}

// 	for k, v := range params {
// 		args = append(args, fmt.Sprintf("%s=%s", k, v))
// 	}

// 	cmd := exec.Command("convox", args...)

// 	cmd.Env = []string{
// 		fmt.Sprintf("AWS_DEFAULT_REGION=%s", regionDisallowed.ReplaceAllString(region, "")),
// 		fmt.Sprintf("AWS_REGION=%s", regionDisallowed.ReplaceAllString(region, "")),
// 		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *cs.AccessKeyId),
// 		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *cs.SecretAccessKey),
// 		fmt.Sprintf("AWS_SESSION_TOKEN=%s", *cs.SessionToken),
// 		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
// 		"HOME=/home/convox",
// 	}

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (i *AWS) RackUninstall(region, name string) error {
// 	cs, err := i.credentials()
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	args := []string{"rack", "uninstall", name}

// 	cmd := exec.Command("convox", args...)

// 	cmd.Env = []string{
// 		fmt.Sprintf("AWS_DEFAULT_REGION=%s", regionDisallowed.ReplaceAllString(region, "")),
// 		fmt.Sprintf("AWS_REGION=%s", regionDisallowed.ReplaceAllString(region, "")),
// 		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *cs.AccessKeyId),
// 		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *cs.SecretAccessKey),
// 		fmt.Sprintf("AWS_SESSION_TOKEN=%s", *cs.SessionToken),
// 		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
// 		"HOME=/home/convox",
// 	}

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		return err
// 	}

// 	return nil
// }
