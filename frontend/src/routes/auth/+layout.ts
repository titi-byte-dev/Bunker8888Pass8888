import { redirect } from "@sveltejs/kit";
import { browser } from "$app/environment";
import { getAuthPhase, isPublicAppPath } from "$lib/auth/guard";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = ({ url }) => {
  if (!browser) return;

  const phase = getAuthPhase();
  const path = url.pathname;

  // Sessão activa + cofre desbloqueado → app
  if (phase === "unlocked") {
    redirect(302, "/vault");
  }

  // Token sem Master Key → unlock (excepto recovery)
  if (phase === "session") {
    const needsUnlock =
      path === "/auth" || path === "/auth/login" || path === "/auth/register";
    if (needsUnlock) {
      const redirectTo = encodeURIComponent("/vault");
      redirect(302, `/auth/unlock?redirect=${redirectTo}`);
    }
  }

  if (phase === "guest" && path === "/auth/unlock") {
    redirect(302, "/auth/login");
  }
};
