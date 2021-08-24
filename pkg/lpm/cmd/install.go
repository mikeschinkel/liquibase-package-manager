package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/pkg/lpm"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

// updateCmd represents the `install` command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Packages from liquibase.json",
	Run: func(cmd *cobra.Command, args []string) {

		var d lpm.Dependencies
		var dep lpm.Dependency
		var fn string
		var files lpm.ClasspathFiles
		var cp string
		var err error

		ctx := lpm.GetContextFromCommand(cmd)
		if ctx.FileSource == lpm.GlobalFiles {
			ctx.Error("cannot install packages from liquibase.json globally")
			goto end
		}

		d = lpm.NewDependencies()

		err = d.ReadManifest(ctx)
		if err != nil {
			ctx.Error("unable to read dependencies file; %w", err)
			goto end
		}

		files, cp, err = ctx.GetClasspathFiles()
		if err != nil {
			ctx.Error("unable to get files for classpath '%s'; %w", cp, err)
			goto end
		}

		for _, dep = range d.Dependencies {
			fn, err = maybeInstall(ctx, dep, cp, files)
			if err != nil {
				err = fmt.Errorf("unable to install dependency %s %s; %w",
					dep.GetName(),
					dep.GetVersion(),
					err)
				continue
			}
			fmt.Printf("%s successfully installed in classpath %s.\n",
				fn,
				cp)
		}

		if err != nil {
			ctx.Error(err)
			goto end
		}

		// Output helper for JAVA_OPTS
		ctx.PrintJavaOptsHelper()

	end:
	},
}

func maybeInstall(ctx *lpm.Context, dep lpm.Dependency, cp string, files lpm.ClasspathFiles) (fn string, err error) {

	p := ctx.GetPackageByName(dep.GetName())
	v := p.GetVersion(dep.GetVersion())
	fn = v.GetFilename()

	if files.VersionExists(v) {
		ctx.Error("%s is already installed for %s", p.Name)
		goto end
	}

	err = CopyToClassPath(v, cp)
	if err != nil {
		err = fmt.Errorf("unable to copy ##WHAT?## to class path '%s' for VersionNumber '%s'; %w",
			cp,
			v.Tag,
			err)
		goto end
	}

end:

	return fn, err

}

func CopyToClassPath(v lpm.Version, cp string) (err error) {
	if v.PathIsHTTP() {
		err = v.DownloadToClassPath(cp)
		goto end
	}
	err = v.CopyToClassPath(cp)
end:
	return err
}