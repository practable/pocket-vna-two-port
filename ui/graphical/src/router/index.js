import { createWebHistory, createRouter } from 'vue-router';
import Parameters from '@/views/Parameters.vue';
import Calibration from '@/views/Calibration.vue';
import Verification from '@/views/Verification.vue';
import Measurement from '@/views/Measurement.vue';


const routes = [
    {
        path: '/',
        redirect: '/parameters'
    },
    {
        path: '/parameters',
        name: 'Parameters',
        component: Parameters
    },
  {
    path: '/calibration',
    name: 'Calibration',
    component: Calibration
  },
  {
    path: '/verification',
    name: 'Verification',
    component: Verification
  },
  {
    path: '/measurement',
    name: 'Measurement',
    component: Measurement
  },

  
  

]

const router = createRouter({
    history: createWebHistory(),   
    base: process.env.VITE_BASE,  
    routes,
})

export default router
