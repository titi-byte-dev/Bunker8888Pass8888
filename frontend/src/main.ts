// Ponto de entrada do frontend. Em Svelte 5 montamos a app com `mount`.
import { mount } from "svelte";
import App from "./App.svelte";
import "./lib/design/tokens.css";
import { initTheme } from "./lib/design";

initTheme();

const target = document.getElementById("app");
if (!target) {
  throw new Error("Elemento #app não encontrado no index.html");
}

const app = mount(App, { target });

export default app;
