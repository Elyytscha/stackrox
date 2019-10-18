package central

import (
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/roxctl/central/db"
	"github.com/stackrox/rox/roxctl/central/debug"
	"github.com/stackrox/rox/roxctl/central/deploy"
	"github.com/stackrox/rox/roxctl/central/license"
	"github.com/stackrox/rox/roxctl/central/userpki"
)

// Command defines the central command tree
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "central",
		Short: "Central is the list of commands that pertain to the Central service",
		Long:  "Central is the list of commands that pertain to the Central service",
		Run: func(c *cobra.Command, _ []string) {
			_ = c.Help()
		},
	}

	c.AddCommand(
		deploy.Command(),
		db.Command(),
		debug.Command(),
		license.Command(),
		userpki.Command(),
	)
	return c
}
