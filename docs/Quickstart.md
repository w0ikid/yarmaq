# Quickstart

## 1. Setup

```bash
cp .env.example .env
make zitadel-up
```

Wait until Zitadel is healthy, then open [http://zitadel.localhost:8080/ui/console](http://zitadel.localhost:8080/ui/console).

---

## 2. Login to Zitadel Console

```
Email:    yarmaq@gmail.com
Password: RootPassword1!
```

---

## 3. Grant IAM Owner to Login Client

1. Go to **Users → Service Users**
2. Click on **"Automatically Initialized IAM_LOGIN_CLIENT"**
3. Open **Memberships** tab
4. Click **+** → select role **IAM Owner**
5. Save

---

## 4. Run Bootstrap

```bash
make zitadel-bootstrap
```

This will create the project, service users, and generate JWT keys in `secrets/`.

---

## 5. Grant Project Owner to Service Users

1. Go to **Users → Service Users**
2. For each service user (`accounts`, `transaction`, `notification`):
   - Click on the user
   - Open **Memberships** tab
   - Click **+** → select your project → role **Project Owner**
   - Save

---

## 6. Create Application

1. Go to your **Project**
2. Click **+ New Application**
3. Name it (e.g. `yarmaq-app`)
4. Type: **User Agent** → Auth Method: **PKCE**
5. Redirect URI: `https://oauth.pstmn.io/v1/callback`
6. Click through and save

Then go to application settings and enable:
- ✅ **Assert Roles on Authentication**

Copy the **Client ID**.

---

## 7. Configure Postman

In Postman → your collection → **Variables**:

```
clientId = <your Client ID from step 6>
```

Authorization tab:
- Type: **OAuth 2.0**
- Grant Type: **Authorization Code (PKCE)**
- Auth URL: `http://zitadel.localhost:8080/oauth/v2/authorize`
- Token URL: `http://zitadel.localhost:8080/oauth/v2/token`
- Scope: `openid profile email`

---

## 8. Start Services

```bash
make infra-up
make apps-up
```

That's it - everything is running.