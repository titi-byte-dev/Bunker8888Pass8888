import { redirect } from "@sveltejs/kit";
import { browser } from "$app/environment";
import { getAuthPhase, isPublicAppPath } from "$lib/auth/guard";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = ({ url }) => {
  if (!browser) return;

  if (isPublicAppPath(url.pathname)) return;

  const phase = getAuthPhase();

  if (phase === "guest") {
    const redirectTo = encodeURIComponent(url.pathname + url.search);
    redirect(302, `/auth/login?redirect=${redirectTo}`);
  }

  if (phase === "session") {
    const redirectTo = encodeURIComponent(url.pathname + url.search);
    redirect(302, `/auth/unlock?redirect=${redirectTo}`);
  }
};
