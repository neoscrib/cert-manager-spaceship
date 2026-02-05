# 🚀 cert-manager DNS-01 Webhook for Spaceship

[![cert-manager](https://img.shields.io/badge/cert--manager-webhook-blue)](https://cert-manager.io/docs/configuration/acme/dns01/)
[![Go Report Card](https://goreportcard.com/badge/github.com/neoscrib/cert-manager-spaceship)](https://goreportcard.com/report/github.com/neoscrib/cert-manager-spaceship)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/neoscrib/cert-manager-webhook-spaceship)](https://hub.docker.com/r/neoscrib/cert-manager-webhook-spaceship)
[![Docker Image Version](https://img.shields.io/docker/v/neoscrib/cert-manager-webhook-spaceship?sort=semver)](https://hub.docker.com/r/neoscrib/cert-manager-webhook-spaceship)

This repository provides a **cert-manager DNS-01 webhook solver** for **Spaceship DNS** (https://spaceship.com).

It allows cert-manager to solve ACME DNS-01 challenges by creating and deleting
`_acme-challenge` TXT records using the Spaceship DNS API.

> This is an independently developed webhook and is not affiliated with or endorsed by the cert-manager project.

---

## ✨ Features

- Native **cert-manager webhook** implementation (gRPC + Protocol Buffers)
- Written in **Go**, no runtime dependencies
- Supports **Let's Encrypt** and any ACME-compatible CA
- Packaged for Kubernetes with **Helm**

---

## 📦 Installation (Helm)

### 1. Add the Helm repository

```bash
helm repo add spaceship-webhook https://neoscrib.github.io/cert-manager-spaceship
helm repo update
```

### 2. Install the webhook

```bash
helm install spaceship-webhook spaceship-webhook/cert-manager-webhook-spaceship \
  --namespace cert-manager \
  --create-namespace
```

---

## 📦 Installation (Rendered Manifest)

If you are not using Helm, you can apply the pre-rendered manifest. Create the
Spaceship API Secret first (the manifest does not include your credentials):

```bash
kubectl create secret generic spaceship-api-key \
  --namespace cert-manager \
  --from-literal=api-key="$SPACESHIP_API_KEY" \
  --from-literal=api-secret="$SPACESHIP_API_SECRET"
```

Then apply the rendered manifest:

```bash
kubectl apply -f https://raw.githubusercontent.com/neoscrib/cert-manager-spaceship/main/deploy/manifest.yaml
```

> Note: the rendered manifest is generated from the chart's default values and
> targets the `cert-manager` namespace.

---

## 🔐 Spaceship API Key

This webhook requires a Spaceship API key and secret. Login to [Spaceship API Manager](https://www.spaceship.com/application/api-manager/) to create an API key.

### Option A: Use an existing Secret (recommended)

Create the Secret manually or via a secret manager:

```bash
kubectl create secret generic spaceship-api-key \
  --namespace cert-manager \
  --from-literal=api-key="$SPACESHIP_API_KEY" \
  --from-literal=api-secret="$SPACESHIP_API_SECRET" \
```

Then install the chart:

```bash
helm install spaceship-webhook spaceship-webhook/cert-manager-webhook-spaceship \
  --namespace cert-manager \
  --set secrets.secretKeyRef.name=spaceship-api-key \
  --set secrets.secretKeyRef.apiKeyKey=api-key \
  --set secrets.secretKeyRef.apiSecretKey=api-secret
```

### Option B: Let the chart create the Secret

```bash
helm install spaceship-webhook spaceship-webhook/cert-manager-webhook-spaceship \
  --namespace cert-manager \
  --set secrets.createSecret=true \
  --set secrets.apiKey="$SPACESHIP_API_KEY" \
  --set secrets.apiSecret="$SPACESHIP_API_SECRET"
```

---

## 📄 Example ClusterIssuer

> You can use a namespaced `Issuer` instead of a `ClusterIssuer` by changing `kind: Issuer`
> and creating it in the target namespace. The webhook works with both.

Production

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod-spaceship
spec:
  acme:
    email: admin@example.com # replace with your email address
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod-spaceship-key
    solvers:
      - dns:
          webhook:
            groupName: acme.spaceship.neoscrib.com
            solverName: spaceship
```

Staging

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging-spaceship
spec:
  acme:
    email: admin@example.com # replace with your email address
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-staging-spaceship-key
    solvers:
      - dns:
          webhook:
            groupName: acme.spaceship.neoscrib.com
            solverName: spaceship
```

---

## 📄 Example Certificate

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-cert
  namespace: default
spec:
  secretName: example-cert-tls
  dnsNames:
    - example.com
    - "*.example.com"
  issuerRef:
    name: letsencrypt-prod-spaceship
    kind: ClusterIssuer
```

---

## 📘 Configuration Values

```yaml
groupName: acme.spaceship.neoscrib.com

secrets:
  createSecret: false
  apiKey: ""
  apiSecret: ""
  secretKeyRef:
    name: spaceship-api-key
    apiKeyKey: api-key
    apiSecretKey: api-secret
```

---

## ✅ Tested Versions

- Kubernetes: k3s v1.34.1+k3s1
- cert-manager: v1.19.2
- webhook image: v0.1.0

---

## 🧭 Supported Versions (cert-manager)

This webhook follows cert-manager's supported Kubernetes version matrix. Refer to cert-manager's "Supported Releases" table for the current support window.

Current supported cert-manager releases (per cert-manager docs):
- cert-manager 1.19: Kubernetes 1.31 → 1.35
- cert-manager 1.18: Kubernetes 1.29 → 1.33

Upcoming release:
- cert-manager 1.20 (upcoming): Kubernetes 1.32 → 1.35

Note: support windows change over time; always defer to the official cert-manager docs for the latest matrix.

---

## 🤝 Contributing

Contributions are welcome!

Ideas:
- Additional Spaceship DNS features
- Improved error handling
- Multi-zone support
- CI pipelines for multi-arch images

---

## 📜 License

MIT License
