package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"package-manager/internal/app"
	"package-manager/internal/app/dependencies"
	"package-manager/internal/app/errors"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Packages from liquibase.json",
	Run: func(cmd *cobra.Command, args []string) {

		if global {
			errors.Exit("Can not install packages from liquibase.json globally", 1)
		}

		d := dependencies.Dependencies{}
		d.Read()

		for _, dep := range d.Dependencies {
			p := packs.GetByName(dep.GetName())
			v := p.GetVersion(dep.GetVersion())

			if v.InClassPath(app.ClasspathFiles) {
				errors.Exit(p.Name+" is already installed.", 1)
			}
			if !v.PathIsHTTP() {
				v.CopyToClassPath(app.Classpath)
			} else {
				v.DownloadToClassPath(app.Classpath)
			}
			fmt.Println(v.GetFilename() + " successfully installed in classpath.")
		}

		// Output helper for JAVA_OPTS
		// TODO Test this on windows
		p := "-cp liquibase_libs/*:" + globalpath + "*:" + liquibaseHome + "liquibase.jar"
		fmt.Println()
		fmt.Println("---------- IMPORTANT ----------")
		fmt.Println("Add the following JAVA_OPTS to your CLI:")
		fmt.Println("export JAVA_OPTS=\"" + p + "\"")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}