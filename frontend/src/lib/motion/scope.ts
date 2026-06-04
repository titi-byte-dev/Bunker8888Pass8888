import { gsap } from "gsap";

/**
 * Executa animações GSAP num scope DOM e devolve cleanup para Svelte onMount.
 *
 * Didático: `gsap.context()` limita selectores ao sub-árvore do componente —
 * evita animar nós fora do componente e facilita `revert()` no unmount.
 */
export function runMotionScope(root: HTMLElement, setup: () => void): () => void {
  const ctx = gsap.context(setup, root);
  return () => ctx.revert();
}
