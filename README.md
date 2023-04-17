# RBAC Slicer

Get some kind of summary info for roles deployed into large cluster.
Large clusters can get a very chaotic and huge collection of rbacs.
This tool gathers gerenal summaries of the situation.

This is a demo that is useful, but needs more work.

For example there might be:

# Examples Data

## verbs sorted by apigroup
- admissionregistration.k8s.io/v1 * 1
- apiextensions.k8s.io * 1
- apiextensions.k8s.io get,list,watch,delete 37
- app.k8s.io * 2
- apps * 13
- apps get 7
- apps get,list,watch 1
- apps get,watch,update,list,create,delete 1
- apps get,watch,update,list,create,delete,patch 2
- apps,extensions get,list,watch 1

## resources from the various apigroups sorted by count
- core configmaps 66
- policy podsecuritypolicies 45
- core events 44
- apiextensions.k8s.io customresourcedefinitions 38
- core secrets 36
- core pods 28
- getambassador.io * 22

# Usage

```
Usage of ./rbac-slicer:
  -c    sort the results by count
  -kubeconfig value
        the path to the kubeconfig
  -l string
        a label selector to filter the results
```
