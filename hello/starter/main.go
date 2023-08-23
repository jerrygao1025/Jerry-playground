package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"Jerry-playground/hello"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	var result string
	ctx := context.TODO()

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello_world_workflowID",
		TaskQueue: "hello-world",
	}

	we, err := c.ExecuteWorkflow(ctx, workflowOptions, helloworld.Workflow, "Jerry")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	err = we.Get(ctx, &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}

	log.Println("Workflow result:", result)
}
