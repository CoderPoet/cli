package cmd

import (
	"github.com/rancher/cli/cliclient"
	projectClient "github.com/rancher/types/client/project/v3"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

type WorkloadData struct {
	ID       string
	Workload projectClient.Workload

	Type  string
	State string
	Image string
	Scale string
}

func WorkloadCommond() cli.Command {
	return cli.Command{
		Name:    "workloads",
		Aliases: []string{"workload"},
		Usage:   "Operations on workload",
		Action:  defaultAction(workloadLs),
		Subcommands: []cli.Command{
			{
				Name:        "ls",
				Usage:       "List Workloads",
				Description: "\nList all workloads in the current project.",
				ArgsUsage:   "None",
				Action:      workloadLs,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "format",
						Usage: "'json','yaml' or Custom format: '{{.Workload.ID}} {{.Workload.Name}}'",
					},
					quietFlag,
				},
			},

			{
				Name:        "create",
				Usage:       "Create a workload",
				Description: "\nCreaate a workload in the current project.",
				ArgsUsage:   "[NEWWORKLOAD...]",
				Action:      workloadCreate,
				Flags: []cli.Flag{
					// todo
				},
			},
		},
	}
}

func workloadLs(ctx *cli.Context) error {
	c, err := GetClient(ctx)
	if err != nil {
		return err
	}

	collection, err := getWorkloadList(ctx, c)
	if err != nil {
		return err
	}

	writer := NewTableWriter([][]string{
		{"ID", "ID"},
		{"NAME", "Workload.Name"},
		{"STATE", "Workload.State"},
		{"TYPE", "Type"},
		{"STATE", "State"},
		{"IMAGE", "Image"},
		{"SCALE", "Scale"},
	}, ctx)

	defer writer.Close()

	for _, item := range collection.Data {
		var scale string

		if item.Scale == nil {
			scale = "-"
		} else {
			scale = strconv.Itoa(int(*item.Scale))
		}

		item.Type = strings.Title(item.Type)

		writer.Write(&WorkloadData{
			ID:       item.ID,
			Workload: item,
			Image:    item.Containers[0].Image,
			Type:     item.Type,
			Scale:    scale,
		})
	}
	return writer.Err()
}

func workloadCreate(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return cli.ShowSubcommandHelp(ctx)
	}

	c, err := GetClient(ctx)
	if err != nil {
		return err
	}

	projectID := c.UserConfig.FocusedProject()
	if ctx.String("project") != "" {
		resource, err := Lookup(c, ctx.String("project"), "project")
		if err != nil {
			return err
		}
		projectID = resource.ID
	}

	// todo 读取yaml来部署
	newWorkload := &projectClient.Workload{
		Name:      ctx.Args().First(),
		ProjectID: projectID,
	}

	_, err = c.ProjectClient.Workload.Create(newWorkload)
	if err != nil {
		return err
	}
	return nil
}

func getWorkloadList(ctx *cli.Context, c *cliclient.MasterClient) (*projectClient.WorkloadCollection, error) {

	collection, err := c.ProjectClient.Workload.List(defaultListOpts(ctx))
	if err != nil {
		return nil, err
	}

	return collection, nil
}
