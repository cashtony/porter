<template>
  <div>
    <el-form
      v-loading="loading"
      label-width="120px"
      element-loading-text="拼命加载中"
    >
      <el-form-item label="账号信息">
        <el-input
          v-model="content"
          type="textarea"
          :autosize="{ minRows: 5 }"
          placeholder="百度bduss 抖音分享链接"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          @click="onSubmit"
        >提交</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { bindAdd } from '@/api/table'

export default {
  name: 'AddAccount',
  data() {
    return {
      content: '',
      loading: false
    }
  },
  methods: {
    onSubmit() {
      this.loading = true
      bindAdd({ content: this.content })
        .then(response => {
          this.$message.success('成功')
          this.content = ''

          this.loading = false
        })
        .catch(message => {
          this.loading = false
        })
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
