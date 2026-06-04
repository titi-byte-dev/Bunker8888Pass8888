<script lang="ts">
  /**
   * Playground de desenvolvimento (INFRA-006+) — UI mínima para testar o cofre.
   * Não é o produto final; serve para validar lib + API end-to-end.
   */
  import { loginUser, loginAfterRegister, registerUser, fetchKdfParams, deriveMasterKeyBytes } from "./lib/auth";
  import { saveSessionToken, loadSessionToken, clearSessionToken } from "./lib/session";
  import { VaultAPI } from "./lib/vault/api";
  import { blobFromBase64, openItem, sealItem, blobToBase64 } from "./lib/vault/items";
  import { setMasterKey, purgeMasterKey, getMasterKey } from "./lib/vault/masterKeyStore";
  import { generatePassword } from "./lib/vault/password";
  import { importCsvToVaultInputs, uploadImportItems } from "./lib/vault/import";
  import {
    generateRecoveryCode,
    wrapMasterKeyBytes,
    uploadRecoveryBackup,
    fetchRecoveryBackupStatus,
    recoverMasterKeyFromEmail,
  } from "./lib/vault/recovery";
  import {
    passkeysSupported,
    registerPasskey,
    unlockWithPasskeyAndPassword,
    listPasskeys,
    type PasskeyMeta,
  } from "./lib/passkey";
  import type { LoginItem, VaultItemMeta } from "./lib/vault/types";

  const API = ""; // Vite proxy → localhost:8080

  let screen = $state<"login" | "recover">("login");
  let email = $state("");
  let masterPassword = $state("");
  let token = $state<string | null>(loadSessionToken());
  let status = $state("");
  let busy = $state(false);

  let items = $state<Array<{ meta: VaultItemMeta; login?: LoginItem }>>([]);

  let newTitle = $state("");
  let newUser = $state("");
  let newPass = $state("");
  let importPreview = $state("");

  let recoveryConfigured = $state(false);
  let recoveryConfirmPw = $state("");
  let newRecoveryCode = $state("");
  let recoverCode = $state("");

  let passkeyName = $state("");
  let passkeys = $state<PasskeyMeta[]>([]);
  const webAuthnOk = passkeysSupported();

  async function handleRegister() {
    busy = true;
    status = "A derivar chaves (Argon2id)…";
    try {
      await registerUser(API, email, masterPassword);
      status = "Registado — a iniciar sessão…";
      const { masterKey, token: t } = await loginAfterRegister(API, email, masterPassword);
      token = t;
      saveSessionToken(t);
      setMasterKey(masterKey);
      status = "Sessão iniciada.";
      recoveryConfigured = await fetchRecoveryBackupStatus(API, t);
      await refreshVault();
    } catch (e) {
      status = e instanceof Error ? e.message : "Erro no registo";
    } finally {
      busy = false;
    }
  }

  async function handleLogin() {
    busy = true;
    status = "A derivar chaves…";
    try {
      const { masterKey, token: t } = await loginUser(API, email, masterPassword);
      token = t;
      saveSessionToken(t);
      setMasterKey(masterKey);
      status = "Sessão iniciada.";
      recoveryConfigured = await fetchRecoveryBackupStatus(API, t);
      passkeys = await listPasskeys(API, t).catch(() => []);
      await refreshVault();
    } catch (e) {
      status = e instanceof Error ? e.message : "Erro no login";
    } finally {
      busy = false;
    }
  }

  async function handlePasskeyLogin() {
    if (!email || !masterPassword) {
      status = "Email e Master Password são necessários (passkey autentica o servidor; password desbloqueia o cofre ZK).";
      return;
    }
    busy = true;
    status = "A autenticar com passkey…";
    try {
      const { masterKey, token: t } = await unlockWithPasskeyAndPassword(API, email, masterPassword);
      token = t;
      saveSessionToken(t);
      setMasterKey(masterKey);
      status = "Sessão iniciada via passkey.";
      recoveryConfigured = await fetchRecoveryBackupStatus(API, t);
      passkeys = await listPasskeys(API, t).catch(() => []);
      await refreshVault();
    } catch (e) {
      status = e instanceof Error ? e.message : "Passkey falhou";
    } finally {
      busy = false;
    }
  }

  async function handleRegisterPasskey() {
    if (!token || !passkeyName) return;
    busy = true;
    try {
      await registerPasskey(API, token, passkeyName);
      passkeyName = "";
      passkeys = await listPasskeys(API, token);
      status = "Passkey registada.";
    } catch (e) {
      status = e instanceof Error ? e.message : "Registo passkey falhou";
    } finally {
      busy = false;
    }
  }

  function handleLogout() {
    token = null;
    clearSessionToken();
    purgeMasterKey();
    items = [];
    status = "Sessão terminada — Master Key expurgada.";
  }

  async function refreshVault() {
    const key = getMasterKey();
    if (!token || !key) return;
    const api = new VaultAPI(API, token);
    const metas = await api.list("login");
    const decoded: typeof items = [];
    for (const meta of metas) {
      let login: LoginItem | undefined;
      if (meta.blob) {
        login = (await openItem(key, blobFromBase64(meta.blob))) as LoginItem;
      }
      decoded.push({ meta, login });
    }
    items = decoded;
  }

  async function handleAddLogin() {
    const key = getMasterKey();
    if (!token || !key) return;
    busy = true;
    try {
      const payload: LoginItem = {
        kind: "login",
        title: newTitle,
        username: newUser,
        password: newPass,
      };
      const blob = blobToBase64(await sealItem(key, payload));
      const api = new VaultAPI(API, token);
      await api.create({ type: "login", blob });
      newTitle = newUser = newPass = "";
      await refreshVault();
      status = "Login guardado (cifrado).";
    } catch (e) {
      status = e instanceof Error ? e.message : "Erro ao guardar";
    } finally {
      busy = false;
    }
  }

  async function handleImport(ev: Event) {
    const key = getMasterKey();
    if (!token || !key) return;
    const file = (ev.target as HTMLInputElement).files?.[0];
    if (!file) return;
    busy = true;
    try {
      const text = await file.text();
      const { parse, inputs } = await importCsvToVaultInputs(key, text);
      importPreview = `${parse.rows.length} logins (${parse.format}), ${parse.skipped} ignorados`;
      const api = new VaultAPI(API, token);
      const result = await uploadImportItems((i) => api.create(i), inputs);
      status = `Importação: ${result.created} criados, ${result.failed} falharam.`;
      await refreshVault();
    } catch (e) {
      status = e instanceof Error ? e.message : "Importação falhou";
    } finally {
      busy = false;
    }
  }

  function fillGeneratedPassword() {
    newPass = generatePassword({ length: 20 });
  }

  async function handleCreateRecovery() {
    if (!token) return;
    busy = true;
    try {
      const kdf = await fetchKdfParams(API, email);
      const mk = await deriveMasterKeyBytes(recoveryConfirmPw, kdf.salt, kdf);
      const code = generateRecoveryCode();
      const envelope = await wrapMasterKeyBytes(mk, code);
      await uploadRecoveryBackup(API, token, envelope);
      newRecoveryCode = code;
      recoveryConfigured = true;
      recoveryConfirmPw = "";
      status = "Chave de recuperação criada — guarda o código abaixo offline.";
    } catch (e) {
      status = e instanceof Error ? e.message : "Falha ao criar recuperação";
    } finally {
      busy = false;
    }
  }

  async function handleRecoverMasterKey() {
    busy = true;
    try {
      await recoverMasterKeyFromEmail(API, email, recoverCode);
      status =
        "Master Key recuperada com sucesso. Ainda precisas de redefinir a Master Password (em breve) para voltar a entrar.";
    } catch (e) {
      status = e instanceof Error ? e.message : "Recuperação falhou";
    } finally {
      busy = false;
    }
  }
</script>

<main>
  <header>
    <h1>AegisPass</h1>
    <p class="tag">Playground de desenvolvimento — Zero-Knowledge no browser</p>
  </header>

  {#if !token}
    <div class="row" style="margin-bottom: 1rem">
      <button class:secondary={screen !== "login"} onclick={() => (screen = "login")}>Entrar</button>
      <button class:secondary={screen !== "recover"} onclick={() => (screen = "recover")}>Recuperar</button>
    </div>

    {#if screen === "login"}
    <section class="card">
      <h2>Iniciar sessão</h2>
      <label>
        Email
        <input type="email" bind:value={email} autocomplete="username" />
      </label>
      <label>
        Master Password
        <input type="password" bind:value={masterPassword} autocomplete="current-password" />
      </label>
      <div class="row">
        <button onclick={handleLogin} disabled={busy || !email || !masterPassword}>Entrar</button>
        <button class="secondary" onclick={handleRegister} disabled={busy || !email || !masterPassword}>
          Registar
        </button>
        {#if webAuthnOk}
          <button class="secondary" onclick={handlePasskeyLogin} disabled={busy || !email || !masterPassword}>
            Passkey
          </button>
        {/if}
      </div>
      {#if webAuthnOk}
        <p class="muted">Passkey autentica o servidor; a Master Password continua necessária para desbloquear o cofre (Zero-Knowledge).</p>
      {/if}
    </section>
    {:else}
    <section class="card">
      <h2>Recuperar Master Key</h2>
      <p class="muted">Usa a chave de recuperação guardada offline. Não substitui login — valida o backup cifrado.</p>
      <label>Email <input type="email" bind:value={email} /></label>
      <label>Código de recuperação <input type="text" bind:value={recoverCode} placeholder="XXXXX-XXXXX-XXXXX-XXXXX" /></label>
      <button onclick={handleRecoverMasterKey} disabled={busy || !email || !recoverCode}>Recuperar</button>
    </section>
    {/if}
  {:else}
    <section class="card">
      <div class="row spread">
        <span class="muted">Sessão: {email}</span>
        <button class="secondary" onclick={handleLogout}>Sair</button>
      </div>
    </section>

    <section class="card">
      <h2>Cofre ({items.length} logins)</h2>
      <button class="secondary" onclick={refreshVault} disabled={busy}>Actualizar</button>
      <ul>
        {#each items as { meta, login }}
          <li>
            <strong>{login?.title ?? meta.id.slice(0, 8)}</strong>
            {#if login}
              <span class="muted"> — {login.username || "(sem user)"}</span>
            {/if}
          </li>
        {:else}
          <li class="muted">Nenhum item. Adiciona abaixo ou importa CSV.</li>
        {/each}
      </ul>
    </section>

    <section class="card">
      <h2>Novo login</h2>
      <label>Título <input bind:value={newTitle} /></label>
      <label>Utilizador <input bind:value={newUser} /></label>
      <label>
        Password
        <div class="row">
          <input type="text" bind:value={newPass} />
          <button type="button" class="secondary" onclick={fillGeneratedPassword}>Gerar</button>
        </div>
      </label>
      <button onclick={handleAddLogin} disabled={busy || !newTitle || !newPass}>Guardar cifrado</button>
    </section>

    <section class="card">
      <h2>Importar CSV</h2>
      <p class="muted">Bitwarden ou CSV genérico (title, url, username, password). Cifragem local.</p>
      <input type="file" accept=".csv,text/csv" onchange={handleImport} disabled={busy} />
      {#if importPreview}<p>{importPreview}</p>{/if}
    </section>

    <section class="card">
      <h2>Chave de recuperação</h2>
      {#if recoveryConfigured && !newRecoveryCode}
        <p class="muted">Backup configurado. Podes criar um novo (substitui o anterior).</p>
      {/if}
      {#if newRecoveryCode}
        <p class="warn">Guarda este código offline — não voltará a ser mostrado:</p>
        <code class="recovery-code">{newRecoveryCode}</code>
      {/if}
      <label>
        Confirma Master Password
        <input type="password" bind:value={recoveryConfirmPw} autocomplete="current-password" />
      </label>
      <button onclick={handleCreateRecovery} disabled={busy || !recoveryConfirmPw}>
        {recoveryConfigured ? "Regenerar chave" : "Criar chave de recuperação"}
      </button>
    </section>

    {#if webAuthnOk}
    <section class="card">
      <h2>Passkeys</h2>
      <p class="muted">Autenticação phishing-resistant ao servidor. Regista após login com password.</p>
      <ul>
        {#each passkeys as pk}
          <li>{pk.name} <span class="muted">({pk.created_at.slice(0, 10)})</span></li>
        {:else}
          <li class="muted">Nenhuma passkey registada.</li>
        {/each}
      </ul>
      <label>Nome deste dispositivo <input bind:value={passkeyName} placeholder="ex: MacBook" /></label>
      <button onclick={handleRegisterPasskey} disabled={busy || !passkeyName}>Registar passkey</button>
    </section>
    {/if}
  {/if}

  {#if status}
    <p class="status" role="status">{status}</p>
  {/if}

  <footer>
    <p class="muted">
      Backend: <code>docker compose up</code> · Frontend: <code>npm run dev</code> →
      <a href="http://localhost:5173">localhost:5173</a>
    </p>
  </footer>
</main>

<style>
  :global(body) {
    margin: 0;
    background: #0f1419;
    color: #e7ecf1;
    font-family: system-ui, sans-serif;
  }
  main {
    max-width: 42rem;
    margin: 0 auto;
    padding: 1.5rem 1rem 3rem;
  }
  header h1 {
    margin: 0;
    font-size: 1.75rem;
  }
  .tag {
    color: #8b9aab;
    margin: 0.25rem 0 1.5rem;
    font-size: 0.9rem;
  }
  .card {
    background: #1a2332;
    border: 1px solid #2a3544;
    border-radius: 8px;
    padding: 1rem 1.25rem;
    margin-bottom: 1rem;
  }
  h2 {
    margin: 0 0 0.75rem;
    font-size: 1.1rem;
  }
  label {
    display: block;
    margin-bottom: 0.75rem;
    font-size: 0.85rem;
    color: #a8b5c4;
  }
  input[type="text"],
  input[type="email"],
  input[type="password"] {
    display: block;
    width: 100%;
    margin-top: 0.25rem;
    padding: 0.5rem;
    border: 1px solid #3d4f63;
    border-radius: 4px;
    background: #0f1419;
    color: inherit;
    box-sizing: border-box;
  }
  button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    background: #3d8bfd;
    color: #fff;
    cursor: pointer;
    font-weight: 500;
  }
  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  button.secondary {
    background: #2a3544;
    color: #e7ecf1;
  }
  .row {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    flex-wrap: wrap;
  }
  .row.spread {
    justify-content: space-between;
  }
  .muted {
    color: #8b9aab;
    font-size: 0.9rem;
  }
  ul {
    padding-left: 1.25rem;
    margin: 0.75rem 0 0;
  }
  .status {
    padding: 0.75rem;
    background: #1e3a2f;
    border-radius: 4px;
    font-size: 0.9rem;
  }
  footer {
    margin-top: 2rem;
    font-size: 0.85rem;
  }
  code {
    background: #1a2332;
    padding: 0.1rem 0.35rem;
    border-radius: 3px;
  }
  a {
    color: #6cb6ff;
  }
  .warn {
    color: #f0ad4e;
    font-size: 0.9rem;
  }
  .recovery-code {
    display: block;
    font-size: 1.1rem;
    letter-spacing: 0.05em;
    padding: 0.75rem;
    margin: 0.5rem 0 1rem;
    background: #0f1419;
    border-radius: 4px;
  }
</style>
