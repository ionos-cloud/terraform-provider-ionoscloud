# Kubernetes Cluster with PostgreSQL DB

This example demonstrates how to provision a Kubernetes cluster with a PostgreSQL database. The setup is designed to support workloads running within the Kubernetes cluster that require database access. The Terraform manifests provision:

- Kubernetes cluster and nodepool
- Kubernetes Secret to access a private Docker registry
- public Service for the Deployment of an example application
- replicated, backed-up PostgreSQL instance
- VLAN for both Kubernetes workloads and the database
- example PostgreSQL database, user and credentials
- connection string for PostgreSQL to be used by the Kubernetes Deployment
