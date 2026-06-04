/**
 * Detecção de prefers-reduced-motion para GSAP e CSS.
 *
 * > 💡 **Conceito:** utilizadores com distúrbios vestibulares podem pedir ao SO
 * que reduza animações — devemos honrar isso (WCAG, Apple HIG).
 */
export function prefersReducedMotion(): boolean {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return false;
  }
  return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

/** Duração efectiva: zero quando reduced motion está activo. */
export function motionDuration(seconds: number): number {
  return prefersReducedMotion() ? 0 : seconds;
}

/** Stagger efectivo: zero quando reduced motion está activo. */
export function motionStagger(seconds: number): number {
  return prefersReducedMotion() ? 0 : seconds;
}
