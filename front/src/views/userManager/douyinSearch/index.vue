<template>
  <div>
    <el-form
      v-loading="loading"
      label-width="120px"
      element-loading-text="拼命加载中"
    >
      <el-form-item label="抖音用户搜索">
        <el-input
          v-model="content"
          type="textarea"
          :autosize="{ minRows: 5 }"
          placeholder="每行一个关键字"
        />
      </el-form-item>

      <el-form-item label="数量">
        <el-input-number
          v-model="num"
          :min="1"
          :max="500"
          size="mini"
          label="增加搜索到的前多少个号"
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
import { searchDouyin } from '@/api/douyin'

export default {
  name: 'AddAccount',
  data() {
    return {
      content: '',
      loading: false,
      num: 100
    }
  },
  methods: {
    onSubmit() {
      this.loading = true
      searchDouyin({ content: this.content, total: this.num })
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
