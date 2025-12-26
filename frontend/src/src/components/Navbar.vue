<template>
  <el-menu
    :default-active="currentPath"
    mode="horizontal"
    :ellipsis="false"
    background-color="#667eea"
    text-color="#fff"
    active-text-color="#ffd04b"
    router
  >
    <el-menu-item index="/" class="brand">Digital Hub</el-menu-item>
    <div style="flex-grow: 1"></div>
    <el-menu-item index="/">首页</el-menu-item>
    <el-menu-item index="/music">音乐</el-menu-item>
    <el-menu-item index="/movie">电影</el-menu-item>
    <el-menu-item index="/computer">计算机</el-menu-item>
    <el-menu-item index="/math">数学</el-menu-item>
  </el-menu>

</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router';

export default defineComponent({
  name: 'Navbar',
  setup() {
    const route = useRoute()
    const currentPath = computed(() => route.path)
    return {
      currentPath
    }
  }
})
</script>

<style scoped>
.brand {
  font-size: 20px;
  font-weight: bold;
  pointer-events: none;
}
</style>
