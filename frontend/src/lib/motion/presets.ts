import { gsap } from "gsap";
import { MOTION } from "./config";
import { motionDuration, motionStagger, prefersReducedMotion } from "./reduced";

/** Entrada suave de painéis (auth card, modais). */
export function animatePanelEnter(target: gsap.TweenTarget): gsap.core.Tween {
  return gsap.from(target, {
    autoAlpha: 0,
    y: 12,
    duration: motionDuration(MOTION.panel.duration),
    ease: MOTION.panel.ease,
  });
}

/** Cross-fade entre secções (navegação app shell). */
export function animatePageEnter(target: gsap.TweenTarget): gsap.core.Tween {
  return gsap.fromTo(
    target,
    { autoAlpha: 0 },
    {
      autoAlpha: 1,
      duration: motionDuration(MOTION.panel.duration),
      ease: MOTION.panel.ease,
    },
  );
}

/** Fade simples para micro-transições. */
export function animateFadeIn(target: gsap.TweenTarget): gsap.core.Tween {
  return gsap.from(target, {
    autoAlpha: 0,
    duration: motionDuration(MOTION.micro.duration),
    ease: MOTION.micro.ease,
  });
}

/** Stagger de linhas de lista (cofre, tabelas). */
export function animateListStagger(
  container: HTMLElement,
  itemSelector = "[data-motion-item]",
): gsap.core.Tween | gsap.core.Timeline {
  const items = container.querySelectorAll(itemSelector);
  if (!items.length) {
    return gsap.timeline();
  }
  return gsap.from(items, {
    autoAlpha: 0,
    y: 8,
    duration: motionDuration(MOTION.list.duration),
    ease: MOTION.list.ease,
    stagger: motionStagger(MOTION.list.stagger),
  });
}

/** Highlight breve após guardar item (sem revelar passwords). */
export function animateSavedHighlight(target: gsap.TweenTarget): gsap.core.Tween {
  if (prefersReducedMotion()) {
    return gsap.set(target, { clearProps: "backgroundColor" });
  }
  return gsap.fromTo(
    target,
    { backgroundColor: "color-mix(in srgb, var(--color-success-fg) 18%, transparent)" },
    {
      backgroundColor: "transparent",
      duration: 0.55,
      ease: "power2.out",
      clearProps: "backgroundColor",
    },
  );
}
