/**
 * Cliente WebSocket para sync do cofre (VAULT-006).
 *
 * Didático: quando outro dispositivo altera o cofre, o servidor envia um evento
 * push. Este cliente reconecta automaticamente se a ligação cair (backoff
 * exponencial até 30s). O conteúdo sensível NÃO vem pelo WS — só metadados;
 * usa VaultAPI.get() para obter o blob cifrado actualizado.
 */
import type { RealtimeMessage, RemoteWipeEvent, VaultSyncEvent } from "./types";
import { isRemoteWipeEvent } from "./types";
import type { GeoPosition } from "./geofence";

const MAX_BACKOFF_MS = 30_000;

export class VaultSync {
  private ws: WebSocket | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private backoffMs = 1_000;
  private closed = false;

  constructor(
    private baseURL: string,
    private token: string,
    private onVaultEvent: (ev: VaultSyncEvent) => void,
    private onStatus?: (connected: boolean) => void,
    private onRemoteWipe?: (ev: RemoteWipeEvent) => void,
    private geo: GeoPosition | null = null,
  ) {}

  /** Abre (ou reabre) a ligação WebSocket. */
  connect(): void {
    this.closed = false;
    this.openSocket();
  }

  /** Fecha a ligação e cancela reconexões. */
  disconnect(): void {
    this.closed = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.ws?.close();
    this.ws = null;
    this.onStatus?.(false);
  }

  /** Actualiza GPS usado na ligação WebSocket (query lat/lon). */
  setGeo(pos: GeoPosition | null): void {
    this.geo = pos;
  }

  private openSocket(): void {
    const url = buildWsURL(this.baseURL, this.token, this.geo);
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this.backoffMs = 1_000;
      this.onStatus?.(true);
    };

    this.ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(String(ev.data)) as RealtimeMessage;
        if (isRemoteWipeEvent(msg)) {
          this.onRemoteWipe?.(msg);
          return;
        }
        this.onVaultEvent(msg);
      } catch {
        // Mensagem inválida — ignoramos (dados ≠ instruções).
      }
    };

    this.ws.onclose = () => {
      this.onStatus?.(false);
      if (!this.closed) this.scheduleReconnect();
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.reconnectTimer = setTimeout(() => {
      this.backoffMs = Math.min(this.backoffMs * 2, MAX_BACKOFF_MS);
      this.openSocket();
    }, this.backoffMs);
  }
}

/** Converte base URL HTTP(S) em URL WebSocket com token (e GPS opcional na query). */
export function buildWsURL(baseURL: string, token: string, geo: GeoPosition | null = null): string {
  const url = new URL(baseURL);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  url.pathname = "/api/ws/vault";
  const params = new URLSearchParams({ token });
  if (geo) {
    params.set("lat", String(geo.lat));
    params.set("lon", String(geo.lon));
  }
  url.search = params.toString();
  return url.toString();
}
