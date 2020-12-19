<template>
  <div>
    <el-form
      label-width="120px"
      v-loading="loading"
      element-loading-text="拼命加载中"
    >
      <el-form-item label="账号信息">
        <el-input
          v-model="content"
          type="textarea"
          :autosize="{ minRows: 5 }"
          placeholder="需要按照从excel表中默认复制出来的格式: 百度bduss  抖音分享链接"
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
import { syncUser } from '@/api/baidu'

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
      syncUser({ content: this.content })
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
