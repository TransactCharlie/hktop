# hktop
Console UI kubernetes diagnostic tool

```
┌─Summary────────────────────────┐┌────────────────────────────────────────────────────────────────┐
│Nodes                       │1  ││Hello World                                                     │
│Pods                        │19 ││                                                                │
│Deployments                 │9  ││                                                                │
│Services                    │5  ││                                                                │
│Daemon Sets                 │1  ││                                                                │
│Persistent Volumes          │0  ││                                                                │
└────────────────────────────────┘└────────────────────────────────────────────────────────────────┘
┌─K8S Nodes────────────────────────────┐┌─K8S Pods─────────────────────────────────────────────────┐
│minikube                              ││ubuntu-charlie-191-699bc747b9-4drr2                       │
│                                      ││ubuntu-charlie-199-6bd6d6f5d5-csb6j                       │
│                                      ││ubuntu-charlie-76994d54bf-v4dbs                           │
│                                      ││ubuntu-charlie-1-68b44fdd85-w4v24                         │
│                                      ││ubuntu-charlie-100-57955fc76f-gbmk8                       │
│                                      ││ubuntu-charlie-199-6bd6d6f5d5-qlnxt                       │
│                                      ││ubuntu-charlie-2-56bff68979-hjlrf                         │
│                                      ││ubuntu-charlie-6-855c8c64dc-z28qf                         │
│                                      ││coredns-fb8b8dccf-fxvtz                                   │
│                                      ││coredns-fb8b8dccf-sprw5                                   │
│                                      ││etcd-minikube                                             │
│                                      ││heapster-qxf2x                                           ▼│
│                                      │└──────────────────────────────────────────────────────────┘
│                                      │┌─Pods per Namespace───────────────────────────────────────┐
│                                      ││                                                          │
│                                      ││                          ░░░░░░░                         │
│                                      ││                         ░░░░░alpha: 3                    │
│                                      ││                         ░░kube-system: 11                │
│                                      ││                         ░░░░░░charlie: 5                 │
│                                      ││                          ░░░░░░░                         │
│                                      ││                                                          │
└──────────────────────────────────────┘└──────────────────────────────────────────────────────────┘
```
