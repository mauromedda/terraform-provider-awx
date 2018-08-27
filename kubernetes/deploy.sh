#!/bin/bash

rc=$(minikube status 2&>/dev/null; echo $?)
if [[ $rc -ne 0 ]]; then
  echo "You MUST have minikube installed and fully functional"
fi

echo "Restart minikube setting 8 cpus and 8 GiB of RAM"
minikube stop && minikube start --cpus 4 --memory 8192

kubectl apply -f awx-namespace.yml
kubectl apply -f tiller-role-awx.yml
kubectl apply -f tiller-service-account.yaml

helm init --service-account tiller --tiller-namespace awx --kube-context minikube

git clone https://github.com/ansible/awx.git

cat >> awx/installer/inventory <<-EOF
# Kubernetes Install
kubernetes_context=minikube
kubernetes_namespace=awx
tiller_namespace=awx
EOF

ansible-playbook -i awx/installer/inventory awx/installer/install.yml
