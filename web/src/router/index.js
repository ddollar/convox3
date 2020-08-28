import Vue from "vue";
import VueRouter from "vue-router";

import Blank from "@/views/Private/Organization/Blank.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: () => import(/* webpackChunkName: "home" */ "@/views/Home.vue"),
    meta: { public: true },
  },
  {
    path: "/",
    component: () => import(/* webpackChunkName: "public" */ "@/views/Public.vue"),
    children: [
      {
        path: "/login",
        name: "login",
        component: () => import(/* webpackChunkName: "public/login" */ "@/views/Public/Login.vue"),
        meta: { public: true },
      },
    ],
  },
  {
    path: "/",
    name: "private",
    component: () => import(/* webpackChunkName: "private" */ "@/views/Private.vue"),
    children: [
      {
        path: "/account",
        name: "account",
        component: () => import(/* webpackChunkName: "public/login" */ "@/views/Private/Account.vue"),
      },
      {
        path: "/organizations/:oid",
        name: "organization",
        component: () => import(/* webpackChunkName: "organization" */ "@/views/Private/Organization.vue"),
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
            component: () => import(/* webpackChunkName: "organization/integrations" */ "@/views/Private/Organization/Integrations.vue"),
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
            component: () => import(/* webpackChunkName: "organization/racks" */ "@/views/Private/Organization/Racks.vue"),
            meta: { role: "developer" },
          },
          {
            path: "racks/:rid",
            name: "organization/rack",
            component: () => import(/* webpackChunkName: "organization/rack" */ "@/views/Private/Organization/Rack.vue"),
            meta: { role: "developer" },
            children: [
              {
                path: "apps",
                name: "organization/rack/apps",
                component: () => import(/* webpackChunkName: "organization/rack/apps" */ "@/views/Private/Organization/Rack/Apps.vue"),
                meta: { expand: true, role: "developer" },
              },
              {
                path: "instances",
                name: "organization/rack/instances",
                component: () =>
                  import(/* webpackChunkName: "organization/rack/instances" */ "@/views/Private/Organization/Rack/Instances.vue"),
                meta: { expand: true, role: "developer" },
              },
              {
                path: "logs",
                name: "organization/rack/logs",
                component: () => import(/* webpackChunkName: "organization/rack/logs" */ "@/views/Private/Organization/Rack/Logs.vue"),
                meta: { expand: true, role: "developer" },
              },
              {
                path: "processes",
                name: "organization/rack/processes",
                component: () =>
                  import(/* webpackChunkName: "organization/rack/processes" */ "@/views/Private/Organization/Rack/Processes.vue"),
                meta: { expand: true, role: "developer" },
              },
              {
                path: "resources",
                name: "organization/rack/resources",
                component: () =>
                  import(/* webpackChunkName: "organization/rack/resources" */ "@/views/Private/Organization/Rack/Resources.vue"),
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
                component: () =>
                  import(/* webpackChunkName: "organization/rack/updates" */ "@/views/Private/Organization/Rack/Updates.vue"),
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
            component: () => import(/* webpackChunkName: "organization/users" */ "@/views/Private/Organization/Users.vue"),
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
  if (to.meta.public) return next();
  if (!store.getters.authenticated) return next({ name: "login" });
  if (!to.meta.role) return next();

  try {
    if (!to.params.oid) return next({ name: "login" });

    const { organization } = (
      await apollo.query({
        query: require("@/queries/Organization.graphql"),
        variables: {
          id: to.params.oid,
        },
      })
    ).data;

    if (!accessible(to.meta.role, organization.access)) {
      return next({
        name: "organization/racks",
        params: { oid: organization.id },
      });
    }

    next();
  } catch {
    next({ name: "login" });
  }
});
