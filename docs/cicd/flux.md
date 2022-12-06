---
id: flux
title: Flux
sidebar_label: Flux
sidebar_position: 1
---

## Bootstrap

### Export the credentials

```bash
export GITHUB_TOKEN=<your-token>
export GITHUB_USER=<your-username>
```

### Precheck the Kubernetes cluster

```bash
flux check --pre
```

### Install Flux onto the cluster

```bash
flux bootstrap github \
  --owner=$GITHUB_USER \
  --repository=fleet-infra \
  --branch=main \
  --path=./clusters/my-cluster \
  --personal
```

The bootstrap command above does following:

+ Creates a git repository fleet-infra on your GitHub account
+ Adds Flux component manifests to the repository
+ Deploys Flux Components to your Kubernetes Cluster
+ Configures Flux components to track the path /clusters/my-cluster/ in the repository


## Triggering Flux

Flux will reconcile the repo and its resources at the interval that you've set in the Flux Kustomizations. If you're doing active development and prefer not to wait, you can force Flux to do a run with a combination of flux reconcile source and flux reconcile kustomization. For example, to force a run on the infrastructure resources:

```
➤ flux reconcile source git flux-system
► annotating GitRepository flux-system in flux-system namespace
✔ GitRepository annotated
◎ waiting for GitRepository reconciliation
✔ fetched revision main/8ad0d7b302be3fd253568ad84354a79385c98908

➤ flux reconcile kustomization infrastructure
► annotating Kustomization infrastructure in flux-system namespace
✔ Kustomization annotated
◎ waiting for Kustomization reconciliation
✔ applied revision main/8ad0d7b302be3fd253568ad84354a79385c98908
```

```bash
flux reconcile [command]
```

where command can be: `alert`, `alert-provider`, `helmrelease`, `image`, `kustomization`, `receiver`, `source`.
        
## Watch Kustomize Controller reconciliation events

```
flux get kustomizations --watch 
```

### Examples

#### Reconcile a specific kustomization

1. Get avaible `kustomizations`

```bash
$ flux get kustomizations

live-infra-secrets	main/ff3d806  	False    	True 	Applied revision: main/ff3d806  	
flux-system       	master/eb6ccff	False    	True 	Applied revision: master/eb6ccff	
crds              	master/eb6ccff	False    	True 	Applied revision: master/eb6ccff	
apps              	master/eb6ccff	False    	True 	Applied revision: master/eb6ccff	
infrastructure    	master/eb6ccff	False    	True 	Applied revision: master/eb6ccff	
```

If you add a new top-level Kustomization (i.e. as a peer to infrastructure, sources, or crds), telling Flux to reconcile flux-system will pull it into the active configuration, at which point Flux will also reconcile it and all the other Kustomizations. Once it shows up in the output of flux get kustomizations, you can reconcile it independently.

### Processing a specific kustomization

```bash
flux reconcile kustomization "apps"
```

#### Check Git Repository Hashes

1. Get all the avaible sources

```bash
flux get sources all
flux get sources git
```

2. Get the last commit for a git repo

```bash
flux reconcile source git live-infra-secrets
```

## Encrypt secrets with sops and age

### Configure age to encrypt secrets with sops

Generate an age key with `age` using `age-keygen`

```
$ age-keygen -o age.agekey
Public key: age1helqcqsh9464r8chnwc2fzj8uv7vr5ntnsft0tn45v2xtz0hpfwq98cmsg
```


```
$ ls -1
...
age.agekey
...
$ cat age.agekey
# created: ...
# public key: age1helqcqsh9464r8chnwc2fzj8uv7vr5ntnsft0tn45v2xtz0hpfwq98cmsg
AGE-SECRET-KEY-...
```

The public key is used to encrypt and it's okay to share, while the secret key decrypts data and must be kept private.

We can specify via configuration which keys to encrypt. Create a `.sops.yaml` file
in the root directory of the flux diretory

```
creation_rules:
  - encrypted_regex: '^(data|stringData)$'
    age: age1helqcqsh9464r8chnwc2fzj8uv7vr5ntnsft0tn45v2xtz0hpfwq98cmsg
```

Notice that we added the age public key in the YAML file. 

### Encrypt secrets with sops and age

Use `sops` and the age public key to encrypt a Kubernetes secret:

1. Create generic plain kubernetes secret

```
kubectl create secret generic sopstest --from-literal=foo=bar -o yaml \
    --dry-run=client | tee sops-test-secret.yaml
```

2. Encrypt previous kubernetes secret with sops and the age public key

```bash
sops --age=age1helqcqsh9464r8chnwc2fzj8uv7vr5ntnsft0tn45v2xtz0hpfwq98cmsg \ -e sops-test-secret.yaml | tee sops-test-secret-encrypted.yaml
```

The secret `sops-test-secret-encrypted.yaml` can be commited to the Git repository.

Finally, create a secret with the age private key, the key name must end with `.agekey` to be detected as an age key

```
cat age.agekey |
kubectl create secret generic sops-age \
--namespace=flux-system \
--from-file=age.agekey=/dev/stdin
```

### Dencrypt sops encrypted secrets using age

Sops will lookup for a data key to decrypt the sops file. The file is located at 

```
$HOME/.age-keys/age-key.txt
```

```
➜  live-infra-secrets git:(main) ls -ltrha
total 32K
-rw-r--r--   1 oscar oscar   21 ago  7 16:50 README.md
-rw-r--r--   1 oscar oscar  115 ago  7 18:01 kustomization.yaml
drwxr-xr-x   3 oscar oscar 4,0K ago  7 18:18 .
-rw-r--r--   1 oscar oscar 1,3K ago  7 18:27 sops-test-secret-encrypted.yaml
drwxr-xr-x   8 oscar oscar 4,0K ago  7 18:27 .git
drwxr-xr-x 128 oscar oscar  12K ago 21 12:34 ..
```

```bash
➜  live-infra-secrets git:(main) sops -d sops-test-secret-encrypted.yaml
apiVersion: v1
data:
    foo: YmFy
kind: Secret
metadata:
    creationTimestamp: null
    name: sopstest
    namespace: monitoring
➜  live-infra-secrets git:(main)
```

### Gitops workflow

A cluster admin should create the Kubernetes secret with the PGP keys on each cluster and add the GitRepository/Kustomization manifests to the fleet repository.

Git repository manifest:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: GitRepository
metadata:
  name: live-infra-secrets
  namespace: flux-system
spec:
  interval: 1m
  url: https://github.com/my-org/my-secrets
```

Kustomization manifest:

```yaml
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  name: live-infra-secrets
  namespace: flux-system
spec:
  interval: 10m0s
  sourceRef:
    kind: GitRepository
    name: live-infra-secrets
  path: ./
  prune: true
  decryption:
    provider: sops
    secretRef:
      name: sops-age
```