import Vue from "vue";
import VueRouter from "vue-router";

import Home from "@/views/Home.vue";
import Test from "@/views/Test.vue";

import Public from "@/views/Public.vue";
import Login from "@/views/Public/Login.vue";

import Organization from "@/views/Organization.vue";
import Blank from "@/views/Organization/Blank.vue";
import Integrations from "@/views/Organization/Integrations.vue";
import Rack from "@/views/Organization/Rack.vue";
import Racks from "@/views/Organization/Racks.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: Home,
    meta: { public: true },
  },
  {
    path: "/",
    component: Public,
    children: [
      {
        path: "/login",
        name: "login",
        component: Login,
        meta: { public: true },
      },
    ],
  },
  // {
  //   path: "/about",
  //   name: "About",
  //   // route level code-splitting
  //   // this generates a separate chunk (about.[hash].js) for this route
  //   // which is lazy-loaded when the route is visited.
  //   component: () =>
  //     import(/* webpackChunkName: "about" */ "../views/About.vue")
  // },
  {
    path: "/test",
    name: "Test",
    component: Test,
  },
  {
    path: "/organizations/:oid",
    name: "organization",
    component: Organization,
    children: [
      {
        path: "audits",
        name: "organization/audits",
        component: Blank,
        meta: { role: "operator" },
      },
      {
        path: "billing",
        name: "organization/billing",
        component: Blank,
        meta: { role: "administrator" },
      },
      {
        path: "integrations",
        name: "organization/integrations",
        component: Integrations,
        meta: { role: "operator" },
      },
      {
        path: "jobs",
        name: "organization/jobs",
        component: Blank,
        meta: { role: "developer" },
      },
      {
        path: "racks",
        name: "organization/racks",
        component: Racks,
        meta: { role: "developer" },
      },
      {
        path: "racks/:rid",
        name: "organization/rack",
        component: Rack,
        children: [
          {
            path: "apps",
            name: "organization/rack/apps",
            component: Blank,
            meta: { expand: true, role: "developer" },
          },
          {
            path: "instances",
            name: "organization/rack/instances",
            component: Blank,
            meta: { expand: true, role: "developer" },
          },
          {
            path: "resources",
            name: "organization/rack/resources",
            component: Blank,
            meta: { expand: true, role: "developer" },
          },
          {
            path: "settings",
            name: "organization/rack/settings",
            component: Blank,
            meta: { role: "developer" },
          },
          {
            path: "updates",
            name: "organization/rack/updates",
            component: Blank,
            meta: { expand: true, role: "developer" },
          },
        ],
      },
      {
        path: "settings",
        name: "organization/settings",
        component: Blank,
        meta: { role: "administrator" },
      },
      {
        path: "support",
        name: "organization/support",
        component: Blank,
        meta: { role: "developer" },
      },
      {
        path: "users",
        name: "organization/users",
        component: Blank,
        meta: { role: "administrator" },
      },
      {
        path: "workflows",
        name: "organization/workflows",
        component: Blank,
        meta: { role: "operator" },
      },
    ],
  },
];

const router = new VueRouter({
  base: process.env.BASE_URL,
  linkActiveClass: "active",
  mode: "history",
  routes,
});

export default router;

import { accessible } from "@/scripts/access";
import { createProvider } from "@/vue-apollo";
import store from "@/store";

const apollo = createProvider().defaultClient;

router.beforeEach(async (to, from, next) => {
  if (to.meta.public) {
    return next();
  }
  console.log("store.getters.token", store.getters.token);

  try {
    if (!to.params.oid) {
      return next({ name: "login" });
    }

    if (!store.getters.authenticated) {
      return next({ name: "login" });
    }

    const { organization } = (
      await apollo.query({
        query: require("@/queries/Organization.graphql"),
        variables: {
          id: to.params.oid,
        },
      })
    ).data;

    if (!accessible(to.meta.role, organization.access)) {
      return next({ name: "organization/racks", params: { oid: organization.id } });
    }

    next();
  } catch {
    next({ name: "login" });
  }
});
