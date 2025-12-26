import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import Home  from "@/views/Home.vue";
import Music from "@/views/Music.vue";
import Movie from "@/views/Movie.vue";
import Computer from "@/views/Computer.vue";
import Math from "@/views/Math.vue";

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/music',
        name: 'Music',
        component: Music
    },
    {
        path: '/movie',
        name: 'Movie',
        component: Movie
    },
    {
        path: '/computer',
        name: 'Computer',
        component: Computer
    },
    {
        path: '/math',
        name: 'Math',
        component: Math
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router