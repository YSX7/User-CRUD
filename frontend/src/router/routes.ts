import MainLayoutVue from 'layouts/MainLayout.vue';
import RedirectBeachVue from 'pages/RedirectBeach.vue';
import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ name:'login', path: '', component: () => import('pages/IndexPage.vue') }],
  },
  {
    path: '/redirectBeach',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ name: 'index', path: '', component: () => import('pages/RedirectBeach.vue')   }],
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  }
];

export default routes;
