# RBAC Slicer

Get some kind of summary info for roles deployed into large cluster.
Large clusters can get a very chaotic and huge collection of rbacs.
This tool gathers gerenal summaries of the situation.

This is a demo that is useful, but needs more work.

For example there might be:

example verbs from the coordination.k8s.io apigroup
- coordination.k8s.io get,watch,list,delete,update,create 1
- coordination.k8s.io leases 21

example resources from the core apigroup
- core configmaps 66
- core configmaps,namespaces,pods,secrets 1
- core configmaps,persistentvolumeclaims,services 1
- core configmaps,secrets 2
