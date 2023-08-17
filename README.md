# kz

A tool like [zoxide](https://github.com/ajeetdsouza/zoxide) but for Kubernetes contexts and namespaces

https://github.com/hpcsc/kz/assets/5293047/48a68950-1b95-410a-8e0f-1ac23321abf1


## Installation

Download from [Github Release](https://github.com/hpcsc/kz/releases) and put in a location in your PATH, .e.g. `/usr/local/bin`

## Examples

```shell
kz ctx sync # copy all context names from your kube config to kz configuration file
kz ns add ns1 ns2 ns3 # ask kz to track 3 namespaces ns1, ns2, ns3
kz ctx list # list contexts tracked by kz
kz ns list  # list namespaces tracked by kz
kz sys 2  # ask kz to search for any context matching the text `sys` and any namespace matching the text `2` and switch to those
kz sys  # ask kz to search for any context matching the text `sys` and switch current context to that
kz ns 2  # ask kz to search for any namespace matching the text `2` and switch namespace of current context to that
```
