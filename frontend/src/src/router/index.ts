import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import Home  from "@/components/views/home/Home.vue";
import Music from "@/components/views/music/Music.vue";
import Movie from "@/components/views/movie/Movie.vue";
import Computer from "@/components/views/computer/Computer.vue";
import Math from "@/components/views/math/Math.vue";

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