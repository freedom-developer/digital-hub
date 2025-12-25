<template>
  <div id="app">
    <h1>用户信息</h1>
    <div v-if="loading">加载中...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else class="user-info">
      <p><strong>姓名: </strong>{{ user.name }}</p>
      <p><strong>年龄: </strong>{{ user.age }}</p>
    </div>
  </div>
</template>>

<script>
  import axios from "axios";
  export  default {
    name: 'App',
    data() {
      return {
        user: {},
        loading: true,
        error: null
      };
    },
    mounted() {
      this.fetchUser();
    },
    methods: {
      async fetchUser() {
        try {
          const response = await axios.get('http://localhost:8888/api/user');
          this.user = response.data;
          this.loading = false;
        } catch (err) {
          this.error = '获取用户信息失败: ' + err.message;
          this.loading = false;
        }
      }
    }
  };
</script>

<style>
  #app {
    font-family: Arial, sans-serif;
    max-width: 600px;
    margin: 50px auto;
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 8px;
  }

  h1 {
    color: #333;
  }

  .user-info {
    background-color: #f9f9f9;
    padding: 15px;
    border-radius: 5px;
    margin-top: 20px;
  }

  p {
    font-size: 18px;
    margin: 10px 0;
  }
</style>