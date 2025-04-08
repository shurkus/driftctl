package middlewares

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/r3labs/diff/v2"
	"github.com/snyk/driftctl/enumeration/resource"
	dctlresource "github.com/snyk/driftctl/pkg/resource"
	"github.com/snyk/driftctl/pkg/resource/aws"
	testresource "github.com/snyk/driftctl/test/resource"
	"github.com/stretchr/testify/mock"
)

func TestAwsSQSQueuePolicyExpander_Execute(t *testing.T) {
	tests := []struct {
		name               string
		resourcesFromState []*resource.Resource
		expected           []*resource.Resource
		mocks              func(factory *dctlresource.MockResourceFactory)
	}{
		{
			"Inline policy, no aws_sqs_queue_policy attached",
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id":     "foo",
						"policy": "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id": "foo",
					},
				},
				{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"queue_url": "foo",
						"id":        "foo",
						"policy":    "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			func(factory *dctlresource.MockResourceFactory) {
				factory.On("CreateAbstractResource", "aws_sqs_queue_policy", "foo", map[string]interface{}{
					"id":        "foo",
					"queue_url": "foo",
					"policy":    "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
				}).Once().Return(&resource.Resource{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"queue_url": "foo",
						"id":        "foo",
						"policy":    "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				}, nil)

			},
		},
		{
			"No inline policy, aws_sqs_queue_policy attached",
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id": "foo",
					},
				},
				{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "foo",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id": "foo",
					},
				},
				{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "foo",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			func(factory *dctlresource.MockResourceFactory) {},
		},
		{
			"Inline policy duplicate aws_sqs_queue_policy",
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id":     "foo",
						"policy": "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
				{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "foo",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id": "foo",
					},
				},
				{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "foo",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			func(factory *dctlresource.MockResourceFactory) {},
		},
		{
			"Inline policy and aws_sqs_queue_policy",
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id":     "foo",
						"policy": "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
				{
					Id:   "bar",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "bar",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			[]*resource.Resource{
				{
					Id:   "foo",
					Type: aws.AwsSqsQueueResourceType,
					Attrs: &resource.Attributes{
						"id": "foo",
					},
				},
				{
					Id:   "bar",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "bar",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYSQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
				{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "foo",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				},
			},
			func(factory *dctlresource.MockResourceFactory) {
				factory.On("CreateAbstractResource", "aws_sqs_queue_policy", "foo", mock.MatchedBy(func(input map[string]interface{}) bool {
					return input["id"] == "foo"
				})).Once().Return(&resource.Resource{
					Id:   "foo",
					Type: aws.AwsSqsQueuePolicyResourceType,
					Attrs: &resource.Attributes{
						"id":        "foo",
						"queue_url": "foo",
						"policy":    "{\"Id\":\"MYINLINESQSPOLICY\",\"Statement\":[{\"Action\":\"sqs:SendMessage\",\"Effect\":\"Allow\",\"Principal\":\"*\",\"Resource\":\"arn:aws:sqs:eu-west-3:047081014315:foo\",\"Sid\":\"Stmt1611769527792\"}],\"Version\":\"2012-10-17\"}",
					},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			factory := &dctlresource.MockResourceFactory{}
			if tt.mocks != nil {
				tt.mocks(factory)
			}

			repo := testresource.InitFakeSchemaRepository("aws", "5.94.1")
			aws.InitResourcesMetadata(repo)

			m := NewAwsSQSQueuePolicyExpander(factory, repo)
			err := m.Execute(&[]*resource.Resource{}, &tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}
			changelog, err := diff.Diff(tt.expected, tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s got = %v, want %v", strings.Join(change.Path, "."), awsutil.Prettify(change.From), awsutil.Prettify(change.To))
				}
			}
		})
	}
}
