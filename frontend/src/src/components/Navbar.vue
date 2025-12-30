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

    <!-- '我的' 下拉菜单 -->
    <div class="menu-dropdown-wrapper">
      <el-dropdown @command="handleCommand">
        <span class="el-dropdown-link" style="color: #fff; cursor: pointer; ">
          我的 <el-icon class="el-icon--right"><arrow-down /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="/user/register">注册</el-dropdown-item>
            <el-dropdown-item command="/user/login">登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

  </el-menu>

</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useRoute, useRouter } from 'vue-router';
import { ArrowDown } from '@element-plus/icons-vue';

export default defineComponent({
  name: 'Navbar',
  components: {
    ArrowDown
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const currentPath = computed(() => route.path)
    const handleCommand = (command: string) => {
      router.push(command)
    };

    return {
      currentPath,
      handleCommand
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
.menu-dropdown-wrapper {
  display: flex;
  align-items: center;
  height: 100%; /* 让容器高度与菜单栏一致 */
  padding: 0 20px;
}

.el-dropdown-link {
  display: flex;
  align-items: center;
  line-height: normal; /* 重置行高 */
}
</style>
