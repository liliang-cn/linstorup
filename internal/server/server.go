package server

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"

	"github.com/liliang-cn/linstorup/pkg/config"
	"github.com/liliang-cn/linstorup/pkg/playbook"
)

// Server holds the dependencies for the web server.
type Server struct {
	Templates *template.Template
	Config    *config.ClusterConfig
	Port      int
}

// NewServer creates a new server instance.
func NewServer(port int) (*Server, error) {
	tpls, err := template.ParseGlob("internal/web/*.html")

	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	return &Server{
		Templates: tpls,
		Config:    &config.ClusterConfig{},
		Port:      port,
	}, nil
}

// Start begins listening for web requests.
func (s *Server) Start() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web"))))

	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/setup/controller", s.controllerHandler)
	http.HandleFunc("/setup/satellites", s.satellitesHandler)
	http.HandleFunc("/setup/components", s.componentsHandler)
	http.HandleFunc("/review", s.reviewHandler)
	http.HandleFunc("/deploy", s.deployHandler)
	http.HandleFunc("/deployment-log", s.deploymentLogHandler)
	http.HandleFunc("/deploy-stream", s.deployStreamHandler)

	addr := fmt.Sprintf(":%d", s.Port)
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Reset config at the beginning
	s.Config = &config.ClusterConfig{}
	s.Templates.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server) controllerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.Config.ControllerIP = r.FormValue("controller_ip")
		fmt.Printf("Controller IP saved: %s\n", s.Config.ControllerIP)
		http.Redirect(w, r, "/setup/satellites", http.StatusFound)
		return
	}
	s.Templates.ExecuteTemplate(w, "controller.html", nil)
}

func (s *Server) satellitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		s.Config.SatelliteIPs = r.Form["satellite_ips[]"]
		fmt.Printf("Satellite IPs saved: %v\n", s.Config.SatelliteIPs)
		http.Redirect(w, r, "/setup/components", http.StatusFound)
		return
	}
	s.Templates.ExecuteTemplate(w, "satellites.html", nil)
}

func (s *Server) componentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.Config.InstallGUI = r.FormValue("install_gui") == "on"
		s.Config.InstallReactor = r.FormValue("install_reactor") == "on"
		fmt.Printf("Optional components saved: GUI=%t, Reactor=%t\n", s.Config.InstallGUI, s.Config.InstallReactor)
		http.Redirect(w, r, "/review", http.StatusFound)
		return
	}
	s.Templates.ExecuteTemplate(w, "components.html", nil)
}

func (s *Server) reviewHandler(w http.ResponseWriter, r *http.Request) {
	s.Templates.ExecuteTemplate(w, "review.html", s.Config)
}

func (s *Server) deployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := playbook.GeneratePlaybook(s.Config); err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate playbook: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/deployment-log", http.StatusFound)
}

func (s *Server) deploymentLogHandler(w http.ResponseWriter, r *http.Request) {
	s.Templates.ExecuteTemplate(w, "deploy.html", nil)
}

func (s *Server) deployStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	cmd := exec.Command("ansible-playbook", "-i", "inventory.ini", "playbook.yml")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Fprintf(w, "data: %s\n\n", scanner.Text())
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command finished with error: %v", err)
		fmt.Fprintf(w, "data: Command finished with error: %v\n\n", err)
	} else {
		fmt.Fprintf(w, "data: \n\n--- DEPLOYMENT COMPLETE ---\n\n")
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
