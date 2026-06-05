<script lang="ts">
  import { HERO, PROOF_POINTS } from "$lib/content";
</script>

<svelte:head>
  <title>AegisPass — identidade e segredos da tua empresa</title>
</svelte:head>

<main>
  <section class="hero" aria-labelledby="hero-title">
    <div class="hero-inner">
      <p class="eyebrow">{HERO.eyebrow}</p>
      <h1 id="hero-title">
        {HERO.title}
        <span class="highlight">{HERO.highlight}</span>
      </h1>
      <p class="subtitle">{HERO.subtitle}</p>

      <div class="cta-row">
        <a class="btn btn-primary" href={HERO.primaryCta.href}>
          {HERO.primaryCta.label}
        </a>
        <a class="btn btn-ghost" href={HERO.secondaryCta.href}>
          {HERO.secondaryCta.label}
        </a>
      </div>

      <dl class="proof">
        {#each PROOF_POINTS as p (p.label)}
          <div class="proof-item">
            <dt class="proof-value">{p.value}</dt>
            <dd class="proof-label">{p.label}</dd>
          </div>
        {/each}
      </dl>
    </div>

    <div class="hero-glow" aria-hidden="true"></div>
  </section>

  <!-- Destino do CTA secundario. Conteudo (Zero-Knowledge + ecossistema) chega no Par B (SITE-002). -->
  <section id="como-funciona" class="placeholder" aria-label="Como funciona">
    <p>Como funciona — em breve.</p>
  </section>
</main>

<style>
  .hero {
    position: relative;
    min-height: 88vh;
    display: grid;
    place-items: center;
    padding: var(--space-16) var(--space-6);
    overflow: hidden;
  }

  .hero-inner {
    position: relative;
    z-index: 1;
    max-width: 56rem;
    text-align: center;
  }

  .eyebrow {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    font-weight: 500;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    color: var(--color-accent);
  }

  h1 {
    margin: 0;
    font-size: clamp(2.25rem, 6vw, 4.25rem);
    font-weight: 700;
    letter-spacing: -0.02em;
  }

  .highlight {
    display: block;
    background: linear-gradient(
      120deg,
      var(--color-accent),
      var(--color-link)
    );
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }

  .subtitle {
    margin: var(--space-6) auto 0;
    max-width: var(--prose-max);
    font-size: clamp(1rem, 2.2vw, 1.25rem);
    color: var(--color-text-muted);
  }

  .cta-row {
    margin-top: var(--space-8);
    display: flex;
    gap: var(--space-3);
    justify-content: center;
    flex-wrap: wrap;
  }

  .btn {
    display: inline-flex;
    align-items: center;
    padding: var(--space-3) var(--space-6);
    border-radius: var(--radius-md);
    font-size: var(--text-base);
    font-weight: 500;
    text-decoration: none;
    transition:
      transform var(--duration-fast) var(--ease-out),
      background-color var(--duration-fast) var(--ease-out);
  }

  .btn:hover {
    transform: translateY(-1px);
  }

  .btn-primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
  }

  .btn-ghost {
    background: var(--color-bg-surface);
    color: var(--color-text);
    border: 1px solid var(--color-border);
  }

  .btn-ghost:hover {
    border-color: var(--color-border-strong);
  }

  .proof {
    margin: var(--space-12) 0 0;
    display: flex;
    gap: var(--space-8);
    justify-content: center;
    flex-wrap: wrap;
  }

  .proof-item {
    text-align: center;
  }

  .proof-value {
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    font-weight: 700;
    color: var(--color-text);
  }

  .proof-label {
    margin: var(--space-1) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  /* Brilho de fundo subtil (a la Linear) — decorativo, respeita reduced-motion. */
  .hero-glow {
    position: absolute;
    top: -20%;
    left: 50%;
    width: min(48rem, 90vw);
    height: 48rem;
    transform: translateX(-50%);
    background: radial-gradient(
      circle,
      color-mix(in srgb, var(--color-accent) 18%, transparent),
      transparent 60%
    );
    filter: blur(40px);
    pointer-events: none;
  }

  .placeholder {
    min-height: 40vh;
    display: grid;
    place-items: center;
    padding: var(--space-16) var(--space-6);
    border-top: 1px solid var(--color-border);
    color: var(--color-text-muted);
  }

  @media (prefers-reduced-motion: reduce) {
    .btn:hover {
      transform: none;
    }
  }
</style>
