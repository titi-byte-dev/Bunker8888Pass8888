const STORAGE_KEY = "aegispass-sidebar-collapsed";

/** Lê preferência de sidebar colapsada (só desktop). */
export function loadSidebarCollapsed(): boolean {
  if (typeof localStorage === "undefined") return false;
  return localStorage.getItem(STORAGE_KEY) === "1";
}

export function saveSidebarCollapsed(collapsed: boolean): void {
  if (typeof localStorage === "undefined") return;
  localStorage.setItem(STORAGE_KEY, collapsed ? "1" : "0");
}
