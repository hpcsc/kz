# kz

A tool like [zoxide](https://github.com/ajeetdsouza/zoxide) but for Kubernetes contexts and namespaces

https://github.com/hpcsc/kz/assets/5293047/48a68950-1b95-410a-8e0f-1ac23321abf1


## Installation

Download from [Github Release](https://github.com/hpcsc/kz/releases) and put in a location in your PATH, .e.g. `/usr/local/bin`

### Bash/ZSH autocomplete

- Clone this repository to local machine, or download `./autocomplete/[bash|zsh]` file
 
#### Bash
- Make sure `bash-completion` is installed, .e.g. by Brew on Mac: `brew install bash-completion`
- Copy `kz` bash completion file to correct location :
```shell
sudo cp /path/to/autocomplete/bash /etc/bash_completion.d/kz
source /etc/bash_completion.d/kz
```

#### ZSH
- Put `source /path/to/autocomplete/zsh` into `.zshrc`

## Examples

```shell
kz ctx sync # copy all context names from your kube config to kz configuration file
kz ns add ns1 ns2 ns3 # track 3 namespaces ns1, ns2, ns3
kz ctx list # list contexts tracked by kz
kz ns list  # list namespaces tracked by kz
kz sys 2  # switch to context matching `sys` and namespace matching the `2`
kz sys  # switch to context matching `sys`
kz ns 2  # switch to namespace matching `2` (using current context)
kz - 2  # same like `kz ns 2`
```
