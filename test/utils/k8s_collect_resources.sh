#!/bin/bash

NAMESPACE=${1:-"keptn"}

# Gathers resource limits by deployment for keptn namespace

KUBE_GET_DEPS="kubectl get deployments -n $NAMESPACE"
echo -e ""
echo -e "| Pod | Container | lim.cpu | lim.mem | req.cpu | req.mem | Images |"
echo -e "|-----|-----------|---------|---------|---------|---------|--------|"
$KUBE_GET_DEPS | sed '1d' | awk '{print $1}' | sort | while read DEPLOYMENT; do
  $KUBE_GET_DEPS $DEPLOYMENT -o jsonpath='{range .spec.template.spec.containers[*]}{"| "}{"'$DEPLOYMENT'"}{" | "}{.name}{" | "}--{.resources.limits.cpu}{" | "}--{.resources.limits.memory}{" | "}--{.resources.requests.cpu}{" | "}--{.resources.requests.memory}{" | "}{.image}{" | "}{"\n"}{end}' | sed -E -e 's/--([0-9])/\1/g' -e 's/--/-/g'
done

echo -e ""
echo -e "Summary (whole cluster):"
echo -e "\`\`\`"
echo -e '$ kubectl describe node | grep -A5 "Allocated"'
kubectl describe node | grep -A5 "Allocated"
echo -e "\`\`\`"
echo -e "Please note: Depending on the setup, the above includes usage for Istio aswell as the Kubernetes control-plane"
echo -e ""

# print PVC data
echo -e "| Name | Size |"
echo -e "|------|------|"
kubectl get pvc -n $NAMESPACE | sed '1d' | awk '{print $1}' | sort | while read PVC; do
  kubectl get pvc $PVC -n $NAMESPACE -o jsonpath='{range .spec}{"| "}{"'$PVC'"}{" | "}--{.resources.requests.storage}{" | "}{"\n"}{end}' | sed -E -e 's/--([0-9])/\1/g' -e 's/--/-/g'
done

echo -e ""
