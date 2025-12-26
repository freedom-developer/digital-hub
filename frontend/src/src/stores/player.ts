import { defineStore } from 'pinia'
import { ref, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { Music } from '@/types'

export const usePlayerStore = defineStore('player', () => {
  // 状态
  const currentMusic = ref<Music | null>(null)
  const isPlaying = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const volume = ref(50)
  const isMuted = ref(false)
  const musicList = ref<Music[]>([])
  
  // 音频元素引用
  let audioElement: HTMLAudioElement | null = null

  // 计算属性
  const currentIndex = computed(() => {
    if (!currentMusic.value) return -1
    return musicList.value.findIndex(m => m.id === currentMusic.value!.id)
  })

  const hasPrevious = computed(() => currentIndex.value > 0)
  const hasNext = computed(() => currentIndex.value < musicList.value.length - 1)

  const currentMusicUrl = computed(() => {
    if (!currentMusic.value) return ''
    return `/api/music/download/${currentMusic.value.id}`
  })

  // 方法
  function setAudioElement(audio: HTMLAudioElement) {
    console.log('设置音频元素:', audio)
    audioElement = audio
    audio.volume = volume.value / 100
  }

  function setMusicList(list: Music[]) {
    musicList.value = list
  }

  async function play(music?: Music) {
    if (music) {
      // 播放新音乐
      console.log('准备播放:', music.name)
      currentMusic.value = music
      isPlaying.value = false
      
      await nextTick()
      
      if (audioElement) {
        console.log('音频元素已就绪，开始播放')
        audioElement.load()
        audioElement.volume = volume.value / 100
        
        try {
          await audioElement.play()
          isPlaying.value = true
          ElMessage.success({
            message: `正在播放：${currentMusic.value?.name}`,
            grouping: true
          })
          console.log('播放成功')
        } catch (err: any) {
          ElMessage.error(`播放失败：${err.message}`)
          console.error('播放错误:', err)
          isPlaying.value = false
        }
      } else {
        console.error('音频元素未初始化')
        ElMessage.error('播放器未就绪，请刷新页面')
      }
    } else {
      // 继续播放当前音乐
      if (audioElement && currentMusic.value) {
        try {
          await audioElement.play()
          isPlaying.value = true
        } catch (err: any) {
          ElMessage.error(`播放失败：${err.message}`)
          console.error('播放错误:', err)
        }
      }
    }
  }

  function pause() {
    if (audioElement) {
      audioElement.pause()
      isPlaying.value = false
    }
  }

  function togglePlay() {
    if (isPlaying.value) {
      pause()
    } else {
      play()
    }
  }

  function playPrevious() {
    if (hasPrevious.value) {
      const prevMusic = musicList.value[currentIndex.value - 1]
      play(prevMusic)
    }
  }

  function playNext() {
    if (hasNext.value) {
      const nextMusic = musicList.value[currentIndex.value + 1]
      play(nextMusic)
    }
  }

  function seek(time: number) {
        console.log('🔍 Seek 到:', time)
    if (audioElement) {
      try {
        audioElement.currentTime = time
        currentTime.value = time
        console.log('✅ Seek 成功，当前时间:', audioElement.currentTime)
      } catch (err) {
        console.error('❌ Seek 失败:', err)
      }
    } else {
      console.error('❌ 音频元素不存在，无法 seek')
    }
  }

  function setVolume(val: number) {
    volume.value = val
    if (audioElement) {
      audioElement.volume = val / 100
      isMuted.value = false
    }
  }

  function toggleMute() {
    if (audioElement) {
      if (isMuted.value) {
        console.log('恢复音量')
        audioElement.volume = volume.value / 100
        isMuted.value = false
      } else {
        console.log('设置静音')
        audioElement.volume = 0
        isMuted.value = true
      }
    } else {
        console.log('没有音乐元素')
    }
  }

  function handleTimeUpdate() {
    if (audioElement) {
      currentTime.value = audioElement.currentTime
    }
  }

  function handleLoadedMetadata() {
    if (audioElement) {
      duration.value = audioElement.duration
      console.log('音频已加载，时长:', duration.value)
    }
  }

  function handleEnded() {
    isPlaying.value = false
    if (hasNext.value) {
      playNext()
    } else {
      ElMessage.info('播放列表已结束')
    }
  }

  function handleError(event: Event) {
    console.error('音频加载错误:', event)
    ElMessage.error('音频加载失败，请检查文件是否存在')
    isPlaying.value = false
  }

  return {
    // 状态
    currentMusic,
    isPlaying,
    currentTime,
    duration,
    volume,
    isMuted,
    musicList,
    // 计算属性
    currentIndex,
    hasPrevious,
    hasNext,
    currentMusicUrl,
    // 方法
    setAudioElement,
    setMusicList,
    play,
    pause,
    togglePlay,
    playPrevious,
    playNext,
    seek,
    setVolume,
    toggleMute,
    handleTimeUpdate,
    handleLoadedMetadata,
    handleEnded,
    handleError
  }
})
