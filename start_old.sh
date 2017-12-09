#!/bin/bash

# startup the kubernetes cluster
# n1-standard-2: Standard machine type with 2 virtual CPUs and 7.5 GB of memory.
gcloud container clusters create backend --cluster-version=1.8.2-gke.0 --machine-type=n1-standard-1 \
    --num-nodes=1 --zone=europe-west1-c

# Create the database service.
kubectl create -f kubernetes/mysql.yaml

# Create the backend including the service and the ingress controller.
kubectl create -f kubernetes/server_mysql.yaml

# create the load kubernetes cluster
gcloud container clusters create clients --cluster-version=1.8.2-gke.0 --machine-type=n1-standard-2 \
   --num-nodes=3 --zone=europe-west1-c

# switch to the new cluster.
# gcloud container clusters get-credentials backend --zone europe-west1-c --project default-1296

# create the load clients
# DO NOT FORGET TO EDIT THE CLIENT YAML WITH THE CORRECT IP.
kubectl create -f kubernetes/client.yaml

# take it down
gcloud container clusters delete backend
gcloud container clusters delete clients