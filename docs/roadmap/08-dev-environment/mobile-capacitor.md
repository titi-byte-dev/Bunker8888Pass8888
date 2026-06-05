# UI-009 — App mobile com Capacitor

> **Objetivo:** empacotar o SvelteKit como app iOS/Android com desbloqueio biométrico
> e suporte a remote wipe em dispositivo real (DoD Fase 1).

> 💡 **Conceito:** Capacitor envolve a WebView existente — reutilizamos o mesmo código
> Svelte; plugins nativos expõem biometria, push e secure storage.

## Pré-requisitos

- Node 22+ (igual ao CI)
- Android Studio (APK) e/ou Xcode (iOS) no portátil de build
- `ssr = false` já activo no projeto (obrigatório para WebView)

## 1. Build estático

```bash
cd frontend
npm install @capacitor/core @capacitor/cli @capacitor/android @capacitor/ios
npm install -D @sveltejs/adapter-static
```

Em `svelte.config.js`, usar `adapter-static` com `fallback: 'index.html'` para SPA.

```bash
npm run build
npx cap init AegisPass com.aegispass.app --web-dir build
npx cap add android
npx cap add ios
npx cap sync
```

## 2. Biometria (fase seguinte)

```bash
npm install @capacitor-community/biometric-auth
```

Ligar em `frontend/src/lib/platform/native.ts` → `unlockWithBiometric()`.

> ⚠️ **Segurança:** biometria **não** guarda a Master Key no Keychain por omissão —
> só confirma identidade antes de carregar a chave já derivada em memória.

## 3. Remote wipe mobile

O wipe existente (`VAULT-012`) usa WebSocket push — na app nativa:

1. Registar device token (FCM/APNs) — task futura `UI-009b`
2. Ao receber evento `security.remote_wipe`, chamar `purgeMasterKey()` + limpar storage

## 4. Desenvolvimento

```bash
npm run dev          # browser
npm run build && npx cap run android   # emulador/dispositivo
```

## Estado actual

| Item | Estado |
|---|---|
| `platform/native.ts` (stubs) | 🟢 |
| Capacitor init + adapters | 🟡 documentado |
| Biometria real | ⚪ |
| Push wipe mobile | ⚪ |
