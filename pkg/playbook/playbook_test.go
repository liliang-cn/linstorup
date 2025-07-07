package playbook

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/liliang-cn/linstorup/pkg/config"
)

// TestGeneratePlaybook tests the logic of playbook and inventory file generation.
func TestGeneratePlaybook(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		config         config.ClusterConfig
		expectInPlaybook []string
		expectInInventory []string
	}{
		{
			name: "Base Install",
			config: config.ClusterConfig{
				ControllerIP: "192.168.1.10",
				SatelliteIPs: []string{"192.168.1.11", "192.168.1.12"},
			},
			expectInPlaybook: []string{"linstor-controller", "linstor-satellite"},
			expectInInventory: []string{"[controller]", "192.168.1.10", "[satellites]", "192.168.1.11", "192.168.1.12"},
		},
		{
			name: "Install with GUI",
			config: config.ClusterConfig{
				ControllerIP: "10.0.0.1",
				SatelliteIPs: []string{"10.0.0.2"},
				InstallGUI:     true,
			},
			expectInPlaybook: []string{"linstor-gui"},
			expectInInventory: []string{"10.0.0.1", "10.0.0.2"},
		},
		{
			name: "Install with Reactor",
			config: config.ClusterConfig{
				ControllerIP: "10.0.0.1",
				SatelliteIPs: []string{"10.0.0.2"},
				InstallReactor: true,
			},
			expectInPlaybook: []string{"drbd-reactor"},
			expectInInventory: []string{"10.0.0.1", "10.0.0.2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Cleanup generated files after the test
			t.Cleanup(func() {
				os.Remove("playbook.yml")
				os.Remove("inventory.ini")
			})

			err := GeneratePlaybook(&tc.config)
			if err != nil {
				t.Fatalf("GeneratePlaybook() failed: %v", err)
			}

			// Check playbook content
			playbookBytes, err := ioutil.ReadFile("playbook.yml")
			if err != nil {
				t.Fatalf("Failed to read playbook.yml: %v", err)
			}
			playbookContent := string(playbookBytes)
			for _, expected := range tc.expectInPlaybook {
				if !strings.Contains(playbookContent, expected) {
					t.Errorf("Expected to find '%s' in playbook, but did not", expected)
				}
			}

			// Check inventory content
			inventoryBytes, err := ioutil.ReadFile("inventory.ini")
			if err != nil {
				t.Fatalf("Failed to read inventory.ini: %v", err)
			}
			inventoryContent := string(inventoryBytes)
			for _, expected := range tc.expectInInventory {
				if !strings.Contains(inventoryContent, expected) {
					t.Errorf("Expected to find '%s' in inventory, but did not", expected)
				}
			}
		})
	}
}
