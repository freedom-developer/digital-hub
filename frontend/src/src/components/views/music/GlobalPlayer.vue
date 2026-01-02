<template>
  <transition name="slide-up">
    <div v-if="playerStore.currentMusic" class="global-player">
      <el-card shadow="always" :body-style="{ padding: '16px' }">
        <!-- 使用 Element Plus 的 Row 和 Col 布局 -->
        <el-row :gutter="20" align="middle">
          <!-- 左侧：歌曲信息（固定宽度） -->
          <el-col :xs="24" :sm="8" :md="6">
            <el-space :size="12">
              <div class="album-cover">
                <div class="cover-bg">
                  <el-icon :size="32"><Headset /></el-icon>
                </div>
                <div v-if="playerStore.isPlaying" class="playing-animation">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
              
              <div class="music-info" style="width: 200px">
                <el-text class="music-title" truncated>{{ playerStore.currentMusic.name }}</el-text>
                <el-text type="info" size="small">
                  <el-icon><VideoPlay v-if="playerStore.isPlaying" /><VideoPause v-else /></el-icon>
                  {{ playerStore.isPlaying ? '正在播放' : '已暂停' }}
                </el-text>
              </div>
            </el-space>
          </el-col>

            <!-- 中间：播放控制和进度条（占据剩余空间，水平排列） -->
            <el-col :xs="24" :sm="10" :md="12">
                <div style="display: flex; align-items: center; gap: 16px; width: 100%;">
                    <!-- 播放控制按钮 -->
                    <el-space :size="8">
                        <!-- 播放模式按钮 -->
                        <el-button circle size="small" @click="playerStore.togglePlayMode()" :title="getPlayModeTitle()">
                          <el-icon>
                            <Refresh v-if="playerStore.playMode === 'order'" />
                            <Sort v-else-if="playerStore.playMode === 'random'" />
                            <RefreshLeft v-else />
                          </el-icon>
                        </el-button>

                        <el-button circle size="small" :disabled="!playerStore.hasPrevious" @click="playerStore.playPrevious()">
                            <el-icon><DArrowLeft /></el-icon>
                        </el-button>
                    
                        <el-button type="primary" circle @click="playerStore.togglePlay()">
                            <el-icon :size="20">
                            <VideoPlay v-if="!playerStore.isPlaying" />
                            <VideoPause v-else />
                            </el-icon>
                        </el-button>
            
                        <el-button circle size="small" :disabled="!playerStore.hasNext" @click="playerStore.playNext()">
                            <el-icon><DArrowRight /></el-icon>
                        </el-button>
                    </el-space>

                    <!-- 进度条 - 使用 flex: 1 占据剩余空间 -->
                    <div style="display: flex; align-items: center; gap: 8px; flex: 1;">
                        <el-text size="small" type="info">{{ formatTime(localCurrentTime) }}</el-text>
                        <el-slider
                            v-model="localCurrentTime"
                            :max="playerStore.duration"
                            :show-tooltip="true"
                            :format-tooltip="formatTime"
                            @input="onProgressInput"
                            @change="onProgressChange"
                            style="flex: 1; min-width: 0;"
                        />
                        <el-text size="small" type="info">{{ formatTime(playerStore.duration) }}</el-text>
                    </div>
                </div>
            </el-col>

          <!-- 右侧：音量控制（固定宽度） -->
          <el-col :xs="24" :sm="6" :md="6">
            <el-space :size="8" style="justify-content: flex-end; width: 100%">
              <el-button text circle @click="playerStore.toggleMute()">
                <el-icon>
                  <Microphone v-if="!(playerStore.volume === 0 || playerStore.isMuted)" />
                  <Mute v-else />
                </el-icon>
              </el-button>
              <el-slider
                v-model="localVolume"
                :max="100"
                :show-tooltip="true"
                :format-tooltip="(val: number) => `${val}%`"
                @change="onVolumeChange"
                style="width: 100px"
              />
            </el-space>
          </el-col>
        </el-row>
      </el-card>
    </div>
  </transition>

  <!-- 音频元素 -->
  <audio
    ref="audioPlayer"
    :src="playerStore.currentMusicUrl"
    @loadedmetadata="playerStore.handleLoadedMetadata"
    @timeupdate="playerStore.handleTimeUpdate"
    @ended="playerStore.handleEnded"
    @error="playerStore.handleError"
    @canplay="onCanPlay"
    preload="auto"
    style="display: none;"
  ></audio>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { 
  VideoPlay, 
  VideoPause, 
  DArrowLeft, 
  DArrowRight,
  Microphone,
  Mute,
  Headset,
  Refresh,
  Sort,
  RefreshLeft
} from '@element-plus/icons-vue'
import { PlayMode, usePlayerStore } from '@/stores/player'

const playerStore = usePlayerStore()
const audioPlayer = ref<HTMLAudioElement>()
const localCurrentTime = ref(0)
const localVolume = ref(50)
const isDraggingProgress = ref(false)

onMounted(() => {
  console.log('🎵 GlobalPlayer 组件已挂载')
  localVolume.value = playerStore.volume
  
  setTimeout(() => {
    if (audioPlayer.value) {
      console.log('✅ 音频元素已找到: ', audioPlayer.value)
      playerStore.setAudioElement(audioPlayer.value)
    } else {
      console.error('❌ 音频元素未找到')
    }
  }, 0)
})

watch(() => playerStore.currentMusic, (newMusic) => {
  if (newMusic) {
    console.log('🎵 当前音乐变化:', newMusic.name)
    setTimeout(() => {
      if (audioPlayer.value) {
        playerStore.setAudioElement(audioPlayer.value)
      }
    }, 0)
  }
})

watch(() => playerStore.currentTime, (newTime) => {
  if (!isDraggingProgress.value) {
    localCurrentTime.value = newTime
  }
})

watch(() => playerStore.volume, (newVolume) => {
  localVolume.value = newVolume
})

function onProgressInput(value: number) {
  console.log('🎯 拖动进度条到:', value)
  isDraggingProgress.value = true
  localCurrentTime.value = value
}

function onProgressChange(value: number) {
  console.log('✅ 进度条最终值:', value)
  isDraggingProgress.value = false
  localCurrentTime.value = value
  playerStore.seek(value)
}

function onVolumeChange(value: number) {
  console.log('🔊 音量改变到:', value)
  playerStore.setVolume(value)
}

function onCanPlay() {
  console.log('✅ 音频可以播放了')
}

function formatTime(seconds: number): string {
  if (isNaN(seconds)) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// 获取播放模式提示文字
function getPlayModeTitle() {
  const modeNames = {
    [PlayMode.ORDER]: '顺序播放',
    [PlayMode.RANDOM]: '随机播放',
    [PlayMode.LOOP]: '单曲循环'
  }
  return modeNames[playerStore.playMode]
}
</script>

<style scoped>
.global-player {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  padding: 0 20px 20px;
  background: linear-gradient(to top, rgba(255,255,255,0.95) 0%, transparent 100%);
  backdrop-filter: blur(10px);
}

/* 封面样式 */
.album-cover {
  position: relative;
  width: 50px;
  height: 50px;
  flex-shrink: 0;
}

.cover-bg {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

/* 播放动画 */
.playing-animation {
  position: absolute;
  bottom: 4px;
  right: 4px;
  display: flex;
  gap: 2px;
  align-items: flex-end;
  height: 12px;
}

.playing-animation span {
  width: 2px;
  background: white;
  border-radius: 2px;
  animation: playing 0.8s ease-in-out infinite;
}

.playing-animation span:nth-child(1) { animation-delay: 0s; }
.playing-animation span:nth-child(2) { animation-delay: 0.2s; }
.playing-animation span:nth-child(3) { animation-delay: 0.4s; }

@keyframes playing {
  0%, 100% { height: 4px; }
  50% { height: 12px; }
}

/* 音乐信息 */
.music-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.music-title {
  font-size: 14px;
  font-weight: 600;
}

/* 过渡动画 */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .global-player {
    padding: 0 10px 10px;
  }
}
</style>

