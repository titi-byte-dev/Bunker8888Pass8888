# PRD — Site institucional AegisPass (aegispass.com)

**Versao:** v1 · **Estado:** aprovado para escopo · **Data:** 2026-06-06
**Tipo:** propriedade web publica, independente da app (`app.aegispass.com`)

> Fonte canonica da narrativa publica do AegisPass. Toda a comunicacao
> (pitch, parceiros, email, social) reutiliza a frase-ancora e o posicionamento
> definidos aqui.

---

## Decisoes ja fechadas (nao reabrir sem motivo)

| Decisao | Escolha | Porque |
|---|---|---|
| Independente vs. dentro da app | **Independente** (`aegispass.com` ⟂ `app.aegispass.com`) | Objetivos opostos: convencer estranho vs. dar poder a quem fez login. SEO/velocidade vs. densidade/estado. Habilita white-label de parceiro. |
| Objetivo da v1 | **Institucional** (apresentar/credibilidade) | Sem funil agressivo, sem trials, sem captura de leads. |
| Camada Ops (RH/Fin/CRM/Mail) | **Escondida na v1** | Narrativa de 2 camadas (Core + Workspace). Reforca o foco "somos o Cofre", evita a dispersao de "fazemos tudo". |
| Frase-ancora | **"A camada de identidade e segredos da tua empresa. Tudo assenta no Cofre."** | Foca a visao de plataforma/ecossistema. |
| Acoplamento com a app | So `tokens.css` + paletas | Codigo separado, linguagem visual partilhada. |

---

## Problem Statement

O AegisPass existe como app autenticada mas **nao tem rosto publico**. Um
estranho (cliente potencial, parceiro white-label, candidato) nao tem como
perceber em 30 segundos o que e, porque e diferente (zero-knowledge + BYOD +
white-label) e como as pecas encaixam. Sem isto, cada conversa comeca do zero e
a estrategia de "ecossistema em camadas sobre o Cofre" fica invisivel. Custo de
nao resolver: zero credibilidade inbound, dependencia total de venda 1-a-1,
nenhum canal para recrutar parceiros.

## Goals

1. Comunicar o ecossistema em camadas (Core / Workspace) de forma que o
   visitante perceba "tudo assenta no Cofre" sem ler documentacao.
2. Estabelecer credibilidade tecnica — zero-knowledge legivel e crivel, nao jargao.
3. Ser a fonte canonica da narrativa (frase-ancora + posicionamento).
4. Coerencia de marca com a app (partilha `tokens.css`/paletas).
5. Performance e SEO de topo (estatico, <1s, indexavel).

## Non-Goals

- **Nao e funil de conversao agressivo** — sem trials, formularios, A/B de CTAs.
- **Nao tem login nem logica de sessao** — "Entrar" e um link para a app.
- **Nao documenta features ao detalhe** — isso vive na doc dentro da app.
- **Nao mostra a camada Ops na v1** — escondida para reforcar o foco no Core.
- **Nao e multi-tenant/white-label na v1** — o site e a marca AegisPass; subdominios de parceiro sao v2.
- **Nao inclui blog/CMS** — conteudo editorial e v2.

## User Stories

**Visitante avaliador (decisor PME / IT)**
- Como decisor que nunca ouviu falar do produto, quero perceber numa frase o que
  o AegisPass faz, para decidir se continuo a ler.
- Como avaliador cetico, quero perceber o que "zero-knowledge" significa na
  pratica (o servidor nunca ve os meus segredos), para confiar.
- Como decisor, quero ver um ecossistema coerente (nao ferramentas soltas), para
  perceber o valor de adotar a plataforma.

**Parceiro potencial (MSP / contabilista)**
- Como MSP, quero perceber que existe uma camada de plataforma (Cofre) a qual
  outras coisas ligam, para imaginar revender isto aos meus clientes. *(semente white-label v2)*

**Candidato / curioso tecnico**
- Como visitante tecnico, quero ver a app real a mexer (nao screenshots mortos),
  para avaliar a qualidade do produto.

## Requirements

### Must-Have (P0)

| # | Requisito | Criterio de aceitacao |
|---|---|---|
| P0-1 | **Hero com frase-ancora** | Above-the-fold: headline "A camada de identidade e segredos da tua empresa. Tudo assenta no Cofre.", sub-linha de 1 frase, CTA primario "Entrar" -> `app.aegispass.com`, CTA secundario "Ver como funciona" (scroll). |
| P0-2 | **Zero-Knowledge explicado** | Diagrama: cliente cifra -> servidor guarda blob opaco -> servidor nunca decifra. Linguagem de negocio, nao criptografica. |
| P0-3 | **Ecossistema em 2 camadas** | Bloco visual: Core (Cofre/Equipa/Seguranca) como base, Workspace (Trabalho: turnos, sandbox, CLI, inventario, Google) como camada que assenta nele. Mostra hierarquia, nao menu plano. Ops omitida. |
| P0-4 | **Coerencia de marca** | Importa `tokens.css` + paleta `aegis`. Inter/Outfit, dark por omissao. |
| P0-5 | **Estatico e rapido** | Build estatico (Astro recomendado). Lighthouse Performance >=95, sem JS de app. Repo/pasta separados. |
| P0-6 | **Responsivo + acessivel** | Mobile-first, WCAG AA (contraste herdado dos tokens validados), navegavel por teclado. |
| P0-7 | **Rodape institucional** | Quem somos (1 linha), contacto, links legais (placeholder), link app. |

### Nice-to-Have (P1)

- **Produto a mexer** — clip/loop curto da app real (a la Linear) ou screenshot interativo.
- **Seccao "para parceiros"** — teaser do white-label (semente do canal).
- **Social/OG cards** — partilha bonita em LinkedIn/X.
- **Animacoes de scroll subtis** — respeitando `prefers-reduced-motion`.

### Future Considerations (P2)

- **Subdominios white-label de parceiro** (`secureco.aegispass.com`) — manter o template parametrizavel por marca desde ja.
- **Camada Ops como integracoes** (Moloni/Factorial/Pipedrive) quando houver provas.
- **Blog/CMS** para SEO de conteudo.
- **Pagina de precos** quando o modelo comercial estabilizar.
- **i18n (PT/EN)** — estruturar copy em ficheiros separados desde o inicio.

## Success Metrics

**Leading (dias->semanas)**
- Lighthouse Performance >=95, Acessibilidade >=95, SEO >=95.
- Tempo ate interativo <1s em 4G.
- Taxa de clique no CTA "Entrar" (proxy de interesse).

**Lagging (semanas->meses)**
- Nº de conversas de parceria/venda que comecam com "vi o vosso site" (qualitativo).
- Reducao do tempo de explicacao numa demo (narrativa ja absorvida).

## Open Questions

- **[Stack — eng]** Astro vs. SvelteKit `adapter-static`. Recomendacao: **Astro** (mais leve para conteudo; mantem a app SvelteKit pura). Nao-bloqueante.
- **[Legal — negocio]** Entidade/morada/politica de privacidade para o rodape, ou placeholders?
- **[Asset — negocio]** Existe clip/gravacao da app para o P1, ou screenshots na v1?

## Timeline Considerations

- **Sem deadline duro** (institucional). Candidato ao loop didatico:
  - Par A: scaffold do site (pasta separada, tokens partilhados) + hero (P0-1, P0-4, P0-5).
  - Par B: seccao ZK (P0-2) + ecossistema 2 camadas (P0-3).
  - Par C: responsivo/a11y (P0-6) + rodape (P0-7).
- **Dependencia leve:** extrair/partilhar `tokens.css` + paletas (pequeno trabalho de fronteira).
- **Faseamento:** v1 = P0 (one-pager estatico). v1.1 = P1 (clip + teaser parceiros). v2 = white-label/Ops/blog/precos.

---

## Arquitetura de fronteira (referencia)

```
aegispass.com  (Astro estatico, publico, SEO)
   |  importa tokens.css + paleta aegis
   |  CTA "Entrar" / "Experimentar"
   v
app.aegispass.com  (SvelteKit, autenticado, ZK)
   └─ resolve tenant -> paleta white-label
```

So acoplamento: o token de design. Codigo, deploy e CI separados.
