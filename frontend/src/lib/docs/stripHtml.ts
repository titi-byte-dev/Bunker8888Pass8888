/** Remove tags HTML e normaliza espaços — útil para pesquisa e tooltips. */
export function stripHtml(html: string): string {
  return html
    .replace(/<[^>]+>/g, " ")
    .replace(/&nbsp;/g, " ")
    .replace(/&#(\d+);/g, (_, n) => String.fromCharCode(Number(n)))
    .replace(/&\w+;/g, " ")
    .replace(/\s+/g, " ")
    .trim();
}
