/**
 * Presets de motion AegisPass (UI-005).
 *
 * Didático: centralizamos duração e easing para manter consistência em toda
 * a app — ver docs/roadmap/09-design/product-ui-vision.md §6.
 */
export const MOTION = {
  /** Hover, toggle — 120–180ms */
  micro: { duration: 0.15, ease: "power2.out" },
  /** Painéis e páginas — 280–380ms */
  panel: { duration: 0.32, ease: "power3.out" },
  /** Modal enter */
  modalEnter: { duration: 0.32, ease: "back.out(1.2)" },
  /** Modal exit */
  modalExit: { duration: 0.22, ease: "power2.in" },
  /** Lista com stagger — 30–50ms/item */
  list: { duration: 0.28, stagger: 0.04, ease: "power1.out" },
} as const;

export type MotionPreset = keyof typeof MOTION;
