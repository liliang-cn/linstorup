# Linstor-Up: A Web-Based LINSTOR Cluster Initializer

This project aims to simplify the initialization of a LINSTOR cluster by providing a web-based, step-by-step wizard.

## Core Functionality

1.  **Web UI:** A Go-based web server will present a user-friendly interface to guide users through the configuration process.
2.  **Configuration:** The user will input essential details, such as the IP addresses for the LINSTOR Controller and Satellite nodes.
3.  **Playbook Generation:** Based on the user's input, the application will dynamically generate an Ansible Playbook.
4.  **Automated Execution:** The application will then execute the generated playbook using `ansible-playbook` to perform tasks like:
    *   Installing LINSTOR software on all nodes.
    *   Configuring the Controller and Satellites.
    *   Initializing the cluster.
5.  **Real-time Feedback:** The web interface will display the output from the Ansible execution, providing real-time status updates to the user.

## Goal

The final product will be a single, self-contained Go binary that can be run on a Linux machine to streamline the setup of a LINSTOR cluster, abstracting away the complexity of manual Ansible commands.
