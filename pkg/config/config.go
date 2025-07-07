package config

// ClusterConfig holds all the configuration data for the LINSTOR cluster.
type ClusterConfig struct {
	ControllerIP   string
	SatelliteIPs  []string
	InstallGUI     bool
	InstallReactor bool
}
