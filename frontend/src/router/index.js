import { createRouter, createWebHistory } from "vue-router";

// 定义路由
const routes = [
  {
    path: "/",
    name: "Home",
    component: () => import("../views/Home.vue"),
    meta: {
      title: "快速查询",
    },
  },
  {
    path: "/batch",
    name: "batch",
    component: () => import("../views/Batch.vue"),
    meta: {
      title: "批量处理",
    },
  },
  {
    path: "/management",
    name: "management",
    component: () => import("../views/Management.vue"),
    meta: {
      title: "管理数据",
    },
  },
  // 404 页面
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("../views/NotFound.vue"),
  },
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  },
});

router.beforeEach((to, from, next) => {
  document.title = to.meta.title;
  const metaTags = to.meta.metaTags || [];
  metaTags.forEach((tag) => {
    const tagElement = document.createElement("meta");
    Object.keys(tag).forEach((key) => {
      tagElement.setAttribute(key, tag[key]);
    });
    document.head.appendChild(tagElement);
  });
  next();
});

console.log("✅ 路由实例已创建，共有路由:", router.getRoutes().length, "条");

export default router;
