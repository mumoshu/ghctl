package repo

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/Ladicle/ghctl/pkg/config"
	"github.com/Ladicle/ghctl/pkg/github"
	"github.com/Ladicle/ghctl/pkg/util"
	"github.com/spf13/cobra"
)

type getRepoOptions struct {
	ClientGenerator func(token string) *github.Client
	Organization    string
	Repository      string
	Cache           bool
	Output          string
}

func newGetRepoCmd(out, errOut io.Writer) *cobra.Command {
	o := getRepoOptions{
		ClientGenerator: github.NewClient,
	}
	cmd := &cobra.Command{
		Use:   "repo [options] <org/repository>",
		Short: "Outputs repository data",
		Run: func(cmd *cobra.Command, args []string) {
			util.HandleCmdError(o.validate(args), errOut)
			util.HandleCmdError(o.execute(out), errOut)
		},
	}

	cmd.Flags().BoolVar(&o.Cache, "cache", false, "Use cache.")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "yaml", "Output format.")
	return cmd
}

func (o *getRepoOptions) validate(args []string) error {
	//if len(args) < 1 {
	//	return fmt.Errorf("target repository is required")
	//}

	var orgRepo []string
	if len(args) > 0 {
		orgRepo = strings.Split(args[0], "/")
	} else {
		orgRepo = []string{}
	}
	if len(orgRepo) == 2 {
		o.Organization, o.Repository = orgRepo[0], orgRepo[1]
	} else if len(orgRepo) == 1 {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = path.Base(dir)
		o.Organization = dir
		o.Repository = orgRepo[0]
	} else if len(orgRepo) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		o.Organization = dir
	} else {
		return fmt.Errorf("%s is not right 'org/repo' format", args[0])
	}

	if o.Output != "yaml" && o.Output != "json" {
		return fmt.Errorf("%s is unknown output format", o.Output)
	}
	return nil
}

func (o *getRepoOptions) execute(out io.Writer) error {
	var repo *github.Repository
	if o.Cache {
		// TODO: get data from cache
		return fmt.Errorf("cache option have not implemented yet")
	}
	cli := o.ClientGenerator(config.GetCurrentContext().AccessToken)
	if o.Repository != "" {
		newRepo, err := cli.GetRepository(o.Organization, o.Repository)
		if err != nil {
			return err
		}
		repo = newRepo

		// TODO: update cache

		d, err := util.GetPrettyOutput(o.Output, *repo)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s", string(d))
	} else {
		newOrg, err := cli.GetOrganization(o.Organization)
		if err != nil {
			return err
		}
		org := newOrg

		// TODO: update cache

		// TODO: filter repositories by match

		d, err := util.GetPrettyOutput(o.Output, *org)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s", string(d))
	}

	return nil
}
