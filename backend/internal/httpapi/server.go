// Package httpapi monta o router HTTP da API do AegisPass.
//
// Didático: separamos a construção do router (aqui) do arranque do processo
// (cmd/server). Assim o router pode ser testado isoladamente, sem abrir portas.
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/burnnotes"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/clidevices"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/crm"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/guardian"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/emergency"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/orchestrator"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/fin"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/openbanking"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/ops"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/hr"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/passkeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/recovery"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/secretlinks"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/security"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sentinel"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharedvaults"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharekeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/shifts"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
)

// Deps são as dependências injetadas no router.
type Deps struct {
	Auth         *auth.Service
	Vault        *vault.Repo
	Hub          *realtime.Hub // nil desactiva WebSocket e notificações push
	Wipe         *security.WipeService
	Users        *users.Repo
	Shifts       *shifts.Repo
	Geofence     *geofence.Repo
	Recovery     *recovery.Repo
	Devices      *clidevices.Repo
	CLIca        *clidevices.CA
	CLICertTTL   time.Duration
	Passkeys     *passkeys.Service
	Emergency    *emergency.Service
	Sentinel     *sentinel.Service
	ShareKeys    *sharekeys.Repo    // nil desactiva endpoints /api/share/*
	SharedVaults *sharedvaults.Repo // nil desactiva endpoints /api/share/vaults*
	SecretLinks  *secretlinks.Store // nil desactiva endpoints /api/share/links*
	BurnNotes    *burnnotes.Store   // nil desactiva endpoints /api/share/notes*
	HR           *hr.Repo           // nil desactiva endpoints /api/hr/*
	Mail         *mail.Repo         // nil desactiva endpoints /api/mail/aliases*
	MailInbox    *mail.InboxRepo       // nil desactiva endpoints /api/mail/inbox*
	MailIngest        *mail.IngestService // nil desactiva webhook Mailpit
	MailRelay         *mail.RelayService  // nil desactiva compose/relay SMTP
	MailRateLimiter   *mail.RateLimiter   // nil desactiva quotas MAIL-005
	MailWebhookSecret string              // vazio desactiva webhook
	Fin          *fin.Repo              // nil desactiva endpoints /api/fin/*
	OpenBanking  *openbanking.Service   // nil desactiva endpoints /api/fin/banking/*
	Ops          *ops.Repo          // nil desactiva endpoints /api/ops/*
	Agent        *agent.Registry    // nil desactiva endpoints /api/agent/*
	AgentRunner  *agent.Runner       // executor de tools (AGENT-001)
	AgentAudit   *guardian.AuditRepo // nil desactiva GET /api/agent/audit
	Prospection  *agent.Prospection  // nil desactiva POST /api/agent/prospection/run
	Recruitment  *agent.Recruitment  // nil desactiva POST /api/agent/recruitment/run
	AgentBus     *eventbus.Bus       // nil desactiva publicação de eventos
	AgentEvents  *eventbus.PGStore      // nil desactiva GET /api/agent/events
	Orchestrator *orchestrator.Orchestrator // nil desactiva GET /api/agent/orchestrator/status
	CRM          *crm.Repo          // nil desactiva endpoints /api/crm/*
	AdminKey     string             // vazio desactiva endpoints /api/admin/*
	Pool         *pgxpool.Pool
}

// ctxKey é um tipo privado para chaves de context (evita colisões entre pacotes).
type ctxKey string

const userIDKey ctxKey = "userID"

// NewRouter devolve o http.Handler com todas as rotas registadas.
func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()
	ap := accessPolicyDeps{Shifts: deps.Shifts, Geofence: deps.Geofence}

	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /api/time", handleServerTime)
	if deps.MailIngest != nil && deps.MailWebhookSecret != "" {
		// Webhook Mailpit (MAIL-002) — sem auth; protegido por segredo na query string.
		mux.HandleFunc("POST /api/mail/webhook/mailpit", handleMailpitWebhook(deps.MailIngest, deps.MailWebhookSecret, deps.AgentBus))
	}

	if deps.Auth != nil {
		mux.HandleFunc("POST /api/auth/register", handleRegister(deps.Auth))
		mux.HandleFunc("POST /api/auth/login", handleLoginWithAccessPolicy(deps.Auth, ap, deps.Sentinel, deps.Users))
		mux.HandleFunc("GET /api/auth/kdf", handleKDFParams(deps.Auth))
	}
	if deps.Auth != nil && deps.Users != nil {
		mux.Handle("GET /api/auth/session", requireAuth(deps.Auth, handleAuthSession(deps.Users)))
		mux.Handle("GET /api/auth/sessions", requireAuth(deps.Auth, handleListSessions(deps.Auth)))
		mux.Handle("DELETE /api/auth/sessions/{id}", requireAuth(deps.Auth, handleRevokeSession(deps.Auth)))
		mux.Handle("POST /api/auth/sessions/revoke-others", requireAuth(deps.Auth, handleRevokeOtherSessions(deps.Auth)))
		mux.Handle("POST /api/auth/logout", requireAuth(deps.Auth, handleLogout(deps.Auth)))
	}
	if deps.Auth != nil && deps.Vault != nil {
		vd := vaultDeps{repo: deps.Vault, hub: deps.Hub}
		mux.Handle("GET /api/vault", requireAuthWithAccessPolicy(deps.Auth, ap, handleListItems(deps.Vault)))
		mux.Handle("POST /api/vault", requireAuthWithAccessPolicy(deps.Auth, ap, handleCreateItem(vd)))
		mux.Handle("GET /api/vault/{id}", requireAuthWithAccessPolicy(deps.Auth, ap, handleGetItem(deps.Vault)))
		mux.Handle("PUT /api/vault/{id}", requireAuthWithAccessPolicy(deps.Auth, ap, handleUpdateItem(vd)))
		mux.Handle("DELETE /api/vault/{id}", requireAuthWithAccessPolicy(deps.Auth, ap, handleDeleteItem(vd)))
	}
	if deps.Auth != nil && deps.Hub != nil {
		mux.HandleFunc("GET /api/ws/vault", handleVaultWS(deps.Auth, deps.Hub, ap))
	}
	if deps.Shifts != nil && deps.Auth != nil {
		mux.Handle("GET /api/access/shift", requireAuth(deps.Auth, handleGetAccessShift(deps.Shifts)))
	}
	if deps.Geofence != nil && deps.Auth != nil {
		mux.Handle("GET /api/access/geofence", requireAuth(deps.Auth, handleGetAccessGeofence(deps.Geofence)))
	}
	registerAdminRoutes(mux, deps)
	if deps.Wipe != nil && deps.Auth != nil {
		mux.Handle(
			"POST /api/security/remote-wipe/self",
			requireAuth(deps.Auth, handleSelfRemoteWipe(deps.Wipe)),
		)
	}
	if deps.Recovery != nil && deps.Auth != nil {
		mux.Handle("PUT /api/vault/recovery-backup", requireAuth(deps.Auth, handlePutRecoveryBackup(deps.Recovery)))
		mux.Handle("GET /api/vault/recovery-backup", requireAuth(deps.Auth, handleGetRecoveryBackupSelf(deps.Recovery)))
		mux.Handle("GET /api/vault/recovery-backup/status", requireAuth(deps.Auth, handleRecoveryBackupStatus(deps.Recovery)))
	}
	if deps.Recovery != nil {
		mux.HandleFunc("GET /api/vault/recovery-backup/lookup", handleGetRecoveryBackupByEmail(deps.Recovery))
	}
	if deps.Emergency != nil && deps.Auth != nil {
		mux.Handle("PUT /api/emergency/config", requireAuth(deps.Auth, handlePutEmergencyConfig(deps.Emergency)))
		mux.Handle("GET /api/emergency/config", requireAuth(deps.Auth, handleGetEmergencyConfig(deps.Emergency)))
		mux.Handle("DELETE /api/emergency/config", requireAuth(deps.Auth, handleDeleteEmergencyConfig(deps.Emergency)))
		mux.Handle("GET /api/emergency/requests", requireAuth(deps.Auth, handleListEmergencyRequests(deps.Emergency)))
		mux.Handle("POST /api/emergency/requests/{id}/reject", requireAuth(deps.Auth, handleRejectEmergencyRequest(deps.Emergency)))
		mux.Handle("POST /api/emergency/requests/{id}/approve", requireAuth(deps.Auth, handleApproveEmergencyRequest(deps.Emergency)))
		mux.Handle("POST /api/emergency/request", requireAuth(deps.Auth, handleCreateEmergencyRequest(deps.Emergency)))
		mux.Handle("GET /api/emergency/request/status", requireAuth(deps.Auth, handleGetEmergencyRequestStatus(deps.Emergency)))
		mux.Handle("GET /api/emergency/access", requireAuth(deps.Auth, handleFetchEmergencyAccess(deps.Emergency)))
	}
	registerCLIRoutes(mux, deps)
	registerPasskeyRoutes(mux, deps, ap)
	if deps.Sentinel != nil && deps.Auth != nil {
		mux.Handle("GET /api/security/sentinel/events", requireAuth(deps.Auth, handleListSentinelEvents(deps.Sentinel)))
	}
	if deps.Sentinel != nil && deps.Passkeys != nil && deps.Users != nil && deps.Auth != nil {
		mux.HandleFunc("POST /api/auth/sentinel/step-up/begin", handleSentinelStepUpBegin(deps.Passkeys, deps.Sentinel, deps.Users))
		mux.HandleFunc("POST /api/auth/sentinel/step-up/finish", handleSentinelStepUpFinish(deps.Auth, deps.Passkeys, deps.Sentinel, deps.Users, ap))
	}
	if deps.ShareKeys != nil && deps.Auth != nil {
		mux.Handle("PUT /api/share/keypair", requireAuth(deps.Auth, handlePutShareKeypair(deps.ShareKeys)))
		mux.Handle("GET /api/share/keypair", requireAuth(deps.Auth, handleGetShareKeypair(deps.ShareKeys)))
		mux.Handle("GET /api/share/keypair/status", requireAuth(deps.Auth, handleShareKeypairStatus(deps.ShareKeys)))
		mux.Handle("GET /api/share/public-key", requireAuth(deps.Auth, handleGetSharePublicKey(deps.ShareKeys)))
	}
	if deps.SharedVaults != nil && deps.Auth != nil {
		// Cofres
		mux.Handle("POST /api/share/vaults", requireAuth(deps.Auth, handleCreateSharedVault(deps.SharedVaults)))
		mux.Handle("GET /api/share/vaults", requireAuth(deps.Auth, handleListSharedVaults(deps.SharedVaults)))
		mux.Handle("GET /api/share/vaults/{id}", requireAuth(deps.Auth, handleGetSharedVault(deps.SharedVaults)))
		mux.Handle("DELETE /api/share/vaults/{id}", requireAuth(deps.Auth, handleDeleteSharedVault(deps.SharedVaults)))
		// Membros (permissões + revogação)
		mux.Handle("GET /api/share/vaults/{id}/members", requireAuth(deps.Auth, handleListSharedVaultMembers(deps.SharedVaults)))
		mux.Handle("POST /api/share/vaults/{id}/members", requireAuth(deps.Auth, handleAddSharedVaultMember(deps.SharedVaults)))
		mux.Handle("DELETE /api/share/vaults/{id}/members/{userId}", requireAuth(deps.Auth, handleRemoveSharedVaultMember(deps.SharedVaults)))
		// Itens do cofre
		mux.Handle("GET /api/share/vaults/{id}/items", requireAuth(deps.Auth, handleListSharedVaultItems(deps.SharedVaults)))
		mux.Handle("POST /api/share/vaults/{id}/items", requireAuth(deps.Auth, handleCreateSharedVaultItem(deps.SharedVaults)))
		mux.Handle("DELETE /api/share/vaults/{id}/items/{itemId}", requireAuth(deps.Auth, handleDeleteSharedVaultItem(deps.SharedVaults)))
		// Anexos cifrados por ficheiro (SHARE-004)
		mux.Handle("GET /api/share/vaults/{id}/attachments", requireAuth(deps.Auth, handleListVaultAttachments(deps.SharedVaults)))
		mux.Handle("POST /api/share/vaults/{id}/attachments", requireAuth(deps.Auth, handleAddVaultAttachment(deps.SharedVaults)))
		mux.Handle("GET /api/share/vaults/{id}/attachments/{attId}", requireAuth(deps.Auth, handleGetVaultAttachment(deps.SharedVaults)))
		mux.Handle("DELETE /api/share/vaults/{id}/attachments/{attId}", requireAuth(deps.Auth, handleDeleteVaultAttachment(deps.SharedVaults)))
	}
	if deps.SecretLinks != nil {
		// Criar exige sessão (só utilizadores criam links).
		if deps.Auth != nil {
			mux.Handle("POST /api/share/links", requireAuth(deps.Auth, handleCreateSecretLink(deps.SecretLinks)))
		}
		// Consumir é PÚBLICO: a chave de cifra vive no fragmento do URL, que
		// nunca chega ao servidor. Qualquer pessoa com o link o pode abrir 1x.
		mux.HandleFunc("POST /api/share/links/{id}", handleConsumeSecretLink(deps.SecretLinks))
	}
	if deps.BurnNotes != nil {
		// Criar exige sessão (só utilizadores criam notas auto-destrutivas).
		if deps.Auth != nil {
			mux.Handle("POST /api/share/notes", requireAuth(deps.Auth, handleCreateBurnNote(deps.BurnNotes)))
		}
		// Ler (queima após leitura) e queimar manualmente (capacidade via token)
		// são PÚBLICOS: a chave de cifra vive no fragmento do URL.
		mux.HandleFunc("POST /api/share/notes/{id}", handleConsumeBurnNote(deps.BurnNotes))
		mux.HandleFunc("POST /api/share/notes/{id}/burn", handleBurnNote(deps.BurnNotes))
	}
	if deps.HR != nil && deps.Auth != nil {
		// Fichas de empregado com cifragem campo-a-campo (HR-001). Tudo exige
		// sessão; cada utilizador só vê e gere as suas próprias fichas.
		mux.Handle("POST /api/hr/employees", requireAuth(deps.Auth, handleCreateEmployeeRecord(deps.HR, deps.AgentBus)))
		mux.Handle("GET /api/hr/employees", requireAuth(deps.Auth, handleListEmployeeRecords(deps.HR)))
		mux.Handle("GET /api/hr/employees/{id}", requireAuth(deps.Auth, handleGetEmployeeRecord(deps.HR)))
		mux.Handle("DELETE /api/hr/employees/{id}", requireAuth(deps.Auth, handleDeleteEmployeeRecord(deps.HR)))
		// Campos individuais (upsert/remover por nome de campo).
		mux.Handle("PUT /api/hr/employees/{id}/fields/{field}", requireAuth(deps.Auth, handlePutEmployeeField(deps.HR)))
		mux.Handle("DELETE /api/hr/employees/{id}/fields/{field}", requireAuth(deps.Auth, handleDeleteEmployeeField(deps.HR)))
		// Crypto-shredding RGPD (HR-003) + certificados de eliminação (HR-004).
		mux.Handle("POST /api/hr/employees/{id}/fields/{field}/shred", requireAuth(deps.Auth, handleShredEmployeeField(deps.HR)))
		mux.Handle("POST /api/hr/employees/{id}/shred", requireAuth(deps.Auth, handleShredEmployeeRecord(deps.HR)))
		mux.Handle("GET /api/hr/certificates", requireAuth(deps.Auth, handleListErasureCertificates(deps.HR)))
		// Logs imutáveis (HR-002) + relatório de conformidade RGPD (HR-008).
		mux.Handle("GET /api/hr/audit", requireAuth(deps.Auth, handleListAuditLog(deps.HR)))
		mux.Handle("GET /api/hr/compliance-report", requireAuth(deps.Auth, handleComplianceReport(deps.HR)))
		// Identidade de assinatura (HR-006).
		mux.Handle("PUT /api/hr/signing-identity", requireAuth(deps.Auth, handlePutSigningIdentity(deps.HR)))
		mux.Handle("GET /api/hr/signing-identity", requireAuth(deps.Auth, handleGetSigningIdentity(deps.HR)))
		mux.Handle("GET /api/hr/signers/{userId}/public-key", requireAuth(deps.Auth, handleGetSignerPublicKey(deps.HR)))
		// Contratos cifrados (HR-005) + assinatura (HR-006).
		mux.Handle("GET /api/hr/employees/{id}/contracts", requireAuth(deps.Auth, handleListContracts(deps.HR)))
		mux.Handle("POST /api/hr/employees/{id}/contracts", requireAuth(deps.Auth, handleAddContract(deps.HR)))
		mux.Handle("GET /api/hr/employees/{id}/contracts/{cid}", requireAuth(deps.Auth, handleGetContract(deps.HR)))
		mux.Handle("POST /api/hr/employees/{id}/contracts/{cid}/sign", requireAuth(deps.Auth, handleSignContract(deps.HR)))
		mux.Handle("DELETE /api/hr/employees/{id}/contracts/{cid}", requireAuth(deps.Auth, handleDeleteContract(deps.HR)))
	}
	if deps.Mail != nil && deps.Auth != nil {
		// Aliases de e-mail (MAIL-001).
		mux.Handle("GET /api/mail/aliases", requireAuth(deps.Auth, handleListAliases(deps.Mail)))
		mux.Handle("POST /api/mail/aliases", requireAuth(deps.Auth, handleCreateAlias(deps.Mail)))
		mux.Handle("PATCH /api/mail/aliases/{id}", requireAuth(deps.Auth, handleSetAliasActive(deps.Mail)))
		mux.Handle("DELETE /api/mail/aliases/{id}", requireAuth(deps.Auth, handleDeleteAlias(deps.Mail)))
		if deps.MailRelay != nil {
			mux.Handle("POST /api/mail/compose", requireAuth(deps.Auth, handleComposeMail(deps.Mail, deps.MailRelay, deps.MailRateLimiter)))
		}
	}
	if deps.MailInbox != nil && deps.Auth != nil {
		// Caixa de entrada simulada (AGENT-003 stub até MAIL-002).
		mux.Handle("GET /api/mail/inbox", requireAuth(deps.Auth, handleListInbox(deps.MailInbox)))
		mux.Handle("POST /api/mail/inbox", requireAuth(deps.Auth, handleCreateInboxMessage(deps.MailInbox, deps.AgentBus)))
		mux.Handle("POST /api/mail/inbox/{id}/processed", requireAuth(deps.Auth, handleMarkInboxProcessed(deps.MailInbox)))
	}
	if deps.Fin != nil && deps.Auth != nil {
		// Monitorizacao de custos SaaS (FIN-001).
		mux.Handle("GET /api/fin/subscriptions", requireAuth(deps.Auth, handleListSubscriptions(deps.Fin)))
		mux.Handle("POST /api/fin/subscriptions", requireAuth(deps.Auth, handleCreateSubscription(deps.Fin)))
		mux.Handle("PUT /api/fin/subscriptions/{id}", requireAuth(deps.Auth, handleUpdateSubscription(deps.Fin)))
		mux.Handle("DELETE /api/fin/subscriptions/{id}", requireAuth(deps.Auth, handleDeleteSubscription(deps.Fin)))
	}
	if deps.OpenBanking != nil && deps.Auth != nil {
		// Open Banking scaffold (FIN-003) — mock provider em dev.
		mux.Handle("GET /api/fin/banking/status", requireAuth(deps.Auth, handleBankingStatus(deps.OpenBanking)))
		mux.Handle("POST /api/fin/banking/connect", requireAuth(deps.Auth, handleBankingConnect(deps.OpenBanking)))
		mux.Handle("POST /api/fin/banking/sync", requireAuth(deps.Auth, handleBankingSync(deps.OpenBanking)))
	}
	if deps.Ops != nil && deps.Auth != nil {
		// Inventário operacional (AGENT-008).
		mux.Handle("GET /api/ops/inventory", requireAuth(deps.Auth, handleListInventory(deps.Ops)))
		mux.Handle("POST /api/ops/inventory", requireAuth(deps.Auth, handleCreateInventory(deps.Ops, deps.AgentBus)))
		mux.Handle("PUT /api/ops/inventory/{id}", requireAuth(deps.Auth, handleUpdateInventory(deps.Ops, deps.AgentBus)))
		mux.Handle("POST /api/ops/inventory/{id}/adjust", requireAuth(deps.Auth, handleAdjustInventory(deps.Ops, deps.AgentBus)))
		mux.Handle("DELETE /api/ops/inventory/{id}", requireAuth(deps.Auth, handleDeleteInventory(deps.Ops)))
	}
	if deps.Agent != nil && deps.AgentRunner != nil && deps.Auth != nil {
		// Sistema de tools para agentes (AGENT-001) + auditoria Guardião (AGENT-002).
		mux.Handle("GET /api/agent/tools", requireAuth(deps.Auth, handleListAgentTools(deps.Agent)))
		mux.Handle("POST /api/agent/tools/{name}/run", requireAuth(deps.Auth, handleRunAgentTool(deps.AgentRunner, deps.AgentAudit, deps.AgentBus)))
		if deps.AgentAudit != nil {
			mux.Handle("GET /api/agent/audit", requireAuth(deps.Auth, handleListAgentAudit(deps.AgentAudit)))
		}
		if deps.Prospection != nil {
			mux.Handle("POST /api/agent/prospection/run", requireAuth(deps.Auth, handleRunProspection(deps.Prospection, deps.AgentAudit, deps.AgentBus)))
		}
		if deps.Recruitment != nil {
			mux.Handle("POST /api/agent/recruitment/run", requireAuth(deps.Auth, handleRunRecruitment(deps.Recruitment, deps.AgentAudit, deps.AgentBus)))
		}
		if deps.AgentBus != nil {
			mux.Handle("POST /api/agent/finance/report-stale", requireAuth(deps.Auth, handleReportStaleSubscriptions(deps.AgentBus)))
			mux.Handle("POST /api/agent/finance/report-sync", requireAuth(deps.Auth, handleReportTransactionsSynced(deps.AgentBus)))
		}
		if deps.AgentEvents != nil {
			mux.Handle("GET /api/agent/events", requireAuth(deps.Auth, handleListAgentEvents(deps.AgentEvents)))
		}
		if deps.Orchestrator != nil {
			mux.Handle("GET /api/agent/orchestrator/status", requireAuth(deps.Auth, handleOrchestratorStatus(deps.Orchestrator)))
		}
		if deps.AgentEvents != nil {
			mux.Handle("POST /api/agent/orchestrator/actions/{id}/approve", requireAuth(deps.Auth, handleDecideOrchestratorAction(deps.AgentEvents, deps.AgentBus, eventbus.DecisionApprove)))
			mux.Handle("POST /api/agent/orchestrator/actions/{id}/reject", requireAuth(deps.Auth, handleDecideOrchestratorAction(deps.AgentEvents, deps.AgentBus, eventbus.DecisionReject)))
		}
	}
	if deps.CRM != nil && deps.Auth != nil {
		mux.Handle("GET /api/crm/leads", requireAuth(deps.Auth, handleListLeads(deps.CRM)))
		mux.Handle("POST /api/crm/leads", requireAuth(deps.Auth, handleCreateLead(deps.CRM)))
		mux.Handle("PUT /api/crm/leads/{id}", requireAuth(deps.Auth, handleUpdateLead(deps.CRM)))
		mux.Handle("DELETE /api/crm/leads/{id}", requireAuth(deps.Auth, handleDeleteLead(deps.CRM)))
	}

	return mux
}

// requireAuth é um middleware: valida o token Bearer e, se válido, injeta o
// userID no context antes de chamar o handler seguinte.
func requireAuth(svc *auth.Service, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "token em falta")
			return
		}
		userID, err := svc.Authenticate(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "sessão inválida")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// bearerToken extrai o token do header "Authorization: Bearer <token>".
func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if strings.HasPrefix(h, prefix) {
		return strings.TrimPrefix(h, prefix)
	}
	return ""
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// writeJSON escreve uma resposta JSON com o status indicado.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError padroniza as respostas de erro: { "error": "mensagem" }.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
