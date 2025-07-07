package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liliang-cn/linstorup/pkg/config"
)

func TestHandlers(t *testing.T) {
	// Dummy variable to ensure config package is considered used
	var _ config.ClusterConfig

	// Change to project root directory for template loading
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	projectRoot := filepath.Join(originalWD, "..", "..") // Go up two levels from internal/server to project root
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Failed to change directory to project root: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Errorf("Failed to change back to original directory: %v", err)
		}
	}()

	srv, err := NewServer(8080) // Use a dummy port for testing
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test indexHandler
	t.Run("IndexHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		srv.indexHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d; got %d", http.StatusOK, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Welcome to Linstor Up") {
			t.Errorf("indexHandler did not render expected content")
		}
	})

	// Test controllerHandler (POST)
	t.Run("ControllerHandlerPost", func(t *testing.T) {
		form := url.Values{}
		form.Add("controller_ip", "1.1.1.1")

		req := httptest.NewRequest("POST", "/setup/controller", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		srv.controllerHandler(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("expected status %d; got %d", http.StatusFound, rr.Code)
		}
		if srv.Config.ControllerIP != "1.1.1.1" {
			t.Errorf("expected controller IP to be '1.1.1.1'; got '%s'", srv.Config.ControllerIP)
		}
	})

	// Test satellitesHandler (POST)
	t.Run("SatellitesHandlerPost", func(t *testing.T) {
		form := url.Values{}
		form.Add("satellite_ips[]", "2.2.2.2")
		form.Add("satellite_ips[]", "3.3.3.3")

		req := httptest.NewRequest("POST", "/setup/satellites", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		srv.satellitesHandler(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("expected status %d; got %d", http.StatusFound, rr.Code)
		}
		if len(srv.Config.SatelliteIPs) != 2 || srv.Config.SatelliteIPs[0] != "2.2.2.2" {
			t.Errorf("unexpected satellite IPs: %v", srv.Config.SatelliteIPs)
		}
	})

	// Test componentsHandler (POST)
	t.Run("ComponentsHandlerPost", func(t *testing.T) {
		form := url.Values{}
		form.Add("install_gui", "on")

		req := httptest.NewRequest("POST", "/setup/components", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		srv.componentsHandler(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("expected status %d; got %d", http.StatusFound, rr.Code)
		}
		if !srv.Config.InstallGUI || srv.Config.InstallReactor {
			t.Errorf("unexpected components selection: GUI=%t, Reactor=%t", srv.Config.InstallGUI, srv.Config.InstallReactor)
		}
	})
}