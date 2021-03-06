# Tidy-Mirror

# Features

Included:

- ocp mirror release
- olm package mirror
- registry login
- catalog source generation and apply to the cluster
- prune index
- push index
- adm catalog mirror 

Not included:

- update trust ca
- icsp updating


# How to use it

You need to prepare a config file before launching the CLI:

In the config file, you can specify the following parameters:

```yaml
ocp_release_version: "4.9.0"            # must be a valid OCP release version with the major.minor.patch format
registry_url: "registry.example.com"    # must be a valid registry URL (maybe with port number)
registry_username: "registry_username"  # must be a valid registry username
registry_password: "registry_password"  # must be a valid registry password
list_packages:                          # list of packages to be mirrored
  - "kubernetes-nmstate-operator"
  - "metallb-operator ocs-operator"
  - "local-storage-operator"
  - "advanced-cluster-management"
extra_images_to_mirror:                  # list of extra images to be mirrored
  - "quay.io/jparrill/registry:3"
  - "registry.access.redhat.com/rhscl/httpd-24-rhel7:latest"
  - "quay.io/ztpfw/ui:latest"


```

Now, we're going to launch the CLI:
```shell
$ t-mirror -h 

$ t-mirror ocp --kubeconfig=/path/to/kubeconfig --config-file=/path/to/config.yaml

$ t-mirror olm --kubeconfig=/path/to/kubeconfig --config-file=/path/to/config.yaml
```


