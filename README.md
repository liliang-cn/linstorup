# LINSTORUP

LINSTORUP is a web-based tool designed to simplify the deployment of LINSTOR clusters using Ansible.
It provides a user-friendly interface to configure your LINSTOR controller and satellite nodes, and then generates and executes Ansible playbooks to automate the installation process.

## Features

-   **Web-based Configuration:** Easily configure LINSTOR deployment parameters through a simple web interface.
-   **Controller & Satellite Setup:** Define controller and satellite IP addresses.
-   **Optional Component Installation:** Choose to install LINSTOR GUI and DRBD Reactor.
-   **Ansible Playbook Generation:** Automatically generates `inventory.ini` and `playbook.yml` based on your input.
-   **Automated Deployment:** Executes the generated Ansible playbook to deploy LINSTOR.
-   **Real-time Deployment Log:** View the deployment progress directly in your browser.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

-   [Go](https://golang.org/doc/install) (version 1.16 or higher)
-   [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) (on the machine where LINSTORUP is run, to execute playbooks)

### Cloning the Repository

```bash
git clone https://github.com/liliang-cn/linstorup.git
cd linstorup
```

### Building the Application

To build the executable:

```bash
go build -o ./bin/linstorup ./cmd/linstorup
```

This will create an executable named `linstorup` (or `linstorup.exe` on Windows) in the `./bin` directory.

### Running the Application

By default, the application runs on port `3374`.

```bash
./bin/linstorup # On Linux/macOS
.\bin\linstorup.exe # On Windows
```

To specify a different port, use the `-port` flag:

```bash
./bin/linstorup -port 8080 # On Linux/macOS
.\bin\linstorup.exe -port 8080 # On Windows
```

Once the server is running, open your web browser and navigate to `http://localhost:3374` (or your specified port).

## Usage

1.  **Controller Setup:** Enter the IP address of your LINSTOR controller.
2.  **Satellite Setup:** Add the IP addresses of your LINSTOR satellite nodes.
3.  **Component Selection:** Choose whether to install the LINSTOR GUI and DRBD Reactor.
4.  **Review:** Review your configuration before deployment.
5.  **Deploy:** Initiate the deployment process and monitor the real-time logs.

## Project Structure

-   `cmd/linstorup`: Contains the main application entry point.
-   `internal/server`: Handles the web server, HTTP routes, and request handlers.
-   `internal/web`: Stores HTML templates and static assets for the web interface.
-   `pkg/config`: Defines the `ClusterConfig` structure for LINSTOR deployment parameters.
-   `pkg/playbook`: Contains logic for generating Ansible `inventory.ini` and `playbook.yml` files.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details. (Note: A LICENSE file is not yet present in the repository. Please add one if you intend to specify a license.)
