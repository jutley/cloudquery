package glue

import (
	"context"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
)

func mlTransformTaskRuns() *schema.Table {
	tableName := "aws_glue_ml_transform_task_runs"
	return &schema.Table{
		Name:        tableName,
		Description: `https://docs.aws.amazon.com/glue/latest/webapi/API_TaskRun.html`,
		Resolver:    fetchGlueMlTransformTaskRuns,
		Transform:   transformers.TransformWithStruct(&types.TaskRun{}, transformers.WithPrimaryKeys("TaskRunId")),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:       "ml_transform_arn",
				Type:       arrow.BinaryTypes.String,
				Resolver:   schema.ParentColumnResolver("arn"),
				PrimaryKey: true,
			},
		},
	}
}

func fetchGlueMlTransformTaskRuns(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)
	svc := cl.Services(client.AWSServiceGlue).Glue
	paginator := glue.NewGetMLTaskRunsPaginator(svc, &glue.GetMLTaskRunsInput{
		TransformId: parent.Item.(types.MLTransform).TransformId,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx, func(options *glue.Options) {
			options.Region = cl.Region
		})
		if err != nil {
			return err
		}
		res <- page.TaskRuns
	}
	return nil
}
