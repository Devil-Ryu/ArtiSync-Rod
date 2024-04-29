import { createRouter, createWebHistory } from "vue-router";
import Settings from "../views/Settings.vue";
import Platform from "../views/Platform.vue";
import Home from "../views/Home.vue";
import Accounts from "../views/Accounts.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "home", component: Home },
    { path: "/platform", name: "platform", component: Platform },
    { path: "/settings", name: "settings", component: Settings },
    { path: "/accounts", name: "accounts", component: Accounts },
  ],
});

export default router;
