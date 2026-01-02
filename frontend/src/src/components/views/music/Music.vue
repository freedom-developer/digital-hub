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
              {{ allMusicList.length }}
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
                      :loading="favoriteLoading[row.id]"
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
import { ref, computed, onMounted, reactive } from 'vue'
import { VideoPlay, VideoPause, Headset, Star, StarFilled, List } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { musicApi } from '@/api/music'
import { usePlayerStore } from '@/stores/player'
import type { Music } from '@/api/music'

const playerStore = usePlayerStore()
const allMusicList = ref<Music[]>([])
const favoriteMusicList = ref<Music[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const favoriteLoading = reactive<Record<number, boolean>>({})

// 收藏的音乐ID集合
const favoriteIds = ref<Set<number>>(new Set())

// 菜单相关
const activeMenu = ref('all')

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)

// 根据当前菜单计算要显示的音乐列表
const displayMusicList = computed(() => {
  if (activeMenu.value === 'favorite') {
    return favoriteMusicList.value
  }
  return allMusicList.value
})

// 计算当前页显示的数据
const paginatedMusicList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return displayMusicList.value.slice(start, end)
})

onMounted(async () => {
  await Promise.all([
    fetchAllMusic(),
    fetchFavoriteIds(),
    fetchFavoriteMusic()
  ])
})

// 获取所有音乐
async function fetchAllMusic() {
  try {
    loading.value = true
    allMusicList.value = await musicApi.getMusicList()
    playerStore.setMusicList(allMusicList.value)
  } catch (err: any) {
    error.value = `获取音乐列表失败：${err.message}`
    console.error('错误详情:', err)
  } finally {
    loading.value = false
  }
}

// 获取收藏的音乐ID列表
async function fetchFavoriteIds() {
  try {
    const response = await musicApi.getFavoriteMusicIds()
    favoriteIds.value = new Set(response)
    console.log('收藏的音乐ID列表:', favoriteIds.value)
  } catch (err: any) {
    console.error('获取收藏列表失败:', err)
  }
}

// 获取收藏的音乐列表
async function fetchFavoriteMusic() {
  try {
    loading.value = true
    const response = await musicApi.getFavoriteMusic()
    favoriteMusicList.value = response
  } catch (err: any) {
    error.value = `获取收藏列表失败：${err.message}`
    console.error('错误详情:', err)
  } finally {
    loading.value = false
  }
}

// 切换收藏状态
async function toggleFavorite(music: Music) {
  favoriteLoading[music.id] = true
  
  try {
    if (favoriteIds.value.has(music.id)) {
      // 取消收藏
      await musicApi.removeFavorite(music.id)
      favoriteIds.value.delete(music.id)
      
      // 如果当前在收藏列表页，需要刷新列表
      // if (activeMenu.value === 'favorite') {
        // favoriteMusicList.value = favoriteMusicList.value.filter(m => m.id !== music.id)
      // }
      await fetchFavoriteMusic()
      
      ElMessage.success(`已取消收藏：${music.name}`)
    } else {
      // 添加收藏
      await musicApi.addFavorite(music.id)
      favoriteIds.value.add(music.id)
      
      // 如果当前在收藏列表页，需要刷新列表
      // if (activeMenu.value === 'favorite') {
      await fetchFavoriteMusic()
      // }
      
      ElMessage.success(`已收藏：${music.name}`)
    }
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '操作失败')
  } finally {
    favoriteLoading[music.id] = false
  }
}

// 判断是否已收藏
function isFavorite(id: number) {
  return favoriteIds.value.has(id)
}

// 菜单选择处理
async function handleMenuSelect(index: string) {
  activeMenu.value = index
  currentPage.value = 1
  
  // 切换到收藏列表时，重新获取数据
  if (index === 'favorite') {
    await fetchFavoriteMusic()
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

@media (max-width: 768px) {
  .music-container {
    flex-direction: column;
  }
  
  .menu-card {
    width: 100%;
  }
}
</style>
