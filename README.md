# 🚀 cert-manager DNS-01 Webhook for Spaceship

[![cert-manager](https://img.shields.io/badge/cert--manager-webhook-blue)](https://cert-manager.io/docs/configuration/acme/dns01/)

This repository provides a **cert-manager DNS-01 webhook solver** for **Spaceship DNS** (https://spaceship.com).

It allows cert-manager to solve ACME DNS-01 challenges by creating and deleting
`_acme-challenge` TXT records using the Spaceship DNS API.

---

## ✨ Features

- Native **cert-manager webhook** implementation (gRPC + Protocol Buffers)
- Written in **Go**, no runtime dependencies
- Supports **Let’s Encrypt** and any ACME-compatible CA
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
