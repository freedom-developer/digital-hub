<template>
  <div class="page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon :size="24"><Headset /></el-icon>
          <span>我的音乐</span>
          <el-tag type="info" size="small">{{ musicList.length }} 首歌曲</el-tag>
        </div>
      </template>

      <div v-loading="loading" style="min-height: 200px;">
        <el-alert
          v-if="error"
          :title="error"
          type="error"
          :closable="false"
          style="margin-bottom: 20px;"
        />

        <div v-else>
          <el-table
            :data="paginatedMusicList"
            stripe
            style="width: 100%;"
            :row-class-name="getRowClassName"
            @row-click="handleRowClick"
          >
            <el-table-column type="index" label="序号" width="80" align="center">
              <template #default="{ $index }">
                {{ (currentPage - 1) * pageSize + $index + 1 }}
              </template>
            </el-table-column>
            
            <el-table-column prop="name" label="歌曲名" min-width="100">
              <template #default="{ row }">
                <el-space>
                  <el-icon v-if="playerStore.currentMusic?.id === row.id && playerStore.isPlaying" color="#67c23a" class="playing-wave">
                    <VideoPlay />
                  </el-icon>
                  <span :style="{ color: playerStore.currentMusic?.id === row.id ? '#409eff' : '', fontWeight: playerStore.currentMusic?.id === row.id ? 'bold' : '' }">
                    {{ row.name }}
                  </span>
                </el-space>
              </template>
            </el-table-column>

            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button
                  :type="playerStore.currentMusic?.id === row.id && playerStore.isPlaying ? 'success' : 'primary'"
                  size="small"
                  circle
                  @click.stop="handlePlay(row)"
                >
                  <el-icon>
                    <VideoPlay v-if="!(playerStore.currentMusic?.id === row.id && playerStore.isPlaying)" />
                    <VideoPause v-else />
                  </el-icon>
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页组件 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :small="false"
              :background="true"
              layout="total, sizes, prev, pager, next, jumper"
              :total="musicList.length"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>

        <el-empty v-if="!loading && !error && musicList.length === 0" description="暂无音乐" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { VideoPlay, VideoPause, Headset } from '@element-plus/icons-vue'
import { musicApi } from '@/api/music'
import { usePlayerStore } from '@/stores/player'
import type { Music } from '@/types'

const playerStore = usePlayerStore()
const musicList = ref<Music[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

// 计算当前页显示的数据
const paginatedMusicList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return musicList.value.slice(start, end)
})

onMounted(async () => {
  await fetchMusicList()
})

async function fetchMusicList() {
  try {
    loading.value = true
    musicList.value = await musicApi.getMusicList()
    playerStore.setMusicList(musicList.value)
    loading.value = false
  } catch (err: any) {
    error.value = `获取音乐列表失败：${err.message}`
    loading.value = false
    console.error('错误详情:', err)
  }
}

function handleRowClick(row: Music) {
  handlePlay(row)
}

function handlePlay(music: Music) {
  if (playerStore.currentMusic?.id === music.id) {
    playerStore.togglePlay()
  } else {
    playerStore.play(music)
  }
}

function getRowClassName({ row }: { row: Music }) {
  return playerStore.currentMusic?.id === row.id ? 'current-row' : ''
}

function handleSizeChange(val: number) {
  pageSize.value = val
  currentPage.value = 1 // 改变每页数量时回到第一页
}

function handleCurrentChange(val: number) {
  currentPage.value = val
}
</script>

<style scoped>
.page {
  max-width: 1000px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: bold;
}

.playing-wave {
  animation: wave 1s ease-in-out infinite;
}

@keyframes wave {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(1.3); }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 10px 0;
}

:deep(.current-row) {
  background-color: #ecf5ff !important;
}

:deep(.current-row:hover > td) {
  background-color: #d9ecff !important;
}
</style>

