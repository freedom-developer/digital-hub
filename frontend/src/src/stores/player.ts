import { defineStore } from 'pinia'
import { ref, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { Music } from '@/types'

// 定义播放模式类型
export enum PlayMode {
  ORDER = 'order',      // 顺序播放
  RANDOM = 'random',    // 随机播放
  LOOP = 'loop'         // 单曲循环
}

export const usePlayerStore = defineStore('player', () => {
  // 状态
  const currentMusic = ref<Music | null>(null)
  const isPlaying = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const volume = ref(50)
  const isMuted = ref(false)
  const musicList = ref<Music[]>([])
  const playMode = ref<PlayMode>(PlayMode.ORDER) // 默认顺序播放
  
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
    // console.log('设置音频元素:', audio)
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
    console.log('当前音乐播放结束，播放下一首')
    isPlaying.value = false
    const nextMusic = getNextMusic()
    if (nextMusic) {
      play(nextMusic)
    } else {
      ElMessage.info('播放列表已结束')
    }
    // if (hasNext.value) {
    //   playNext()
    // } else {
    //   ElMessage.info('播放列表已结束')
    // }
  }

  function handleError(event: Event) {
    console.error('音频加载错误:', event)
    ElMessage.error('音频加载失败，请检查文件是否存在')
    isPlaying.value = false
  }

  // 切换播放模式
  function togglePlayMode() {
    const modes = [PlayMode.ORDER, PlayMode.RANDOM, PlayMode.LOOP]
    const currentModeIndex = modes.indexOf(playMode.value)
    const nextModeIndex = (currentModeIndex + 1) % modes.length
    playMode.value = modes[nextModeIndex]
    
    const modeNames = {
      [PlayMode.ORDER]: '顺序播放',
      [PlayMode.RANDOM]: '随机播放',
      [PlayMode.LOOP]: '单曲循环'
    }
    
    ElMessage.success({
      message: `切换到${modeNames[playMode.value]}`,
      grouping: true
    })
  }

  // 根据播放模式获取下一首歌
  function getNextMusic(): Music | null {
    if (musicList.value.length === 0) return null
    console.log('当前播放模式:', playMode.value)
    switch (playMode.value) {
      case PlayMode.LOOP:
        // 单曲循环，返回当前歌曲
        return currentMusic.value
        
      case PlayMode.RANDOM:
        // 随机播放
        const randomIndex = Math.floor(Math.random() * musicList.value.length)
        return musicList.value[randomIndex]
        
      case PlayMode.ORDER:
      default:
        // 顺序播放
        if (hasNext.value) {
          return musicList.value[currentIndex.value + 1]
        }
        return null
    }
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
    playMode,
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
    togglePlayMode,
    seek,
    setVolume,
    toggleMute,
    handleTimeUpdate,
    handleLoadedMetadata,
    handleEnded,
    handleError
  }
})
