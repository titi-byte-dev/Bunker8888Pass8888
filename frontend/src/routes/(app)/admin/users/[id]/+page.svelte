<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import AdminGate from "$lib/admin/AdminGate.svelte";
  import { hasAdminKey } from "$lib/admin/adminKey";
  import {
    formatSimpleSchedule,
    getAdminGeofencePolicy,
    getAdminShiftPolicy,
    getAdminUser,
    parseSimpleSchedule,
    setAdminGeofencePolicy,
    setAdminShiftPolicy,
    triggerRemoteWipe,
  } from "$lib/admin/api";
  import type { ShiftPolicy } from "$lib/vault/shift";
  import type { GeofencePolicy } from "$lib/vault/geofence";
  import { PageShell, StatusBanner } from "$lib/ui";

  const userId = $derived(page.params.id ?? "");

  let unlocked = $state(hasAdminKey());
  let email = $state("");
  let shiftEnabled = $state(false);
  let shiftTz = $state("Europe/Lisbon");
  let shiftSkew = $state(300);
  let shiftScheduleText = $state("");
  let geoEnabled = $state(false);
  let geoCidrs = $state("");
  let geoGpsEnabled = $state(false);
  let geoLat = $state("");
  let geoLon = $state("");
  let geoRadius = $state(500);
  let wipeReason = $state("");
  let busy = $state(true);
  let saving = $state(false);
  let status = $state("");
  let error = $state("");

  async function load() {
    if (!hasAdminKey() || !userId) return;
    busy = true;
    error = "";
    status = "";
    try {
      const [user, shift, geo] = await Promise.all([
        getAdminUser(userId),
        getAdminShiftPolicy(userId),
        getAdminGeofencePolicy(userId),
      ]);
      email = user.email;
      applyShift(shift);
      applyGeo(geo);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      busy = false;
    }
  }

  function applyShift(p: ShiftPolicy) {
    shiftEnabled = p.enabled;
    shiftTz = p.timezone || "Europe/Lisbon";
    shiftSkew = p.max_clock_skew_seconds || 300;
    shiftScheduleText = formatSimpleSchedule(p.schedule);
  }

  function applyGeo(p: GeofencePolicy) {
    geoEnabled = p.enabled;
    geoCidrs = (p.allowed_cidrs ?? []).join("\n");
    geoGpsEnabled = p.gps_enabled;
    geoLat = p.gps_lat != null ? String(p.gps_lat) : "";
    geoLon = p.gps_lon != null ? String(p.gps_lon) : "";
    geoRadius = p.gps_radius_m || 500;
  }

  async function saveShift() {
    saving = true;
    error = "";
    status = "";
    try {
      await setAdminShiftPolicy(userId, {
        enabled: shiftEnabled,
        timezone: shiftTz,
        max_clock_skew_seconds: shiftSkew,
        schedule: parseSimpleSchedule(shiftScheduleText),
      });
      status = "Turno gravado.";
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao gravar turno";
    } finally {
      saving = false;
    }
  }

  async function saveGeo() {
    saving = true;
    error = "";
    status = "";
    try {
      await setAdminGeofencePolicy(userId, {
        enabled: geoEnabled,
        allowed_cidrs: geoCidrs
          .split("\n")
          .map((s) => s.trim())
          .filter(Boolean),
        gps_enabled: geoGpsEnabled,
        gps_lat: geoLat ? Number(geoLat) : null,
        gps_lon: geoLon ? Number(geoLon) : null,
        gps_radius_m: geoRadius,
      });
      status = "Geofence gravado.";
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao gravar geofence";
    } finally {
      saving = false;
    }
  }

  async function handleWipe() {
    if (!confirm(`Remote wipe para ${email}? Revoga sessões e apaga dados locais nos dispositivos online.`)) {
      return;
    }
    saving = true;
    error = "";
    status = "";
    try {
      const result = await triggerRemoteWipe(userId, wipeReason.trim());
      status = `Wipe enviado: ${result.devices_notified} dispositivo(s), ${result.sessions_revoked} sessão(ões) revogada(s).`;
      wipeReason = "";
    } catch (e) {
      error = e instanceof Error ? e.message : "Remote wipe falhou";
    } finally {
      saving = false;
    }
  }

  onMount(load);

  function onGateChange() {
    unlocked = hasAdminKey();
    if (unlocked) load();
  }
</script>

<svelte:head>
  <title>{email || "Utilizador"} — Admin</title>
</svelte:head>

<PageShell
  title={email || "Utilizador"}
  leaf={email || undefined}
  description="Políticas de turno, geofence e remote wipe para este colaborador."
>
  <p class="id mono">{userId}</p>

  <AdminGate onUnlocked={onGateChange} />

  {#if unlocked}
    {#if status}
      <StatusBanner variant="success">{status}</StatusBanner>
    {/if}
    {#if error}
      <StatusBanner variant="error">{error}</StatusBanner>
    {/if}

    {#if busy}
      <p class="muted">A carregar…</p>
    {:else}
      <section class="block">
        <h2>Turno (VAULT-010)</h2>
        <label class="check">
          <input type="checkbox" bind:checked={shiftEnabled} />
          Restrição horária activa
        </label>
        <label>
          Fuso horário
          <input type="text" bind:value={shiftTz} placeholder="Europe/Lisbon" />
        </label>
        <label>
          Desvio máximo relógio (s)
          <input type="number" bind:value={shiftSkew} min="30" max="3600" />
        </label>
        <label>
          Horário por dia (seg–dom, uma linha por dia)
          <span class="hint">Formato: 09:00-17:00 ou «-» para dia livre</span>
          <textarea rows="7" bind:value={shiftScheduleText}></textarea>
        </label>
        <button type="button" disabled={saving} onclick={saveShift}>Gravar turno</button>
      </section>

      <section class="block danger-zone">
        <h2>Geofence (VAULT-011)</h2>
        <label class="check">
          <input type="checkbox" bind:checked={geoEnabled} />
          Geofence activo
        </label>
        <label>
          CIDRs permitidos (um por linha)
          <textarea rows="3" bind:value={geoCidrs} placeholder="192.168.1.0/24"></textarea>
        </label>
        <label class="check">
          <input type="checkbox" bind:checked={geoGpsEnabled} />
          Exigir GPS (círculo)
        </label>
        <div class="row">
          <label>
            Latitude
            <input type="text" bind:value={geoLat} placeholder="38.7223" />
          </label>
          <label>
            Longitude
            <input type="text" bind:value={geoLon} placeholder="-9.1393" />
          </label>
          <label>
            Raio (m)
            <input type="number" bind:value={geoRadius} min="50" />
          </label>
        </div>
        <button type="button" disabled={saving} onclick={saveGeo}>Gravar geofence</button>
      </section>

      <section class="block danger-zone">
        <h2>Remote wipe (VAULT-012)</h2>
        <p class="hint">
          Envia push WebSocket aos dispositivos online e revoga todas as sessões HTTP.
          Acção irreversível nos dados locais do cofre.
        </p>
        <label>
          Motivo (auditoria)
          <input type="text" bind:value={wipeReason} placeholder="Dispositivo perdido" />
        </label>
        <button type="button" class="danger" disabled={saving} onclick={handleWipe}>
          Executar remote wipe
        </button>
      </section>
    {/if}
  {/if}
</PageShell>

<style>
  .id {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    margin: 0;
    word-break: break-all;
  }

  .mono {
    font-family: var(--font-mono);
  }

  .block {
    margin-bottom: var(--space-8);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
  }

  .danger-zone {
    border-color: color-mix(in srgb, var(--color-danger) 25%, var(--color-border));
  }

  h2 {
    margin: 0 0 var(--space-4);
    font-size: var(--text-lg);
  }

  label {
    display: block;
    font-size: var(--text-sm);
    margin-bottom: var(--space-3);
  }

  .check {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  input[type="text"],
  input[type="number"],
  textarea {
    display: block;
    width: 100%;
    margin-top: var(--space-1);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    font-size: var(--text-sm);
    box-sizing: border-box;
    font-family: inherit;
  }

  textarea {
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    resize: vertical;
  }

  .row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(8rem, 1fr));
    gap: var(--space-3);
  }

  .hint {
    display: block;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    margin-top: var(--space-1);
  }

  button {
    margin-top: var(--space-2);
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  button.danger {
    border-color: var(--color-danger);
    color: var(--color-danger);
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
</style>
