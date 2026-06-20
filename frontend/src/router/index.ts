import { createRouter, createWebHistory } from 'vue-router'
import AnalyzePage from '@/pages/AnalyzePage.vue'
import HistoryPage from '@/pages/HistoryPage.vue'
import AppLayout from '@/layouts/AppLayout.vue'

const routes = [
  {
    path: '/',
    component: AppLayout,
    children: [
      {
        path: '',
        name: 'analyze',
        component: AnalyzePage,
      },
      {
        path: 'history',
        name: 'history',
        component: HistoryPage,
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
