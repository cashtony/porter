<template>
  <div>
    <el-button @click="dailyUpdate"> 手动执行每日更新任务</el-button>
    <el-button @click="newlyUpdate"> 检测当天最新视频并更新</el-button>
  </div>
</template>

<script>
import { manuallyDailyUpdate, manuallyNewlyUpdate } from '@/api/manage'

export default {
  methods: {
    dailyUpdate() {
      manuallyDailyUpdate().then(response => {
        if (response.code !== 1000) {
          this.$message.error('操作失败:', response.code)
          return
        }

        this.$message.success('操作成功, 请到统计页面查看详情')
      })
    },
    newlyUpdate() {
      manuallyNewlyUpdate().then(response => {
        if (response.code !== 1000) {
          this.$message.error('操作失败:', response.code)
          return
        }

        this.$message.success('操作成功, 请到统计页面查看详情')
      })
    }
  }
}
</script>
