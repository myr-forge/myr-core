// adapters/in/cli/role.go — gestion des rôles RBAC (permissions personnalisées)
package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"myr-core/domain/role"
)

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage RBAC roles and their permissions",
	Long: `Commands to create, inspect and manage the network's role table and the
permissions attached to each role. The table is published on the blockchain
and applied by every node; changing it requires the wallet of an identity
whose role carries role.admin (myr session).

Until a table is published, the initial table applies: reader, contributor,
auditor and admin. Every role can be modified or deleted except admin, which
always keeps every permission.

A role marked auto-grant is granted on request (myr identity request-role)
without the administrator's approval; a role carrying an administration
permission is never auto-granted.

Permissions: read, write, network.admin, role.admin, identity.admin, admin.`,
}

// grantLabel décrit le réglage d'attribution d'un rôle demandé (UCA08).
func grantLabel(r *role.Role) string {
	if r.AutoGrant {
		return "automatique"
	}
	return "sur approbation"
}

// ── list ──────────────────────────────────────────────────────────────────────

func runRoleList(w io.Writer, svc role.RoleService) error {
	roles, err := svc.List()
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NOM\tINITIAL\tATTRIBUTION\tPERMISSIONS")
	for _, r := range roles {
		perms := make([]string, len(r.Permissions))
		for i, p := range r.Permissions {
			perms[i] = string(p)
		}
		builtIn := ""
		if r.BuiltIn {
			builtIn = "*"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Name, builtIn, grantLabel(r), strings.Join(perms, ", "))
	}
	return tw.Flush()
}

var roleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles and their permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoleList(cmd.OutOrStdout(), roleSvc)
	},
}

// ── show ──────────────────────────────────────────────────────────────────────

var roleShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show a single role's permissions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := roleSvc.Get(args[0])
		if err != nil {
			return err
		}
		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "Rôle : %s\n", r.Name)
		fmt.Fprintf(w, "Table initiale : %v\n", r.BuiltIn)
		fmt.Fprintf(w, "Attribution sur demande : %s\n", grantLabel(r))
		fmt.Fprintln(w, "Permissions :")
		for _, p := range r.Permissions {
			fmt.Fprintf(w, "  - %s\n", p)
		}
		return nil
	},
}

// toPermissions convertit les valeurs du drapeau --permission.
func toPermissions(values []string) []role.Permission {
	perms := make([]role.Permission, len(values))
	for i, p := range values {
		perms[i] = role.Permission(p)
	}
	return perms
}

// ── create ────────────────────────────────────────────────────────────────────

var (
	roleCreatePerms     []string
	roleCreateAutoGrant bool
)

var roleCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Add a role with a given set of permissions",
	Long: `Add a role to the network's role table with an arbitrary subset of the
permission catalogue, then publish the table on the blockchain.

Examples:
  myr role create consumer --permission read --auto-grant
  myr role create manufacturer --permission read --permission write`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := roleSvc.Create(args[0], toPermissions(roleCreatePerms), roleCreateAutoGrant)
		if err != nil {
			return withSignerHint(err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Rôle %q créé.\n", r.Name)
		return nil
	},
}

// ── update ────────────────────────────────────────────────────────────────────

var (
	roleUpdatePerms     []string
	roleUpdateAutoGrant bool
)

var roleUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Change the permissions or the request setting of a role",
	Long: `Change a role of the network's role table, then publish the table on the
blockchain. Only the flags provided are changed. admin cannot be changed.

Examples:
  myr role update contributor --permission read
  myr role update consumer --auto-grant=false`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var change role.Change
		if cmd.Flags().Changed("permission") {
			perms := toPermissions(roleUpdatePerms)
			change.Permissions = &perms
		}
		if cmd.Flags().Changed("auto-grant") {
			change.AutoGrant = &roleUpdateAutoGrant
		}
		r, err := roleSvc.Update(args[0], change)
		if err != nil {
			return withSignerHint(err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Rôle %q mis à jour.\n", r.Name)
		return nil
	},
}

// ── delete ────────────────────────────────────────────────────────────────────

var roleDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a role",
	Long: `Delete a role from the network's role table, then publish the table on the
blockchain. admin cannot be deleted. An identity that still carries the
deleted role keeps it in its certificate but no longer has any permission.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := roleSvc.Delete(args[0]); err != nil {
			return withSignerHint(err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Rôle %q supprimé.\n", args[0])
		return nil
	},
}

// ── init ──────────────────────────────────────────────────────────────────────

func init() {
	roleCmd.AddCommand(roleListCmd, roleShowCmd, roleCreateCmd, roleUpdateCmd, roleDeleteCmd)

	roleCreateCmd.Flags().StringArrayVar(&roleCreatePerms, "permission", nil, "Permission à accorder (répétable) : read, write, network.admin, role.admin, identity.admin, admin")
	roleCreateCmd.Flags().BoolVar(&roleCreateAutoGrant, "auto-grant", false, "Attribuer ce rôle sur simple demande, sans approbation (jamais pour un rôle d'administration)")
	roleUpdateCmd.Flags().StringArrayVar(&roleUpdatePerms, "permission", nil, "Nouvel ensemble de permissions (répétable, remplace l'existant)")
	roleUpdateCmd.Flags().BoolVar(&roleUpdateAutoGrant, "auto-grant", false, "Attribuer ce rôle sur simple demande (true) ou sur approbation (false)")
}
