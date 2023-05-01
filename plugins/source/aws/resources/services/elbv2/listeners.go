package elbv2

import (
	"context"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/transformers"
)

func listeners() *schema.Table {
	tableName := "aws_elbv2_listeners"
	return &schema.Table{
		Name:        tableName,
		Description: `https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_Listener.html`,
		Resolver:    fetchElbv2Listeners,
		Multiplex:   client.ServiceAccountRegionMultiplexer(tableName, "elasticloadbalancing"),
		Transform:   transformers.TransformWithStruct(&types.Listener{}),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:     "arn",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("ListenerArn"),
				CreationOptions: schema.ColumnCreationOptions{
					PrimaryKey: true,
				},
			},
			{
				Name:     "tags",
				Type:     schema.TypeJSON,
				Resolver: resolveElbv2listenerTags,
			},
		},

		Relations: []*schema.Table{
			listenerCertificates(),
			listenerRules(),
		},
	}
}

func fetchElbv2Listeners(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	lb := parent.Item.(types.LoadBalancer)
	config := elbv2.DescribeListenersInput{
		LoadBalancerArn: lb.LoadBalancerArn,
	}
	c := meta.(*client.Client)
	svc := c.Services().Elasticloadbalancingv2
	paginator := elbv2.NewDescribeListenersPaginator(svc, &config)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx, func(options *elbv2.Options) {
			options.Region = c.Region
		})
		if err != nil {
			if c.IsNotFoundError(err) {
				return nil
			}
			return err
		}
		res <- page.Listeners
	}
	return nil
}

func resolveElbv2listenerTags(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	region := meta.(*client.Client).Region
	cl := meta.(*client.Client)
	svc := cl.Services().Elasticloadbalancingv2
	listener := resource.Item.(types.Listener)
	tagsOutput, err := svc.DescribeTags(ctx, &elbv2.DescribeTagsInput{
		ResourceArns: []string{
			*listener.ListenerArn,
		},
	}, func(o *elbv2.Options) {
		o.Region = region
	})
	if err != nil {
		return err
	}
	if len(tagsOutput.TagDescriptions) == 0 {
		return nil
	}
	tags := make(map[string]*string)
	for _, td := range tagsOutput.TagDescriptions {
		for _, s := range td.Tags {
			tags[*s.Key] = s.Value
		}
	}

	return resource.Set(c.Name, tags)
}
