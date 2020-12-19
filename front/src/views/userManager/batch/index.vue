<template>
  <div>
    <el-form
      ref="form"
      :model="form"
      label-width="120px"
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
import { bindAdd } from '@/api/baidu'

export default {
  name: 'AddAccount',
  data() {
    return {
      content: ''
    }
  },
  methods: {
    onSubmit() {
      bindAdd({ content: this.content }).then(response => {
        if (response.code === 1000) {
          this.$message.success('成功')
          this.content = ''
        } else {
          this.$message.error('操作失败,请联系管理员')
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
