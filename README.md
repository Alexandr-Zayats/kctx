# kctx — Kubernetes Context Switcher for Multi-Cloud (Enterprise Ready)

`kctx` is a CLI tool for fast and secure switching between Kubernetes contexts in a multi-cloud environment:

* AWS (SSO, multi-account)
* GCP (gcloud)
* DigitalOcean (teams)

---

## 🚀 Features

* 🔐 AWS SSO support (no static credentials)
* 🧠 Automatic generation of `~/.aws/config`
* 🌍 Multi-account support (dev / stg / prod / root / services)
* ☁️ Multi-cloud: AWS, GCP, DO
* ⚡ Fast switching between clusters and namespaces
* 🧩 Project-based configuration (`.kctx/aws.yaml`)
* 🛡 Secure: no access keys stored in code

---

## 📦 Installation

```bash
git clone <repo>
cd kctx
go build -o kctx
cp kctx /usr/local/bin/
```

---

## ⚙️ Configuration (AWS — primary use case)

### 📁 `.kctx/aws.yaml`

Create this file in your project root:

```yaml
project: bethunder

sso:
  start_url: https://your-org.awsapps.com/start
  region: eu-central-1
  role_name: AdministratorAccess

accounts:
  dev:
    account_id: "123456789012"
    region: eu-central-1

  stg:
    account_id: "234567890123"
    region: eu-central-1

  prod:
    account_id: "345678901234"
    region: eu-central-1

  root:
    account_id: "456789012345"
    region: eu-central-1

  services:
    account_id: "567890123456"
    region: eu-central-1
```

---

## 🧠 How it works

When running:

```bash
kctx switch aws
```

kctx performs:

1. 📂 Detects `.kctx/aws.yaml`
2. 🛠 Generates / updates `~/.aws/config`
3. 🔐 Validates SSO session:

   * if valid → continue
   * if expired → runs `aws sso login`
4. 📡 Fetches EKS clusters
5. 🎯 Prompts selection:

   * account
   * cluster
   * namespace
6. 🔄 Updates kubeconfig

---

## 🔐 AWS SSO (Enterprise Best Practice)

kctx relies on:

```bash
aws sso login
aws sts get-caller-identity
```

### Requirements

* AWS Identity Center (SSO) configured
* User has access to accounts
* Permission set assigned (e.g. `AdministratorAccess`)

---

## ☁️ Usage

### AWS

```bash
kctx switch aws
```

---

### GCP

```bash
kctx switch gcp
```

Requirements:

```bash
gcloud auth login
gcloud config set project <project>
```

---

### DigitalOcean

```bash
kctx switch do
```

Requirements:

```bash
doctl auth init
```

---

## ⚡ Quick Switch

```bash
kctx switch
```

Opens interactive selection:

* provider
* account
* cluster

---

## 🧩 Aliases

```yaml
aliases:
  dev: aws-bethunder-dev-cluster
```

```bash
kctx switch dev
```

---

## 🔧 Environment Variables

kctx automatically manages:

```bash
AWS_PROFILE
AWS_DEFAULT_PROFILE
AWS_SDK_LOAD_CONFIG
```

---

## 🛡 Security

* ❌ No access keys stored
* ❌ No credentials persisted
* ✅ Uses AWS SSO
* ✅ AWS CLI is the source of truth

---

## 🧱 Enterprise Recommendations

### 1. Project Structure

```
repo/
 ├── .kctx/
 │    └── aws.yaml
 ├── services/
 ├── infra/
```

---

### 2. GitOps Integration

Use together with:

* Terraform / Terragrunt
* FluxCD
* Helm

---

### 3. RBAC Strategy

Use different roles:

```yaml
role_name: ReadOnlyAccess
role_name: PowerUserAccess
role_name: AdministratorAccess
```

---

### 4. Multi-Team Setup

Separate projects:

```yaml
project: payments
project: platform
project: analytics
```

---

## 🧪 Troubleshooting

### Check AWS identity

```bash
aws sts get-caller-identity
```

---

### Check clusters

```bash
aws eks list-clusters --profile <profile>
```

---

### Reset config

```bash
rm ~/.aws/config
```

---

## 🧠 Architecture

```
kctx
 ├── provider/
 │    ├── aws
 │    ├── gcp
 │    └── do
 ├── core/
 ├── ui/
 └── kube/
```

---

## 🚀 Roadmap

* [ ] Auto-detect SSO roles
* [ ] Cluster caching improvements
* [ ] Parallel multi-cloud queries
* [ ] UI improvements (fzf previews)
* [ ] Plugin system

---

## 📌 TL;DR

```bash
kctx switch aws
```

Everything else is handled automatically.

---

## 🧑‍💻 Author
Alexander Zaiats alex@zaiats.de
DevOps / Platform Engineering tooling for production environments

