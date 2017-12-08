#!/bin/bash

# startup the spanner instance
gcloud spanner instances create ticketshop --config=regional-europe-west1 --description="Ticketshop" --nodes=1

# create the server kubernetes cluster
gcloud container clusters create backend --cluster-version=1.8.2-gke.0 --machine-type=n1-standard-1 \
    --num-nodes=3 --zone=europe-west1-c

# add the service account for spanner as a secret into the kubernetes cluster.
kubectl create secret generic spanner --from-file kubernetes/service-account.json

# Create the backend including the service and the ingress controller.
# CHANGE DATABASE CONFIG!
kubectl create -f kubernetes/server_spanner.yaml

# create the load kubernetes cluster
gcloud container clusters create clients --cluster-version=1.8.2-gke.0 --machine-type=n1-standard-2 \
   --num-nodes=3 --zone=europe-west1-c

# switch to the new cluster.
# kubectl config set-cluster clients

# create the load clients
# DO NOT FORGET TO EDIT THE CLIENT YAML WITH THE CORRECT IP.
kubectl create -f kubernetes/client.yaml

# take it down
gcloud spanner instances delete ticketshop
gcloud container clusters delete backend
gcloud container clusters delete clients
