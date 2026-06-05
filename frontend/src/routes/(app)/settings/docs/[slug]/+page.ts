import { error } from "@sveltejs/kit";
import { loadDocPage } from "$lib/docs/loader";
import type { PageLoad } from "./$types";

export const load: PageLoad = ({ params }) => {
  const page = loadDocPage(params.slug);
  if (!page || !page.in_app) {
    error(404, "Documentação não encontrada");
  }
  return { page };
};
