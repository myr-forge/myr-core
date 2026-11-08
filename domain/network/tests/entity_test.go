// domain/network/tests/entity_test.go — tests unitaires des entités réseau
package network_test

import (
	"strings"
	"testing"

	"myr-core/domain/network"
)

// / @brief  AllChannels liste les canaux du réseau, canal par défaut en tête, sans doublon ni canal vide
// / @input  FabricChannel "sandbox", Channels ["fabrication", "sandbox", ""] ; puis profil sans canal
// / @expect "sandbox,fabrication" ; liste vide pour le profil sans canal
func TestAllChannels_DefaultFirstNoDuplicate(t *testing.T) {
	np := &network.NetworkProfile{FabricChannel: "sandbox", Channels: []string{"fabrication", "sandbox", ""}}
	if got := strings.Join(np.AllChannels(), ","); got != "sandbox,fabrication" {
		t.Errorf("AllChannels = %q", got)
	}
	if got := (&network.NetworkProfile{}).AllChannels(); len(got) != 0 {
		t.Errorf("profil sans canal : %v", got)
	}
}
