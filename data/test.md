% test-app 8
# NAME
test-app - interact with config map and secret manager variables
# SYNOPSIS
test-app


# COMMAND TREE

- [env, e](#env-e)
    - [diff, d](#diff-d)
        - [namespace, ns](#namespace-ns)
        - [ansible, legacy](#ansible-legacy)
    - [view, v](#view-v)
        - [configmap, c](#configmap-c)
        - [ansible, legacy](#ansible-legacy)
        - [namespace, ns](#namespace-ns)
- [s3](#s3)
    - [get](#get)
- [version, v](#version-v)

**Usage**:
```
test-app [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## env, e

Commands to interact with environment variables, both local and on cluster.

### diff, d

Print out detailed diff reports comparing local and running Pod

#### namespace, ns

View diff of local vs. namespace

```
View the diff of the local ansible-vault encrypted Kubenetes Secret file
against a given dotenv file on a pod within a namespace.

The local file will use the contents of the 'data.<accsessor flag>' block.
This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>

It will then grab contents of the dotenv filr on a Pod in a given Namespace.

This defaults to inspecting the '$PWD/.env on' when executing a 'cat' command.
This method uses '/bin/bash -c' as the base command to perform inspection.
```

**--cmd**="": Command to inspect (default: node)

**--configmap, -c**="": Path to configmap.yaml

**--exclude**="": List (csv) of specific env vars to exclude values from display. Set to `""` to remove any exclusions. (default: PATH,SHLVL,HOSTNAME)

**--filter-prefix, -f**="": List of prefixes (csv) used to filter values from display. Set to `""` to remove any filters. (default: npm_,KUBERNETES_,API_PORT)

**--namespace, -n**="": Kube Namespace to list Pods from for inspection

**--secret-suffix**="": Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup (default: _NAME)

**--secrets, -s**="": Path to secrets.yml (default: .docker/secrets.yml)

#### ansible, legacy

View diff of local (ansible encrypted) vs. namespace

```
View the diff of the local ansible-vault encrypted Kubenetes Secret file
against a given dotenv file on a pod within a namespace.

The local file will use the contents of the 'data.<accsessor flag>' block.
This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>

It will then grab contents of the dotenv filr on a Pod in a given Namespace.

This defaults to inspecting the '$PWD/.env on' when executing a 'cat' command.
This method uses '/bin/bash -c' as the base command to perform inspection.
```

**--accessor, -a**="": Accessor key to pull data out of Data block. (default: .env)

**--dotenv**="": Path to `.env` file on Pod (default: $PWD/.env)

**--encrypted-env-file, -e**="": Path to encrypted Kube Secret file

**--namespace, -n**="": Kube Namespace list Pods from for inspection

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

### view, v

View configured environment for either local or running on a Pod

#### configmap, c

View env values based on local settings in a ConfigMap and secrets.yml

>A single line of UsageText

**--configmap, -c**="": Path to configmap.yaml

**--secret-suffix**="": Suffix used to find ENV variables that denote the Secret Manager Secrets to lookup (default: _NAME)

**--secrets, -s**="": Path to secrets.yml (default: .docker/secrets.yml)

#### ansible, legacy

View env values from ansible-vault encrypted Secret file.

>A single line of UsageText

**--accessor, -a**="": Accessor key to pull data out of Data block. (default: .env)

**--encrypted-env-file, -e**="": Path to encrypted Kube Secret file

**--vault-password-file**="": vault password file `VAULT_PASSWORD_FILE`

#### namespace, ns

Interact with env on a running Pod within a Namespace

```
View the diff of the local ansible-vault encrypted Kubenetes Secret file
against a given dotenv file on a pod within a namespace.

The local file will use the contents of the 'data.<accsessor flag>' block.
This defaults to 'data..env'.

Supported ansible-vault encryption version: $ANSIBLE_VAULT;1.1;AES256

Example file structure of decrypted file:

---
apiVersion: v1
kind: Secret
type: Opaque
data:
  .env: <BASE64 ENCODED STRING>

It will then grab contents of the dotenv filr on a Pod in a given Namespace.

This defaults to inspecting the '$PWD/.env on' when executing a 'cat' command.
This method uses '/bin/bash -c' as the base command to perform inspection.
```

**--cmd**="": Command to inspect (default: node)

**--exclude**="": List (csv) of specific env vars to exclude values from display. Set to `""` to remove any exclusions. (default: PATH,SHLVL,HOSTNAME)

**--filter-prefix, -f**="": List of prefixes (csv) used to filter values from display. Set to `""` to remove any filters. (default: npm_,KUBERNETES_,API_PORT)

**--namespace, -n**="": Kube Namespace list Pods from

## s3

simple S3 commands

### get

[object path] [destination path]

## version, v

Print version info
