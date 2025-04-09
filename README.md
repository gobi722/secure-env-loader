# Secure Env Loader 🔐

A lightweight Go utility to securely fetch and load environment variables from an encrypted database, designed for backend services that require dynamic configuration without exposing secrets in code or `.env` files.

## 🌟 Features

- Secure storage of DB credentials and API keys in MongoDB (or any DB)
- Encrypted secret key-based access
- Auto-fetch and injects env variables during application start
- Works well with microservices or distributed systems
- Minimalistic and easy to integrate

## 📦 Technologies Used

- Go (Golang)
- MongoDB
- AES Encryption (CTR or ECB mode)
- dotenv (if `.env` fallback is needed)

## 🔐 Use Case

You can't store DB credentials directly in `.env` files or hardcode them.  
Instead, store them in a DB collection in encrypted format and access them securely at runtime.

## 🧠 How It Works

1. On app start, pass a **decryption key** as an input flag or env variable.
2. The app connects to a base DB using minimal details.
3. It decrypts the secrets stored in a collection (like `secrets_config`) using the key.
4. The fetched values are injected into `os.Environ()` so that all parts of the app can use them like normal env variables.

## 🧪 Example

```go
os.Getenv("DB_HOST") // returns 127.0.0.1 from secure loader
