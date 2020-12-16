<template>
  <div>
    <el-form
      ref="form"
      :model="form"
      label-width="120px"
    >
      <el-form-item label="账号信息">
        <el-input
          type="textarea"
          :autosize="{ minRows: 5 }"
          placeholder="需要按照从excel表中默认复制出来的格式: 百度bduss  抖音分享链接"
          v-model="content"
        >
        </el-input>
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
import { syncUser } from '@/api/table'

export default {
  name: 'AddAccount',
  data() {
    return {
      content: ''
    }
  },
  methods: {
    onSubmit() {
      syncUser({ content: this.content }).then(response => {
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
