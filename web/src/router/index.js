import Vue from "vue";
import VueRouter from "vue-router";

import Home from "@/views/Home.vue";
import Organization from "@/views/Organization.vue";
import Test from "@/views/Test.vue";

import Dashboard from "@/views/Organization/Dashboard.vue";
import Racks from "@/views/Organization/Racks.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: Home
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
    component: Test
  },
  {
    path: "/organizations/:oid",
    name: "organization",
    component: Organization,
    children: [
      {
        path: "audits",
        name: "organization/audits",
        component: Dashboard,
        meta: { role: "operator" }
      },
      {
        path: "billing",
        name: "organization/billing",
        component: Dashboard,
        meta: { role: "administrator" }
      },
      {
        path: "dashboard",
        name: "organization/dashboard",
        component: Dashboard
      },
      {
        path: "integrations",
        name: "organization/integrations",
        component: Dashboard,
        meta: { role: "operator" }
      },
      {
        path: "jobs",
        name: "organization/jobs",
        component: Dashboard
      },
      {
        path: "racks",
        name: "organization/racks",
        component: Racks
      },
      {
        path: "racks/:rid",
        name: "organization/rack",
        component: Racks,
        children: [
          {
            path: "apps",
            name: "organization/rack/apps",
            component: Racks,
            meta: { expand: true }
          },
          {
            path: "instances",
            name: "organization/rack/instances",
            component: Racks,
            meta: { expand: true }
          },
          {
            path: "resources",
            name: "organization/rack/resources",
            component: Racks,
            meta: { expand: true }
          },
          {
            path: "settings",
            name: "organization/rack/settings",
            component: Racks
          },
          {
            path: "updates",
            name: "organization/rack/updates",
            component: Racks,
            meta: { expand: true }
          }
        ]
      },
      {
        path: "settings",
        name: "organization/settings",
        component: Dashboard,
        meta: { role: "administrator" }
      },
      {
        path: "support",
        name: "organization/support",
        component: Dashboard
      },
      {
        path: "users",
        name: "organization/users",
        component: Dashboard,
        meta: { role: "administrator" }
      },
      {
        path: "workflows",
        name: "organization/workflows",
        component: Dashboard,
        meta: { role: "operator" }
      }
    ]
  }
];

const router = new VueRouter({
  base: process.env.BASE_URL,
  linkActiveClass: "active",
  mode: "history",
  routes
});

export default router;

import { accessible } from "@/scripts/access";
import { createProvider } from "@/vue-apollo";

const apollo = createProvider().defaultClient;

router.beforeEach(async (to, from, next) => {
  const { role } = to.meta;

  if (!to.params.oid) {
    next();
  }

  const { organization } = (
    await apollo.query({
      query: require("@/queries/Organization.graphql"),
      variables: {
        id: to.params.oid
      }
    })
  ).data;

  if (!accessible(role, organization.access)) {
    next({ name: "organization/dashboard", params: { oid: organization.id } });
  }

  next();
});
