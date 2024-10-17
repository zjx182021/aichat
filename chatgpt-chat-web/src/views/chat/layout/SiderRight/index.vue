<script setup lang="ts">
import type { CSSProperties } from 'vue'
import { computed, ref, watch } from 'vue'
import { NLayoutSider } from 'naive-ui'
import { useAppStore, useChatStore } from '@/store'
import { useBasicLayout } from '@/hooks/useBasicLayout'
import { PromptStore } from '@/components/common'

const appStore = useAppStore()
const chatStore = useChatStore()

const { isMobile } = useBasicLayout()
const show = ref(false)

const collapsed = computed(() => appStore.siderCollapsed)

function handleAdd() {
  chatStore.addHistory({ title: 'New Chat', uuid: Date.now(), isEdit: false })
  if (isMobile.value)
    appStore.setSiderCollapsed(true)
}

function handleUpdateCollapsed() {
  appStore.setSiderCollapsed(!collapsed.value)
}

const getMobileClass = computed<CSSProperties>(() => {
  if (isMobile.value) {
    return {
      position: 'fixed',
      zIndex: 50,
      width: 260,
    }
  }
  return {
  }
})

const getWidth = () => isMobile.value ? 280 : 500

const mobileSafeArea = computed(() => {
  if (isMobile.value) {
    return {
      paddingBottom: 'env(safe-area-inset-bottom)',
    }
  }
  return {}
})

watch(
  isMobile,
  (val) => {
    appStore.setSiderCollapsed(val)
  },
  {
    immediate: true,
    flush: 'post',
  },
)
</script>

<template>
  <NLayoutSider
    :collapsed="collapsed" :collapsed-width="0" collapse-mode="transform" position="static" :width="isMobile ? 260 : 500"
    bordered :style="getMobileClass" class="bottom-0 top-0" @update-collapsed="handleUpdateCollapsed"
  >
    <div class="flex flex-col h-full p-6 bottom-0" :style="mobileSafeArea">
      <main class="flex flex-col flex-1 min-h-0">
        <div class="_sidebar_79eze_17">
          <h3>
            声明：本项目主要是零声教育内部使用，探索AI与c/c++职业技能教育的结合。
          </h3>
          <br>

          <p>项目基于开源，感谢以下开源项目的贡献：</p>
          <p>后端开源：https://github.com/Arvintian/chatgpt-web.git</p>
          <p>前端开源：https://github.com/Chanzhaoyu/chatgpt-web</p>
          <p>敏感词汇过滤：https://github.com/importcjj/sensitive</p>
          <p>喂养方案：https://github.com/03Anmol/-openai_chatbot</p>
          <br>

          <p>云主机使用腾讯云（硅谷地区），chatgpt付费api消耗 $50/day，超出限量会有相关错误提示。 回复的观点，来源chatgpt官方，不代表开源项目，也不代表项目参与者观点。chatgpt属于ai，使用时，内容仅供参考。</p>
          <p>注: 目前本项目有为openai喂养一个关于c/c++开发开源技能点的zvoice_chatbot.json，会提高答案的可供参考的准确性。</p>
          <br>

          <p>
            零声，专注于C/C++，Linux，Nginx，ZeroMQ，MySQL，Redis，
            fastdfs，MongoDB，ZK，流媒体，CDN，P2P，K8S，Docker，
            TCP/IP，协程，DPDK, SPDK, bpf/ebpf等等相关技术探索与分享。
          </p>
          <br>

          <div>
            <div>零声交流群：<a href="https://jq.qq.com/?_wv=1027&amp;k=sEIbk6yO" target="blank">762073882</a></div>
            <div>每晚八点直播：<a href="https://ke.qq.com/course/417774?flowToken=1044591" target="_blank">https://ke.qq.com/course/417774</a></div>
            <div>ChatGPT项目实战教程：<a href="https://cox.xet.tech/s/16n7fK" target="_blank">https://cox.xet.tech/s/16n7fK</a></div>
          </div>
        </div>
      </main>
      <!-- <Footer /> -->
    </div>
  </NLayoutSider>
  <template v-if="isMobile">
    <div v-show="!collapsed" class="fixed inset-0 z-40 bg-black/40" @click="handleUpdateCollapsed" />
  </template>
  <PromptStore v-model:visible="show" />
</template>
