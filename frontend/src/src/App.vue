<template>
  <div id="app">
    <Navbar />
    <div class="content">
    <el-main>
      <router-view />
    </el-main>
    </div>
    <!-- 全局播放器 - 始终存在 -->
    <GlobalPlayer />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import Navbar from './components/Navbar.vue'
import GlobalPlayer from './components/GlobalPlayer.vue';

export default defineComponent({
  name: 'App',
  components: {
    Navbar,
    GlobalPlayer
  }
})
</script>

<style>
#app {
  min-height: 100vh;
  padding-bottom: 160px; /* 为播放器留出空间 */
}
</style>
