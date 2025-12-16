<template>
    <transition name="toast-fade">
        <div v-if="visible"
            class="fixed bottom-8 left-1/2 -translate-x-1/2 px-4 py-2 rounded-xl glass text-gray-800 z-50">
            {{ message }}
        </div>
    </transition>
</template>

<script setup>
import { ref } from 'vue'

const visible = ref(false)
const message = ref('')
let timer = null

// 暴露 toast(message) 方法
const showToast = (msg) => {
    message.value = msg
    visible.value = true

    clearTimeout(timer)
    timer = setTimeout(() => {
        visible.value = false
    }, 2000)
}

defineExpose({ showToast })
</script>

<style scoped>
.toast-fade-enter-active,
.toast-fade-leave-active {
    transition: opacity 0.3s ease;
}

.toast-fade-enter-from,
.toast-fade-leave-to {
    opacity: 0;
}
</style>
