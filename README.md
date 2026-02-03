# 🚀 cert-manager DNS-01 Webhook for Spaceship

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

## 🧠 How It Works

cert-manager communicates with this webhook using **gRPC (HTTP/2 + Protobuf)**.

Flow:

1. cert-manager calls `Present()` to request a TXT record
2. Webhook calls the Spaceship DNS API to create the record
3. cert-manager verifies the challenge via DNS
4. cert-manager calls `CleanUp()` to remove the TXT record

No external HTTP endpoints are exposed — everything runs inside the cluster.

---

## 📦 Installation (Helm)

### 1. Add the Helm repository

```bash
helm repo add spaceship-webhook https://neoscrib.github.io/cert-manager-spaceship
helm repo update
```

### 2. Install the webhook

```bash
helm install spaceship-webhook spaceship-webhook/cert-manager-spaceship \
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

This webhook requires a Spaceship API key and secret.

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
helm install spaceship-webhook spaceship-webhook/cert-manager-spaceship \
  --namespace cert-manager \
  --set secrets.secretKeyRef.name=spaceship-api-key \
  --set secrets.secretKeyRef.apiKeyKey=api-key \
  --set secrets.secretKeyRef.apiSecretKey=api-secret
```

### Option B: Let the chart create the Secret

```bash
helm install spaceship-webhook spaceship-webhook/cert-manager-spaceship \
  --namespace cert-manager \
  --set secrets.createSecret=true \
  --set secrets.apiKey="$SPACESHIP_API_KEY" \
  --set secrets.apiSecret="$SPACESHIP_API_SECRET"
```

---

## 📄 Example ClusterIssuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-spaceship
spec:
  acme:
    email: admin@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-spaceship-key
    solvers:
      - dns:
          webhook:
            groupName: acme.spaceship.neoscrib.com
            solverName: spaceship
```

> No service name or URL is required — the webhook is discovered via
> the aggregated `APIService` registered by this chart.

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

## 🪵 Logging & Debugging

```bash
kubectl logs -l app=cert-manager-spaceship -n cert-manager
```

Look for `Present` and `CleanUp` log entries.

---

## 🧩 Why gRPC?

cert-manager webhooks use **gRPC** to provide:

- Strongly typed contracts
- Efficient communication
- Language-agnostic implementations

While this is heavier than HTTP+JSON, it ensures consistency across solvers.

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
