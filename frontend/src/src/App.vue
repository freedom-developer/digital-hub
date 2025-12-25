<template>
  <div id="app">
    <Navbar 
      :currentMenu="currentMenu" 
      @menu-change="switchMenu"
    />
    <div class="content">
      <component :is="currentComponent" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import Navbar from './components/Navbar.vue'
import Home from './views/Home.vue'
import Music from './views/Music.vue'
import Movie from './views/Movie.vue'
import Computer from './views/Computer.vue'
import Math from './views/Math.vue'

export default defineComponent({
  name: 'App',
  components: {
    Navbar,
    Home,
    Music,
    Movie,
    Computer,
    Math
  },
  data() {
    return {
      currentMenu: 'home' as string
    }
  },
  computed: {
    currentComponent() {
      const componentMap: Record<string, string> = {
        home: 'Home',
        music: 'Music',
        movie: 'Movie',
        computer: 'Computer',
        math: 'Math'
      }
      return componentMap[this.currentMenu] || 'Home'
    }
  },
  methods: {
    switchMenu(menuId: string): void {
      this.currentMenu = menuId
    }
  }
})
</script>

<style>
  * {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
  background: #f5f5f5;
}

#app {
  min-height: 100vh;
}

/* 内容区域 */
.content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

/* 页面样式（所有页面共用） */
.page {
  background: white;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 2px 15px rgba(0, 0, 0, 0.08);
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page h1 {
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 3px solid #667eea;
}

.page p {
  line-height: 1.8;
  color: #666;
  margin: 15px 0;
  font-size: 16px;
}

.error {
  color: #dc3545;
  background: #f8d7da;
  padding: 15px;
  border-radius: 5px;
  border-left: 4px solid #dc3545;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page {
    padding: 25px;
  }
}
</style>
