import { RouteRecordRaw } from 'vue-router';

const MainLayoutComponent = () => import('layouts/MainLayout.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: MainLayoutComponent,
    children: [
      { name: 'login', path: '/login', component: () => import('src/pages/LoginPage.vue'), meta: { requiresAuth: false } },
      { name: 'index', path: '/', component: () => import('src/pages/IndexPage.vue'), meta: { requiresAuth: true } },
      { name: 'users', path: '/users', component: () => import('src/pages/UsersPage.vue'), meta: { requiresAuth: true } },]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  }
];

export default routes;
