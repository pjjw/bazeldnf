package main

import (
	"fmt"
	"os"

	"github.com/rmohr/bazeldnf/cmd/template"
	"github.com/rmohr/bazeldnf/pkg/api/bazeldnf"
	l "github.com/rmohr/bazeldnf/pkg/logger"
	"github.com/rmohr/bazeldnf/pkg/reducer"
	"github.com/rmohr/bazeldnf/pkg/repo"
	"github.com/rmohr/bazeldnf/pkg/sat"
	"github.com/spf13/cobra"
)

type resolveOpts struct {
	in               []string
	lang             string
	nobest           bool
	arch             string
	baseSystem       string
	repofiles        []string
	forceIgnoreRegex []string
}

var resolveopts = resolveOpts{}

func NewResolveCmd() *cobra.Command {

	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "resolves depencencies of the given packages",
		Long:  `resolves dependencies of the given packages with the assumption of a SCRATCH container as install target`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, required []string) error {
			InitLogger(cmd)

			repos := &bazeldnf.Repositories{}
			if len(resolveopts.in) == 0 {
				var err error
				repos, err = repo.LoadRepoFiles(resolveopts.repofiles)
				if err != nil {
					return err
				}
			}
			repo := reducer.NewRepoReducer(repos, resolveopts.in, resolveopts.lang, resolveopts.baseSystem, resolveopts.arch, ".bazeldnf")
			l.Log().Info("Loading packages.")
			if err := repo.Load(); err != nil {
				return fmt.Errorf("failed while loading packages: %w", err)
			}
			l.Log().Info("Loaded packages.")

			l.Log().Infof("Initial reduction of involved packages from required packages %+v.", required)
			matched, involved, err := repo.Resolve(required)
			if err != nil {
				return fmt.Errorf("failed while reducing packages: %w", err)
			}
			l.Log().Infof("packages matched: %d, packages involved: %d", len(matched), len(involved))
			l.Log().Debugf("involved: %+v", involved)
			// involved := repo.DumpPackages()
			// matched := []string{resolveopts.baseSystem}
			for _, r := range required {
				matched = append(matched, r)
			}
			solver := sat.NewResolver(resolveopts.nobest)
			l.Log().Info("Loading involved packages into the resolver.")
			err = solver.LoadInvolvedPackages(involved, resolveopts.forceIgnoreRegex)
			if err != nil {
				return fmt.Errorf("failed while loading involved packages: %w", err)
			}
			l.Log().Info("Adding required packages to the resolver.")
			err = solver.ConstructRequirements(matched)
			if err != nil {
				return fmt.Errorf("failed while constructing requirements: %w", err)
			}
			l.Log().Info("Solving.")
			install, _, forceIgnored, err := solver.Resolve()
			if err != nil {
				return err
			}
			l.Log().Infof("Selected %d packages.", len(install))
			if err := template.Render(os.Stdout, install, forceIgnored); err != nil {
				return err
			}
			return nil
		},
	}

	resolveCmd.Flags().StringArrayVarP(&resolveopts.in, "input", "i", nil, "primary.xml of the repository")
	resolveCmd.Flags().StringVar(&resolveopts.baseSystem, "basesystem", "fedora-release-container", "base system to use (e.g. fedora-release-server, centos-stream-release, ...)")
	resolveCmd.Flags().StringVarP(&resolveopts.arch, "arch", "a", "x86_64", "target architecture")
	resolveCmd.Flags().BoolVarP(&resolveopts.nobest, "nobest", "n", false, "allow picking versions which are not the newest")
	resolveCmd.Flags().StringArrayVarP(&resolveopts.repofiles, "repofile", "r", []string{"repo.yaml"}, "repository information file. Can be specified multiple times. Will be used by default if no explicit inputs are provided.")
	resolveCmd.Flags().StringArrayVar(&resolveopts.forceIgnoreRegex, "force-ignore-with-dependencies", []string{}, "Packages matching these regex patterns will not be installed. Allows force-removing unwanted dependencies. Be careful, this can lead to hidden missing dependencies.")
	// deprecated options
	resolveCmd.Flags().StringVarP(&resolveopts.baseSystem, "fedora-base-system", "f", "fedora-release-container", "base system to use (e.g. fedora-release-server, centos-stream-release, ...)")
	resolveCmd.Flags().MarkDeprecated("fedora-base-system", "use --basesystem instead")
	resolveCmd.Flags().MarkShorthandDeprecated("fedora-base-system", "use --basesystem instead")
	resolveCmd.Flags().MarkShorthandDeprecated("nobest", "use --nobest instead")
	return resolveCmd
}
