// adapters/in/cli/model_link.go — liaisons d'assemblage entre interfaces
package cli

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"myr-core/domain/model"
)

var modelLinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Manage assembly links between interfaces",
}

// ── add ───────────────────────────────────────────────────────────────────────

func runModelLinkAdd(w io.Writer, svc model.ModelService, from, to, label, fromInstance, toInstance, fastener string) error {
	if from == "" || to == "" {
		return fmt.Errorf("--from et --to sont requis")
	}
	conn, err := svc.AddAssemblyLink(from, to, label, fromInstance, toInstance, fastener)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Liaison créée : id=%s\n", conn.ID)
	return nil
}

var modelLinkAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create an assembly link between two interfaces",
	Long: `Create a link between two physical interfaces, directly or through a
fastener asset (screw, cable, clip...). The link is refused when the two
interfaces are incompatible (or, with a fastener, when the fastener does not
fit both), when an interface of the given instance is already linked, or when
an interface would be linked to itself.

Examples:
  myr model link add --from iface-1 --to iface-2
  myr model link add --from iface-1 --to iface-2 --fastener asset-screw-m3 --label "vis M3"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		label, _ := cmd.Flags().GetString("label")
		fastener, _ := cmd.Flags().GetString("fastener")
		fromInstance, _ := cmd.Flags().GetString("from-instance")
		toInstance, _ := cmd.Flags().GetString("to-instance")
		return runModelLinkAdd(cmd.OutOrStdout(), modelSvc, from, to, label, fromInstance, toInstance, fastener)
	},
}

// ── remove ────────────────────────────────────────────────────────────────────

func runModelLinkRemove(w io.Writer, svc model.ModelService, id string) error {
	if err := svc.RemoveConnection(id); err != nil {
		return err
	}
	fmt.Fprintf(w, "Liaison %s supprimée.\n", id)
	return nil
}

var modelLinkRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove an assembly link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelLinkRemove(cmd.OutOrStdout(), modelSvc, args[0])
	},
}

// ── list ──────────────────────────────────────────────────────────────────────

func runModelLinkList(w io.Writer, svc model.ModelService) error {
	conns, err := svc.ListConnections()
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tDE\tVERS\tLABEL\tACCROCHE\tINCOMPATIBLE")
	for _, c := range conns {
		incompatible := ""
		if c.Incompatible {
			incompatible = "*"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", c.ID, c.From, c.To, c.Label, c.FastenerAssetID, incompatible)
	}
	return tw.Flush()
}

var modelLinkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all assembly links",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelLinkList(cmd.OutOrStdout(), modelSvc)
	},
}

func init() {
	modelLinkAddCmd.Flags().String("from", "", "ID interface source (requis)")
	modelLinkAddCmd.Flags().String("to", "", "ID interface cible (requis)")
	modelLinkAddCmd.Flags().String("label", "", "Label de la liaison")
	modelLinkAddCmd.Flags().String("fastener", "", "ID de l'asset d'accroche (optionnel)")
	modelLinkAddCmd.Flags().String("from-instance", "", "ID instance de module pour l'interface source (optionnel)")
	modelLinkAddCmd.Flags().String("to-instance", "", "ID instance de module pour l'interface cible (optionnel)")

	modelLinkCmd.AddCommand(modelLinkAddCmd, modelLinkRemoveCmd, modelLinkListCmd)
	modelCmd.AddCommand(modelLinkCmd)
}
