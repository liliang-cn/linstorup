package playbook

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/liliang-cn/linstorup/pkg/config"
)

// GeneratePlaybook creates the ansible playbook and inventory files.
func GeneratePlaybook(cfg *config.ClusterConfig) error {
	pl := &bytes.Buffer{}

	// --- Inventory ---
	fmt.Fprintf(pl, "[controller]\n%s\n\n", cfg.ControllerIP)
	fmt.Fprintf(pl, "[satellites]\n%s\n\n", strings.Join(cfg.SatelliteIPs, "\n"))
	fmt.Fprintf(pl, "[all:vars]\nansible_user=root\n") // Example user, should be configurable

	// Write inventory to a separate file
	if err := os.WriteFile("inventory.ini", pl.Bytes(), 0644); err != nil {
		return err
	}
	pl.Reset() // Reset buffer for playbook

	// --- Playbook ---
	fmt.Fprintln(pl, "- hosts: controller")
	fmt.Fprintln(pl, "  become: yes")
	fmt.Fprintln(pl, "  tasks:")
	fmt.Fprintln(pl, "    - name: Install base controller components")
	fmt.Fprintln(pl, "      ansible.builtin.package:")
	fmt.Fprintln(pl, "        name: ['linstor-controller', 'linstor-client']")
	fmt.Fprintln(pl, "        state: present")

	if cfg.InstallGUI {
		fmt.Fprintln(pl, "    - name: Install LINSTOR GUI")
		fmt.Fprintln(pl, "      ansible.builtin.package:")
		fmt.Fprintln(pl, "        name: 'linstor-gui'")
		fmt.Fprintln(pl, "        state: present")
	}

	if cfg.InstallReactor {
		fmt.Fprintln(pl, "    - name: Install DRBD Reactor")
		fmt.Fprintln(pl, "      ansible.builtin.package:")
		fmt.Fprintln(pl, "        name: 'drbd-reactor'")
		fmt.Fprintln(pl, "        state: present")
	}

	fmt.Fprintln(pl, "\n- hosts: satellites")
	fmt.Fprintln(pl, "  become: yes")
	fmt.Fprintln(pl, "  tasks:")
	fmt.Fprintln(pl, "    - name: Install satellite components")
	fmt.Fprintln(pl, "      ansible.builtin.package:")
	fmt.Fprintln(pl, "        name: 'linstor-satellite'")
	fmt.Fprintln(pl, "        state: present")

	fmt.Printf("\n--- Generated Playbook (playbook.yml) ---\n%s\n--------------------------\n", pl.String())

	return os.WriteFile("playbook.yml", pl.Bytes(), 0644)
}