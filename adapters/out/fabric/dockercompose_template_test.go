// adapters/out/fabric/dockercompose_template_test.go
// Test boîte blanche (package fabric) : renderDockerCompose et ses paramètres ne
// sont pas exportés — même raison que l'exception des tests de adapters/in/cli/.
package fabric

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type composeService struct {
	Ports       []string `yaml:"ports"`
	Environment []string `yaml:"environment"`
	Command     string   `yaml:"command"`
}

func sampleComposeParams() dockerComposeParams {
	return dockerComposeParams{
		NetworkName:     "myr-net-test",
		CryptoDir:       "/tmp/crypto",
		CAName:          "ca-org1",
		CAAdminID:       "admin",
		CAAdminPwd:      "CASECRETTEST",
		CAPort:          7054,
		CouchDBAdminPwd: "DBSECRETTEST",
		OrgMSPID:        "Org1MSP",
		Orderers:        []dcOrdererParams{{Name: "orderer.example.com", Port: 7050, AdminPort: 7053, OpPort: 9443}},
		Peers: []dcPeerParams{
			{Name: "peer0.org1.example.com", Port: 7051, ChaincodePort: 7052, OperationPort: 9444, CouchDBName: "couchdb0"},
			{Name: "peer1.org1.example.com", Port: 7061, ChaincodePort: 7062, OperationPort: 9445, CouchDBName: "couchdb1"},
		},
	}
}

func renderComposeServices(t *testing.T) (string, map[string]composeService) {
	t.Helper()
	out, err := renderDockerCompose(sampleComposeParams())
	if err != nil {
		t.Fatalf("renderDockerCompose: %v", err)
	}
	var doc struct {
		Services map[string]composeService `yaml:"services"`
	}
	if err := yaml.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("docker-compose.yaml généré invalide: %v\n%s", err, out)
	}
	return out, doc.Services
}

func hasEnv(s composeService, entry string) bool {
	for _, e := range s.Environment {
		if e == entry {
			return true
		}
	}
	return false
}

// / @brief  renderDockerCompose injecte les secrets générés au bootstrap, sans aucun identifiant par défaut
// / @input  CAAdminPwd et CouchDBAdminPwd distincts, 2 peers (donc 2 bases CouchDB)
// / @expect aucune occurrence de "adminpw" ; fabric-ca-server démarre avec le secret CA ; chaque base CouchDB et chaque peer utilisent le secret CouchDB
func TestRenderDockerCompose_NoDefaultCredentials(t *testing.T) {
	out, services := renderComposeServices(t)

	if strings.Contains(out, "adminpw") {
		t.Errorf("le docker-compose généré contient encore l'identifiant par défaut \"adminpw\"")
	}
	ca, ok := services["ca.ca-org1"]
	if !ok {
		t.Fatalf("service CA absent du docker-compose généré")
	}
	if !strings.Contains(ca.Command, "-b admin:CASECRETTEST") {
		t.Errorf("commande CA : got %q, want le secret généré passé via -b", ca.Command)
	}
	for _, db := range []string{"couchdb0", "couchdb1"} {
		if !hasEnv(services[db], "COUCHDB_PASSWORD=DBSECRETTEST") {
			t.Errorf("%s : COUCHDB_PASSWORD n'est pas le secret généré (env=%v)", db, services[db].Environment)
		}
	}
	for _, peer := range []string{"peer0.org1.example.com", "peer1.org1.example.com"} {
		if !hasEnv(services[peer], "CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=DBSECRETTEST") {
			t.Errorf("%s : mot de passe CouchDB du peer différent du secret généré (env=%v)", peer, services[peer].Environment)
		}
	}
}

// / @brief  renderDockerCompose ne publie jamais le port d'une base CouchDB sur l'hôte
// / @input  2 peers, donc 2 bases CouchDB
// / @expect aucune base CouchDB n'a de section ports ; chaque peer joint sa base par le réseau Docker interne
func TestRenderDockerCompose_CouchDBNotPublished(t *testing.T) {
	_, services := renderComposeServices(t)

	for i, db := range []string{"couchdb0", "couchdb1"} {
		svc, ok := services[db]
		if !ok {
			t.Fatalf("service %s absent du docker-compose généré", db)
		}
		if len(svc.Ports) != 0 {
			t.Errorf("%s publie des ports sur l'hôte : %v", db, svc.Ports)
		}
		peer := []string{"peer0.org1.example.com", "peer1.org1.example.com"}[i]
		if !hasEnv(services[peer], "CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS="+db+":5984") {
			t.Errorf("%s ne joint pas %s par le réseau Docker interne (env=%v)", peer, db, services[peer].Environment)
		}
	}
}
