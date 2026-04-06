#!/bin/sh
set -e

# ─────────────────────────────────────────────
# CONFIG
# ─────────────────────────────────────────────
ZITADEL_INTERNAL_URL="${ZITADEL_INTERNAL_URL:-http://zitadel-api:8080}"
ZITADEL_HOST="${ZITADEL_HOST:-zitadel.localhost:8080}"
PAT_FILE="${PAT_FILE:-/zitadel/bootstrap/login-client.pat}"
SECRETS_DIR="${SECRETS_DIR:-/secrets}"

PROJECT_NAME="yarmaq"
WEBHOOK_URL="${WEBHOOK_URL:-http://yarmaq-accounts-service:8081/api/v1/webhook/sync}"

ROLES="support admin"

# username_в_zitadel:имя_файла_в_secrets
SERVICE_USERS="accounts-service:accounts transaction-service:transaction notification-service:notification"

# ─────────────────────────────────────────────
# HELPERS
# ─────────────────────────────────────────────
log()  { echo "[bootstrap] $*"; }
ok()   { echo "[bootstrap] ✓ $*"; }
skip() { echo "[bootstrap] → skip: $*"; }
fail() { echo "[bootstrap] ✗ $*" >&2; exit 1; }

zapi() {
  METHOD="$1"; shift
  ZPATH="$1"; shift
  curl -s \
    -X "$METHOD" \
    -H "Host: $ZITADEL_HOST" \
    -H "Authorization: Bearer $PAT" \
    -H "Content-Type: application/json" \
    "$ZITADEL_INTERNAL_URL$ZPATH" \
    "$@"
}
# ─────────────────────────────────────────────
# WAIT FOR PAT FILE
# ─────────────────────────────────────────────
log "waiting for PAT file at $PAT_FILE ..."
ATTEMPTS=0
until [ -f "$PAT_FILE" ] && [ -s "$PAT_FILE" ]; do
  ATTEMPTS=$((ATTEMPTS + 1))
  if [ "$ATTEMPTS" -ge 30 ]; then
    fail "PAT file not found after 30 attempts. Is Zitadel running?"
  fi
  sleep 2
done

PAT=$(cat "$PAT_FILE")
ok "PAT loaded"

# ─────────────────────────────────────────────
# WAIT FOR ZITADEL API
# ─────────────────────────────────────────────
log "waiting for Zitadel API..."
ATTEMPTS=0
until zapi GET "/healthz" > /dev/null 2>&1; do
  ATTEMPTS=$((ATTEMPTS + 1))
  if [ "$ATTEMPTS" -ge 30 ]; then
    fail "Zitadel API not ready after 30 attempts"
  fi
  sleep 2
done
ok "Zitadel API is ready"

# ─────────────────────────────────────────────
# PROJECT
# ─────────────────────────────────────────────
log "checking project '$PROJECT_NAME'..."

PROJECTS_RESP=$(zapi POST "/management/v1/projects/_search" \
  -d "{\"queries\":[{\"nameQuery\":{\"name\":\"$PROJECT_NAME\",\"method\":\"TEXT_QUERY_METHOD_EQUALS\"}}]}")

PROJECT_ID=$(echo "$PROJECTS_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$PROJECT_ID" ]; then
  skip "project '$PROJECT_NAME' already exists (id=$PROJECT_ID)"
else
  log "creating project '$PROJECT_NAME'..."
  PROJECT_RESP=$(zapi POST "/management/v1/projects" \
    -d "{\"name\":\"$PROJECT_NAME\",\"projectRoleAssertion\":true,\"projectRoleCheck\":true}")
  PROJECT_ID=$(echo "$PROJECT_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
  [ -n "$PROJECT_ID" ] || fail "failed to create project"
  ok "project created (id=$PROJECT_ID)"
fi

# ─────────────────────────────────────────────
# ROLES
# ─────────────────────────────────────────────
log "checking roles..."

EXISTING_ROLES=$(zapi POST "/management/v1/projects/$PROJECT_ID/roles/_search" \
  -d '{}' | grep -o '"roleKey":"[^"]*"' | cut -d'"' -f4)

for ROLE in $ROLES; do
  if echo "$EXISTING_ROLES" | grep -qx "$ROLE"; then
    skip "role '$ROLE' already exists"
  else
    zapi POST "/management/v1/projects/$PROJECT_ID/roles" \
      -d "{\"roleKey\":\"$ROLE\",\"displayName\":\"$ROLE\"}" > /dev/null
    ok "role '$ROLE' created"
  fi
done

# ─────────────────────────────────────────────
# SERVICE USERS
# ─────────────────────────────────────────────
log "checking service users..."

for ENTRY in $SERVICE_USERS; do
  SVC_USERNAME=$(echo "$ENTRY" | cut -d':' -f1)   # accounts-service
  SVC_KEYNAME=$(echo "$ENTRY" | cut -d':' -f2)    # accounts
  KEY_FILE="$SECRETS_DIR/$SVC_KEYNAME.json"

  # Если ключ уже есть — пропускаем
  if [ -f "$KEY_FILE" ] && [ -s "$KEY_FILE" ]; then
    skip "key for '$SVC_USERNAME' already exists at $KEY_FILE"
    continue
  fi

  # Ищем существующего machine user
  USERS_RESP=$(zapi POST "/management/v1/users/_search" \
    -d "{\"queries\":[{\"userNameQuery\":{\"userName\":\"$SVC_USERNAME\",\"method\":\"TEXT_QUERY_METHOD_EQUALS\"}}]}")

  USER_ID=$(echo "$USERS_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

  if [ -n "$USER_ID" ]; then
    skip "service user '$SVC_USERNAME' already exists (id=$USER_ID)"
  else
    log "creating service user '$SVC_USERNAME'..."
    USER_RESP=$(zapi POST "/management/v1/users/machine" \
      -d "{\"userName\":\"$SVC_USERNAME\",\"name\":\"$SVC_USERNAME\",\"description\":\"Service account for $SVC_USERNAME\",\"accessTokenType\":\"ACCESS_TOKEN_TYPE_JWT\"}")
    USER_ID=$(echo "$USER_RESP" | grep -o '"userId":"[^"]*"' | cut -d'"' -f4)
    [ -n "$USER_ID" ] || fail "failed to create service user '$SVC_USERNAME'"
    ok "service user '$SVC_USERNAME' created (id=$USER_ID)"
  fi

  # Генерируем ключ
  log "generating key for '$SVC_USERNAME' -> $KEY_FILE ..."
  KEY_RESP=$(zapi POST "/management/v1/users/$USER_ID/keys" \
    -d '{"type":"KEY_TYPE_JSON","expirationDate":"2099-01-01T00:00:00Z"}')

  KEY_DETAIL=$(echo "$KEY_RESP" | grep -o '"keyDetails":"[^"]*"' | cut -d'"' -f4)
  [ -n "$KEY_DETAIL" ] || fail "failed to generate key for '$SVC_USERNAME'"

  echo "$KEY_DETAIL" | base64 -d > "$KEY_FILE"
  ok "key saved to $KEY_FILE"
done

# ─────────────────────────────────────────────
# WEBHOOK TARGET
# ─────────────────────────────────────────────
log "checking webhook target..."

TARGET_ID=""

TARGET_CREATE_RESP=$(zapi POST "/v2/actions/targets" \
  -d "{\"name\":\"yarmaq-webhook\",\"restWebhook\":{\"interruptOnError\":false},\"endpoint\":\"$WEBHOOK_URL\",\"timeout\":\"10s\"}")

TARGET_ID=$(echo "$TARGET_CREATE_RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$TARGET_ID" ]; then
  ok "webhook target created (id=$TARGET_ID)"
else
  skip "webhook target already exists or failed to create, trying to find id..."
  # fallback: достать из UI невозможно через API, поэтому падаем
  fail "could not get target id — delete existing target manually and re-run"
fi

# ─────────────────────────────────────────────
# WEBHOOK EXECUTION
# ─────────────────────────────────────────────
log "setting execution for AddHumanUser response -> target $TARGET_ID..."
zapi POST "/zitadel.action.v2.ActionService/SetExecution" \
  -d "{
    \"condition\": {
      \"response\": {
        \"method\": \"/zitadel.user.v2.UserService/AddHumanUser\"
      }
    },
    \"targets\": [\"$TARGET_ID\"]
  }" > /dev/null

ok "execution set"

# ─────────────────────────────────────────────
# DONE
# ─────────────────────────────────────────────
echo ""
echo "╔══════════════════════════════════════╗"
echo "║   Bootstrap completed successfully   ║"
echo "╚══════════════════════════════════════╝"
echo ""
echo "  Project:  $PROJECT_NAME (id=$PROJECT_ID)"
echo "  Roles:    $ROLES"
echo "  Secrets:  $SECRETS_DIR/"
echo ""