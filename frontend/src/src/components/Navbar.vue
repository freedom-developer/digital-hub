<template>
  <nav class="navbar">
    <div class="nav-brand">Digital Hub</div>
    <ul class="nav-menu">
      <li 
        v-for="menu in menus" 
        :key="menu.id"
        :class="{ active: currentMenu === menu.id }"
        @click="handleMenuClick(menu.id)"
      >
        {{ menu.name }}
      </li>
    </ul>
  </nav>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import type { MenuItem } from '@/types'

export default defineComponent({
  name: 'Navbar',
  props: {
    currentMenu: {
      type: String,
      required: true
    }
  },
  emits: ['menu-change'],
  data() {
    return {
      menus: [
        { id: 'home', name: '首页' },
        { id: 'music', name: '音乐' },
        { id: 'movie', name: '电影' },
        { id: 'computer', name: '计算机' },
        { id: 'math', name: '数学' }
      ] as MenuItem[]
    }
  },
  methods: {
    handleMenuClick(menuId: string): void {
      this.$emit('menu-change', menuId)
    }
  }
})
</script>

<style scoped>
  .navbar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-brand {
  font-size: 24px;
  font-weight: bold;
  padding: 18px 0;
}

.nav-menu {
  display: flex;
  list-style: none;
  gap: 5px;
}

.nav-menu li {
  padding: 18px 25px;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 5px;
  user-select: none;
}

.nav-menu li:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.nav-menu li.active {
  background: rgba(255, 255, 255, 0.3);
  font-weight: bold;
}

@media (max-width: 768px) {
  .navbar {
    flex-direction: column;
    padding: 10px 20px;
  }

  .nav-menu {
    flex-wrap: wrap;
    justify-content: center;
  }

  .nav-menu li {
    padding: 12px 15px;
  }
}
</style>
