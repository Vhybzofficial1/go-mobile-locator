<template>
  <div class="space-y-6 my-6">
    <section class="p-5 rounded-xl bg-white/60 backdrop-blur-md shadow border border-white/30">
      <h2 class="text-lg font-semibold mb-3">使用说明</h2>
      <ol class="list-decimal list-inside space-y-1.5 text-gray-700 leading-relaxed text-sm">
        <li>
          准备一个 CSV 文件。如果原本是 Excel，请通过「另存为 CSV」保存。
        </li>
        <li>
          将要查询的手机号放在第一列（A 列）。第一行是标题，系统会从第二行开始处理。
        </li>
        <li>
          系统会把省份、城市、运营商结果写入 B、C、D 列；D 列之后的内容会完整保留，不会被覆盖。
        </li>
      </ol>
    </section>
    <section class="relative flex-center flex-col p-6 border-2 border-dashed border-gray-300 rounded-xl
         bg-white/40 backdrop-blur-md shadow transition-all cursor-pointer
         hover:border-primary hover:bg-white/60" @click="openFile" @dragover.prevent="onDragOver"
      @dragleave="onDragLeave" @drop.prevent="onDrop" :class="[
        isDragging ? 'border-primary bg-white/70' : '',
        loading ? 'pointer-events-none' : ''
      ]">
      <div class="text-center">
        <div class="i-material-symbols-upload-file-rounded text-4xl opacity-70 mb-3 mx-a text-primary"></div>
        <p class="text-gray-700 text-sm">点击上传文件 或拖拽文件到这里</p>
        <p class="text-gray-500 mt-1 text-xs">仅支持 CSV 格式</p>
      </div>
      <input type="file" accept=".csv" ref="fileInput" class="hidden" @change="onFileChange" />
      <div v-if="loading" class="absolute inset-0 z-10 flex-center rounded-xl
           bg-white/60 backdrop-blur-sm">
        <div class="flex flex-col items-center gap-3 text-primary">
          <div class="i-svg-spinners-3-dots-fade text-4xl"></div>
          <span class="text-sm">处理中，请稍候…</span>
        </div>
      </div>
    </section>
    <section v-if="uploadedFile" class="relative p-5 rounded-xl bg-white/60 backdrop-blur-md shadow
         border border-white/30 space-y-4" :class="loading ? 'pointer-events-none' : ''">
      <div class="flex items-center space-x-2 text-green-600 font-medium text-sm">
        <div class="i-material-symbols-check-circle-rounded text-xl"></div>
        <span>文件已上传，可下载处理后的结果</span>
      </div>
      <div class="px-5 py-2.5 rounded-lg bg-primary text-white text-sm shadow
           hover:bg-primary/90 transition flex-1 flex-center gap-2 cursor-pointer" @click="downloadCSV">
        <div class="i-material-symbols-download-rounded text-lg"></div>
        下载处理后的 CSV 文件
      </div>
      <div v-if="loading" class="absolute inset-0 z-10 flex-center rounded-xl
           bg-white/60 backdrop-blur-sm">
        <div class="flex flex-col items-center gap-3 text-primary">
          <div class="i-svg-spinners-3-dots-fade text-3xl"></div>
          <span class="text-sm">正在生成文件…</span>
        </div>
      </div>
    </section>
  </div>
  <Toast ref="toastRef" />
</template>

<script setup>
import { ref } from "vue"
import { CarrierProcessCSV, SaveCSV } from '../../wailsjs/go/main/App'
import Toast from '@/components/Toast.vue'
const toastRef = ref(null)

const fileInput = ref(null)
const uploadedFile = ref(null)
const isDragging = ref(false)
const loading = ref(false)

// 点击上传
const openFile = () => {
  fileInput.value?.click()
}

// 进入拖拽
const onDragOver = () => {
  isDragging.value = true
}

// 离开拖拽
const onDragLeave = () => {
  isDragging.value = false
}

// 释放鼠标（确定完成拖拽）
const onDrop = (e) => {
  isDragging.value = false
  handleFile(e.dataTransfer.files[0])
}

// 文件上传
const handleFile = async (file) => {
  if (!file) return
  if (!file.name.toLowerCase().endsWith('.csv')) {
    toastRef.value.showToast('只支持 CSV 文件')
    return
  }
  loading.value = true
  try {
    // 读取文件为 Base64
    const base64 = await new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result.split(',')[1])
      reader.onerror = (err) => reject(err)
      reader.readAsDataURL(file)
    })
    // 调用 Go 接口处理
    const res = await CarrierProcessCSV(base64)
    if (res.code === 0) {
      const csvBytes = Uint8Array.from(atob(res.data), c => c.charCodeAt(0))
      uploadedFile.value = new Blob([csvBytes], { type: 'text/csv' })
      toastRef.value.showToast('文件处理完成，可以下载')
    } else {
      toastRef.value.showToast(res.message)
    }
    loading.value = false
  } catch (err) {
    toastRef.value.showToast('文件处理失败')
    loading.value = false
    console.error(err)
  }
}

// 点击下载按钮
const downloadCSV = async () => {
  if (!uploadedFile.value) return
  loading.value = true
  const arrayBuffer = await uploadedFile.value.arrayBuffer()
  const bytes = Array.from(new Uint8Array(arrayBuffer))
  await SaveCSV(bytes)
  loading.value = false
  fileInput.value = null
  uploadedFile.value = null
}

const onFileChange = (e) => {
  handleFile(e.target.files[0])
}

</script>
