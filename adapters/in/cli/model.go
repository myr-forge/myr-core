// adapters/in/cli/model.go — commandes Cobra pour les modeles 3D
package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"myr-core/domain/model"
	"myr-core/domain/session"
)

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage 3D models",
	Long: `Commands to publish and inspect 3D models on the Myr platform.

Models are stored on the active blockchain network. Each model has a
unique ID and belongs to a channel.`,
}

// ── add ───────────────────────────────────────────────────────────────────────

func runModelAdd(w io.Writer, svc model.ModelService, req model.AddRequest) error {
	m, err := svc.AddFull(req)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Modèle ajouté en brouillon : %s  id=%s (myr model submit %s pour l'engager sur la blockchain)\n", m.Name, m.ID, m.ID)
	return nil
}

var modelAddCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Create a 3D model as a local draft",
	Long: `Create a 3D file as a local draft (no blockchain transaction). The format is
designer-agnostic: .stl, .step, .obj, .3mf and any other format are accepted.
The file is hashed (SHA-256) and stored, and the model is indexed with the
provided name and tags — but it only joins the blockchain network at an
explicit "myr model submit" call. --channel and --category are required: no
default is chosen.

Examples:
  myr model add wheel.stl  --name "Front wheel"  --channel greenchannel --category base --tags "3dprint"
  myr model add pivot.step --name "Left pivot"   --channel greenchannel --tags "mechanical" --category amelioration --parent abc123 --license mit`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		channelID, _ := cmd.Flags().GetString("channel")
		category, _ := cmd.Flags().GetString("category")
		parentID, _ := cmd.Flags().GetString("parent")
		licenseID, _ := cmd.Flags().GetString("license")
		ownerID, _ := cmd.Flags().GetString("owner-id")
		tagsStr, _ := cmd.Flags().GetString("tags")

		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
		}
		ownerID, err := resolveOwner(sessionSvc, ownerID)
		if err != nil {
			return err
		}

		return runModelAdd(cmd.OutOrStdout(), modelSvc, model.AddRequest{
			FilePath:    args[0],
			Name:        name,
			Description: description,
			Category:    model.Category(category),
			ChannelID:   channelID,
			OwnerID:     ownerID,
			ParentID:    parentID,
			LicenseID:   licenseID,
			Tags:        tags,
		})
	},
}

// resolveOwner retourne le propriétaire d'un asset créé en ligne de commande :
// l'identité du compte local CLI (SessionIdentity), dont le wallet signera sa
// soumission sur le ledger (ADR-12). --owner-id, s'il est donné, doit désigner
// cette identité ; sans compte local, il est repris tel quel.
func resolveOwner(ss session.SessionService, ownerFlag string) (string, error) {
	handle, ok := SessionIdentity(ss)
	if !ok {
		return ownerFlag, nil
	}
	if ownerFlag != "" && ownerFlag != handle {
		return "", fmt.Errorf("--owner-id %q ne désigne pas l'identité du compte local %q : un asset appartient à l'identité qui le soumet", ownerFlag, handle)
	}
	return handle, nil
}

// ── submit ────────────────────────────────────────────────────────────────────

// runModelSubmit soumet n'importe quel asset (composant simple ou décomposé) à
// la blockchain. ModelService.Submit route en interne vers la logique
// historique de SubmitModule (RM17, ModuleVersion) dès que l'asset a des
// instances — cette commande n'a pas à connaître cette distinction (ADR-11,
// specs/3-Conception/Conception_intro.md). note n'est utilisée que dans ce cas.
func runModelSubmit(w io.Writer, svc model.ModelService, id, note string) error {
	m, err := svc.Submit(id, note)
	if err != nil {
		return withSignerHint(err)
	}
	fmt.Fprintf(w, "Modèle %s soumis à la blockchain (status=%s).\n", m.ID, m.Status)
	return nil
}

var modelSubmitCmd = &cobra.Command{
	Use:   "submit <id>",
	Short: "Submit a draft asset to the blockchain",
	Long: `Commit the current state of a draft asset (metadata, local interfaces, and —
if the asset has instances — its assembly links) to the blockchain in a
single transaction, then mark it as submitted. Fails if a decomposed asset
has no assembly link (RM17). A submitted asset becomes immutable (fork to
evolve it further).

Example:
  myr model submit abc123def456
  myr model submit abc123def456 --note "v1 release"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		note, _ := cmd.Flags().GetString("note")
		return runModelSubmit(cmd.OutOrStdout(), modelSvc, args[0], note)
	},
}

// ── get ───────────────────────────────────────────────────────────────────────

func runModelGet(w io.Writer, svc model.ModelService, id string) error {
	m, err := svc.Get(id, "")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "ID      : %s\nNom     : %s\nCanal   : %s\nOwner   : %s\nParent  : %s\nLicence : %s\nTags    : %s\nStatus  : %s\n",
		m.ID, m.Name, m.ChannelID, m.OwnerID, m.ParentID, m.LicenseID, strings.Join(m.Tags, ", "), m.Status)
	return nil
}

var modelGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Show a model",
	Long: `Display the metadata of a model by its blockchain identifier.

Example:
  myr model get abc123def456`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelGet(cmd.OutOrStdout(), modelSvc, args[0])
	},
}

// ── list ──────────────────────────────────────────────────────────────────────

func runModelList(w io.Writer, svc model.ModelService, channelID string) error {
	models, err := svc.List(channelID)
	if err != nil {
		if !model.IsDegradedListErr(err) || len(models) == 0 {
			return err
		}
		// Mode dégradé : la blockchain est injoignable mais des brouillons locaux
		// existent — on les affiche quand même, en le signalant explicitement
		// plutôt que de faire comme si la liste était complète (voir REST listGraph).
		fmt.Fprintf(w, "Avertissement : %s (%v)\n\n", model.DegradedListWarning, err)
	}
	return printModelTable(w, models)
}

var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List models in a channel",
	Long: `List all registered models, with optional filtering by channel.

Examples:
  myr model list
  myr model list --channel greenchannel`,
	RunE: func(cmd *cobra.Command, args []string) error {
		channelID, _ := cmd.Flags().GetString("channel")
		return runModelList(cmd.OutOrStdout(), modelSvc, channelID)
	},
}

// ── update ────────────────────────────────────────────────────────────────────

func runModelUpdate(w io.Writer, svc model.ModelService, req model.UpdateRequest) error {
	m, err := svc.UpdateAsset(req)
	if err != nil {
		return withSignerHint(err)
	}
	fmt.Fprintf(w, "Modèle %s mis à jour.\n", m.ID)
	return nil
}

var modelUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update the metadata of an existing asset",
	Long: `Patch the metadata of an existing component or module. Only flags
explicitly provided are changed — everything else is left untouched.

A draft asset is updated locally. On a submitted asset only the name,
description, tags and links can change (written to the blockchain, by its
owner only); its license is frozen — derive a new asset to change it (RM19).

Example:
  myr model update abc123def456 --description "Nouvelle description" --tags "3dprint,mechanical"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		licenseID, _ := cmd.Flags().GetString("license")
		tagsStr, _ := cmd.Flags().GetString("tags")
		linksStr, _ := cmd.Flags().GetString("links")

		var tags, links []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
		}
		if linksStr != "" {
			links = strings.Split(linksStr, ",")
		}

		return runModelUpdate(cmd.OutOrStdout(), modelSvc, model.UpdateRequest{
			ID:          args[0],
			Name:        name,
			Description: description,
			LicenseID:   licenseID,
			Tags:        tags,
			Links:       links,
		})
	},
}

// ── remove ────────────────────────────────────────────────────────────────────

func runModelRemove(w io.Writer, svc model.ModelService, id string) error {
	if err := svc.Remove(id); err != nil {
		return err
	}
	fmt.Fprintf(w, "Modèle %s retiré.\n", id)
	return nil
}

var modelRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove an asset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelRemove(cmd.OutOrStdout(), modelSvc, args[0])
	},
}

// ── children ──────────────────────────────────────────────────────────────────

func runModelChildren(w io.Writer, svc model.ModelService, parentID string) error {
	children, err := svc.GetChildren(parentID)
	if err != nil {
		return err
	}
	return printModelTable(w, children)
}

var modelChildrenCmd = &cobra.Command{
	Use:   "children <parentID>",
	Short: "List the assets derived from a parent asset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runModelChildren(cmd.OutOrStdout(), modelSvc, args[0])
	},
}

func printModelTable(w io.Writer, models []*model.Model3D) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNOM\tCATÉGORIE\tTAGS")
	for _, m := range models {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", m.ID, m.Name, m.Category, strings.Join(m.Tags, ", "))
	}
	return tw.Flush()
}

func init() {
	modelAddCmd.Flags().String("name", "", "Nom du modele (requis)")
	modelAddCmd.Flags().String("description", "", "Description du modèle")
	modelAddCmd.Flags().String("channel", "", "ID du canal (requis)")
	modelAddCmd.Flags().String("category", "", "Catégorie (requis) : base, amelioration, variation, adaptation, derivation, extension, regression")
	modelAddCmd.Flags().String("parent", "", "ID du composant parent (lignée)")
	modelAddCmd.Flags().String("license", "", "ID de licence dans le catalogue")
	modelAddCmd.Flags().String("owner-id", "", "Identité propriétaire (pseudo@org) — par défaut celle du compte local (myr session), seule autorisée à soumettre l'asset")
	modelAddCmd.Flags().String("tags", "", "Tags separes par des virgules")
	modelAddCmd.MarkFlagRequired("name")

	modelListCmd.Flags().String("channel", "", "ID du canal")

	modelSubmitCmd.Flags().String("note", "", "Note de version — utilisée uniquement si l'asset a des instances (ModuleVersion.Note)")

	modelUpdateCmd.Flags().String("name", "", "Nouveau nom")
	modelUpdateCmd.Flags().String("description", "", "Nouvelle description")
	modelUpdateCmd.Flags().String("license", "", "Nouvel ID de licence")
	modelUpdateCmd.Flags().String("tags", "", "Nouveaux tags séparés par des virgules (remplace l'existant)")
	modelUpdateCmd.Flags().String("links", "", "Nouveaux liens séparés par des virgules (remplace l'existant)")

	modelCmd.AddCommand(
		modelAddCmd, modelSubmitCmd, modelGetCmd, modelListCmd,
		modelUpdateCmd, modelRemoveCmd, modelChildrenCmd,
	)
}
