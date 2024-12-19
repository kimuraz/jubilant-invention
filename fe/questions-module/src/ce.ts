import { defineCustomElement } from "vue";
import App from "./App.ce.vue";

const QuestionsModule = defineCustomElement(App);

customElements.define("questions-module", QuestionsModule);