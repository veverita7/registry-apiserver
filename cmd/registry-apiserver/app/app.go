package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"

	"github.com/veverita7/registry-server/cmd/registry-apiserver/options"
	"github.com/veverita7/registry-server/pkg/server"
)

func NewCommand(stopCh <-chan struct{}) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Short: "Launch registry-server",
		Long:  "Launch registry-server",
		RunE: func(c *cobra.Command, args []string) error {
			return runCommand(opts, stopCh)
		},
	}

	fs := cmd.Flags()
	ofs := opts.Flags()
	for _, f := range ofs.FlagSets {
		fs.AddFlagSet(f)
	}
	local := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(local)
	ofs.FlagSet("logg").AddGoFlagSet(local)

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), ofs, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), ofs, cols)
	})

	return cmd
}

func runCommand(opts *options.Options, stopCh <-chan struct{}) error {
	cnf, err := opts.ServerConfig()
	if err != nil {
		return err
	}

	serv, err := server.NewServer(cnf)
	if err != nil {
		return err
	}

	return serv.RunUntil(stopCh)
}
