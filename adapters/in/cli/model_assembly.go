// adapters/in/cli/model_assembly.go — rattacher/détacher une liaison à un asset
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"myr-core/domain/model"
)

var modelAssemblyCmd = &cobra.Command{
	Use:   "assembly",
	Short: "Attach or detach an assembly link on an asset",
}

func runModelAssemblyAdd(w io.Writer, svc model.ModelService, assetID, connID string) error {
	if err := svc.AddAssemblyToModule(assetID, connID); err != nil {
		return err
	}
	fmt.Fprintf(w, "Liaison %s rattachée à l'asset %s.\n", connID, assetID)
	return nil
}

var modelAssemblyAddCmd = &cobra.Command{
	Use:   "add <assetID> <connID>",
	Short: "Attach an assembly link to an asset (local draft, no blockchain effect)",
	Long: `Attach an assembly link to an asset. This only updates the local draft — even
if the asset was already submitted before, this reopens it as a draft; the
change reaches the blockchain only at the next "myr model submit".`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelAssemblyAdd(cmd.OutOrStdout(), modelSvc, args[0], args[1])
	},
}

func runModelAssemblyRemove(w io.Writer, svc model.ModelService, assetID, connID string) error {
	if err := svc.RemoveAssemblyFromModule(assetID, connID); err != nil {
		return err
	}
	fmt.Fprintf(w, "Liaison %s détachée de l'asset %s.\n", connID, assetID)
	return nil
}

var modelAssemblyRemoveCmd = &cobra.Command{
	Use:   "remove <assetID> <connID>",
	Short: "Detach an assembly link from an asset (local draft, no blockchain effect)",
	Long: `Detach an assembly link from an asset. This only updates the local draft —
even if the asset was already submitted before, this reopens it as a draft;
the change reaches the blockchain only at the next "myr model submit".`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelAssemblyRemove(cmd.OutOrStdout(), modelSvc, args[0], args[1])
	},
}

func init() {
	modelAssemblyCmd.AddCommand(modelAssemblyAddCmd, modelAssemblyRemoveCmd)
	modelCmd.AddCommand(modelAssemblyCmd)
}
