import { redirect } from "@sveltejs/kit";

/** Raiz da app → módulo Cofre por omissão */
export function load() {
  redirect(302, "/vault");
}
