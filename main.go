package main

import (
	"context"
	"net/http"

	"github.com/puppetlabs/pct/internal/pkg/pct_config_processor"
	"github.com/puppetlabs/pct/pkg/exec_runner"

	cmd_build "github.com/puppetlabs/pct/cmd/build"
	"github.com/puppetlabs/pct/cmd/completion"
	"github.com/puppetlabs/pct/cmd/explain"
	cmd_install "github.com/puppetlabs/pct/cmd/install"
	"github.com/puppetlabs/pct/cmd/new"
	"github.com/puppetlabs/pct/cmd/root"
	appver "github.com/puppetlabs/pct/cmd/version"
	"github.com/puppetlabs/pct/pkg/build"
	"github.com/puppetlabs/pct/pkg/gzip"
	"github.com/puppetlabs/pct/pkg/install"
	"github.com/puppetlabs/pct/pkg/tar"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	ctx := context.Background()

	var rootCmd = root.CreateRootCommand()

	var verCmd = appver.CreateVersionCommand(version, date, commit)
	v := appver.Format(version, date, commit)
	rootCmd.Version = v
	rootCmd.SetVersionTemplate(v)
	rootCmd.AddCommand(verCmd)

	rootCmd.AddCommand(completion.CreateCompletionCommand())

	// afero setup
	fs := afero.NewOsFs()
	afs := afero.Afero{Fs: fs}
	iofs := afero.IOFS{Fs: fs}

	// build
	buildCmd := cmd_build.BuildCommand{
		ProjectType: "template",
		Builder: &build.Builder{
			Tar:             &tar.Tar{AFS: &afero.Afero{Fs: fs}},
			Gzip:            &gzip.Gzip{AFS: &afero.Afero{Fs: fs}},
			AFS:             &afero.Afero{Fs: fs},
			ConfigProcessor: &pct_config_processor.PctConfigProcessor{AFS: &afero.Afero{Fs: fs}},
			ConfigFile:      "pct-config.yml",
		},
	}
	rootCmd.AddCommand(buildCmd.CreateCommand())

	// install
	installCmd := cmd_install.InstallCommand{
		PctInstaller: &install.Installer{
			Tar:        &tar.Tar{AFS: &afs},
			Gunzip:     &gzip.Gunzip{AFS: &afs},
			AFS:        &afs,
			IOFS:       &iofs,
			HTTPClient: &http.Client{},
			Exec:       &exec_runner.Exec{},
			ConfigProcessor: &pct_config_processor.PctConfigProcessor{
				AFS: &afs,
			},
			ConfigFileName: "pct-config.yml",
		},
		AFS: &afs,
	}
	rootCmd.AddCommand(installCmd.CreateCommand())

	// new
	rootCmd.AddCommand(new.CreateCommand())

	// explain
	rootCmd.AddCommand(explain.CreateCommand())

	// initialize
	cobra.OnInitialize(root.InitLogger, root.InitConfig)

	// execute called command
	err := rootCmd.ExecuteContext(ctx)

	// Handle exiting with/out errors.
	cobra.CheckErr(err)
}
