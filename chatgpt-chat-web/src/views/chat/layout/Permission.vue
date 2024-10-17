<script setup lang='ts'>
import { computed, ref } from 'vue'
import { NButton, NInput, NModal, useMessage } from 'naive-ui'
import { useAuthStore } from '@/store'
import Icon403 from '@/icons/403.vue'
import { fetchCode, login } from '@/api'

interface Props {
  visible: boolean
}

defineProps<Props>()

const authStore = useAuthStore()
const ms = useMessage()
authStore.setToken('helloxx')
// authStore.getSession().then((data) => {
// })
// console.log('authStore.getSession() ', authStore.getSession())
const phoneRe = /1\d{10}$/
const codeRe = /\d{4}/

const token = ref('')
const loading = ref(false)
const phone = ref('')
const code = ref('')
const codeBtnMsg = ref('发送验证码')
const codeSended = ref(false)
const codeLoading = ref(false)

const sendCodeDisabled = computed(() => !phoneRe.test(phone.value.trim()) || codeSended.value === true)
const loginDisabled = computed(() => !phoneRe.test(phone.value.trim()) || !codeRe.test(code.value.trim()))

async function sendCode() {
  fetchCode(phone.value).then((data) => {
    ms.success('验证码已发送')
  }).catch((err) => {
    // console.log('err: ', err)
    ms.error(err.error)
  })

  let time = 120
  codeSended.value = true
  codeBtnMsg.value = `已发送 ( ${time} )`

  const intervalId = setInterval(() => {
    codeBtnMsg.value = `已发送 ( ${--time} )`
    if (time === 0) {
      codeBtnMsg.value = '发送验证码'
      codeSended.value = false
      clearInterval(intervalId)
    }
  }, 1000)
}

async function handleLogin() {
  try {
    loading.value = true
    const data = await login(phone.value.trim(), code.value.trim())
    // authStore.setToken(data.token)
    // console.log('datasss: ', data)
    localStorage.access_token = data.access_token
    ms.success('登录成功')
    window.location.reload()
  }
  catch (error: any) {
    ms.error(error.message ?? 'error')
    // authStore.removeToken()
    localStorage.access_token = ''
    token.value = ''
  }
  finally {
    loading.value = false
  }
}

function handleCodePress(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey)
    event.preventDefault()
}

function handlePhonePress(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey)
    event.preventDefault()
}
</script>

<template>
  <NModal :show="visible" style="width: 90%; max-width: 640px">
    <div class="p-10 bg-white rounded dark:bg-slate-800">
      <div class="space-y-4">
        <header class="space-y-2">
          <h2 class="text-2xl font-bold text-center text-slate-800 dark:text-neutral-200">
            请先登录（试运营）
          </h2>
          <p class="text-base text-center text-slate-500 dark:text-slate-500">
            {{ $t('common.unauthorizedTips') }}
          </p>
          <Icon403 class="w-[200px] m-auto" />
        </header>
        <NInput v-model:value="phone" placeholder="请输入手机号码" @keypress="handlePhonePress" />
        <div class="code-warp">
          <NInput v-model:value="code" class="code" type="text" placeholder="请输入验证码" @keypress="handleCodePress" />
          <div class="send-code-btn">
            <NButton
              block
              type="primary"
              :disabled="sendCodeDisabled"
              :loading="codeLoading"
              @click="sendCode"
            >
              {{ codeBtnMsg }}
            </NButton>
          </div>
        </div>
        <NButton
          block
          type="primary"
          :disabled="loginDisabled"
          :loading="loading"
          @click="handleLogin"
        >
          {{ $t('common.verify') }}
        </NButton>
      </div>
    </div>
  </NModal>
</template>

<style scoped>
.code-warp {
  display: flex;
}

.code-warp .code {
	flex: 1;
}

.code-warp .send-code-btn {
	width: 130px;
	margin-left: 15px;
	flex: 0 0 130px;
}
</style>
