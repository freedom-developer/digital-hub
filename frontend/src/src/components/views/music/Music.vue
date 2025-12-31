<template>
  <div class="page">
    <div class="music-container">
      <!-- 左侧菜单 -->
      <el-card class="menu-card" shadow="hover">
        <el-menu
          :default-active="activeMenu"
          class="music-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="all">
            <el-icon><List /></el-icon>
            <span>所有音乐</span>
            <el-tag size="small" type="info" style="margin-left: auto;">
              {{ musicList.length }}
            </el-tag>
          </el-menu-item>
          
          <el-menu-item index="favorite">
            <el-icon><Star /></el-icon>
            <span>我的收藏</span>
            <el-tag size="small" type="warning" style="margin-left: auto;">
              {{ favoriteMusicList.length }}
            </el-tag>
          </el-menu-item>
        </el-menu>
      </el-card>

      <!-- 右侧音乐列表 -->
      <el-card class="content-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <el-icon :size="24"><Headset /></el-icon>
            <span>{{ activeMenu === 'all' ? '所有音乐' : '我的收藏' }}</span>
            <el-tag type="info" size="small">{{ displayMusicList.length }} 首歌曲</el-tag>
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

              <el-table-column label="操作" width="180" align="center">
                <template #default="{ row }">
                  <el-space>
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
                    
                    <el-button
                      :type="isFavorite(row.id) ? 'warning' : 'default'"
                      size="small"
                      circle
                      @click.stop="toggleFavorite(row)"
                    >
                      <el-icon>
                        <StarFilled v-if="isFavorite(row.id)" />
                        <Star v-else />
                      </el-icon>
                    </el-button>
                  </el-space>
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
                :total="displayMusicList.length"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
          </div>

          <el-empty v-if="!loading && !error && displayMusicList.length === 0" description="暂无音乐" />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { VideoPlay, VideoPause, Headset, Star, StarFilled, List } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { musicApi } from '@/api/music'
import { usePlayerStore } from '@/stores/player'
import type { Music } from '@/types'

const playerStore = usePlayerStore()
const musicList = ref<Music[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

// 菜单相关
const activeMenu = ref('all')
const favoriteIds = ref<Set<number>>(new Set())

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

// 根据当前菜单计算要显示的音乐列表
const displayMusicList = computed(() => {
  if (activeMenu.value === 'favorite') {
    return musicList.value.filter(music => favoriteIds.value.has(music.id))
  }
  return musicList.value
})

// 收藏的音乐列表
const favoriteMusicList = computed(() => {
  return musicList.value.filter(music => favoriteIds.value.has(music.id))
})

// 计算当前页显示的数据
const paginatedMusicList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return displayMusicList.value.slice(start, end)
})

onMounted(async () => {
  await fetchMusicList()
  loadFavorites()
})

async function fetchMusicList() {
  try {
    loading.value = true
    musicList.value = await musicApi.getMusicList()
    console.log('获取到的音乐列表:', musicList.value)
    playerStore.setMusicList(musicList.value)
    loading.value = false
  } catch (err: any) {
    error.value = `获取音乐列表失败：${err.message}`
    loading.value = false
    console.error('错误详情:', err)
  }
}

// 加载收藏列表（从 localStorage）
function loadFavorites() {
  const saved = localStorage.getItem('favoriteMusicIds')
  if (saved) {
    try {
      const ids = JSON.parse(saved)
      favoriteIds.value = new Set(ids)
    } catch (e) {
      console.error('加载收藏列表失败:', e)
    }
  }
}

// 保存收藏列表到 localStorage
function saveFavorites() {
  localStorage.setItem('favoriteMusicIds', JSON.stringify([...favoriteIds.value]))
}

// 切换收藏状态
function toggleFavorite(music: Music) {
  if (favoriteIds.value.has(music.id)) {
    favoriteIds.value.delete(music.id)
    ElMessage.success(`已取消收藏：${music.name}`)
  } else {
    favoriteIds.value.add(music.id)
    ElMessage.success(`已收藏：${music.name}`)
  }
  saveFavorites()
}

// 判断是否已收藏
function isFavorite(id: number) {
  return favoriteIds.value.has(id)
}

// 菜单选择处理
function handleMenuSelect(index: string) {
  activeMenu.value = index
  currentPage.value = 1 // 切换菜单时回到第一页
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
  currentPage.value = 1
}

function handleCurrentChange(val: number) {
  currentPage.value = val
}
</script>

<style scoped>
.page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

.music-container {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.menu-card {
  width: 200px;
  flex-shrink: 0;
}

.content-card {
  flex: 1;
  min-width: 0;
}

.music-menu {
  border: none;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .music-container {
    flex-direction: column;
  }
  
  .menu-card {
    width: 100%;
  }
}
</style>
