/**
 * Tipos partilhados do DataTable (UI-015).
 * Ficheiro .ts separado para export limpo em index.ts e testes.
 */

export type DataColumn<T> = {
  id: string;
  label: string;
  align?: "left" | "right";
  mono?: boolean;
  muted?: boolean;
  accessor?: (row: T) => string;
};
