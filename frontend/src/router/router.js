import { createRouter, createWebHistory } from "vue-router";
import Config from "../views/Config.vue";
import Platform from "../views/Platform.vue";
import Home from "../views/Home.vue";
import Accounts from "../views/Accounts.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "home", component: Home },
    { path: "/platform", name: "platform", component: Platform },
    { path: "/config", name: "config", component: Config },
    { path: "/accounts", name: "accounts", component: Accounts },
  ],
});

export default router;
