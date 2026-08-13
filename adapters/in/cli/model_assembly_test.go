// adapters/in/cli/model_assembly_test.go — tests des handlers CLI assembly
package cli

import (
	"bytes"
	"errors"
	"testing"
)

/// @brief  runModelAssemblyAdd rattache une liaison à l'asset (UCMOD02)
/// @input  assetID="c1", connID="conn-1"
/// @expect AddAssemblyToModule reçoit les deux ids
func TestRunModelAssemblyAdd_NominalCase(t *testing.T) {
	var gotAsset, gotConn string
	svc := &mockModelSvc{addAssemblyToModule: func(assetID, connID string) error {
		gotAsset, gotConn = assetID, connID
		return nil
	}}
	var buf bytes.Buffer
	if err := runModelAssemblyAdd(&buf, svc, "c1", "conn-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAsset != "c1" || gotConn != "conn-1" {
		t.Fatalf("unexpected forwarded args: asset=%q conn=%q", gotAsset, gotConn)
	}
}

/// @brief  runModelAssemblyAdd propage l'erreur du service
/// @input  service retournant une erreur
/// @expect L'erreur est retournée
func TestRunModelAssemblyAdd_ServiceError_Rejected(t *testing.T) {
	wantErr := errors.New("asset introuvable")
	svc := &mockModelSvc{addAssemblyToModule: func(string, string) error { return wantErr }}
	var buf bytes.Buffer
	err := runModelAssemblyAdd(&buf, svc, "c1", "conn-1")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

/// @brief  runModelAssemblyRemove détache une liaison de l'asset
/// @input  assetID="c1", connID="conn-1"
/// @expect RemoveAssemblyFromModule reçoit les deux ids
func TestRunModelAssemblyRemove_NominalCase(t *testing.T) {
	var gotAsset, gotConn string
	svc := &mockModelSvc{removeAssemblyFromModule: func(assetID, connID string) error {
		gotAsset, gotConn = assetID, connID
		return nil
	}}
	var buf bytes.Buffer
	if err := runModelAssemblyRemove(&buf, svc, "c1", "conn-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAsset != "c1" || gotConn != "conn-1" {
		t.Fatalf("unexpected forwarded args: asset=%q conn=%q", gotAsset, gotConn)
	}
}
