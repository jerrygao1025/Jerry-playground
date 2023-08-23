package helloworld

import (
	"context"
	"go.temporal.io/sdk/temporal"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"

	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

// Workflow is a Hello World workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", name)

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	rand := rand.Intn(100)
	time.Sleep(3 * time.Second)
	if rand <= 90 {
		logger.Info("Activity failed...")
		return "To retry", temporal.NewApplicationError("some retryable error", "SomeType")
	}

	logger.Info("Activity succeed!")
	//logger.Info("Attempts:", activity.GetInfo(ctx).Attempt)
	return "Hello " + name + "!", nil
}
