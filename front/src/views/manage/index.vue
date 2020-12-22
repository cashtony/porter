<template>
  <div>
    <el-button @click="onDailyUpdate"> 手动执行每日更新任务</el-button>
    <el-button @click="onNewlyUpdate"> 检测当天最新视频并更新</el-button>
    <el-button @click="onUpdateBaiduUser"> 更新全民小视频数据(钻石,T豆之类的)</el-button>
  </div>
</template>

<script>
import { manuallyDailyUpdate, manuallyNewlyUpdate } from '@/api/manage'
import { updateBaiduUser } from '@/api/baidu'

export default {
  methods: {
    onDailyUpdate() {
      manuallyDailyUpdate().then(response => {
        if (response.code !== 1000) {
          this.$message.error('操作失败:', response.code)
          return
        }

        this.$message.success('操作成功, 请到统计页面查看详情')
      })
    },
    onNewlyUpdate() {
      manuallyNewlyUpdate().then(response => {
        if (response.code !== 1000) {
          this.$message.error('操作失败:', response.code)
          return
        }

        this.$message.success('操作成功, 请到统计页面查看详情')
      })
    },
    onUpdateBaiduUser() {
      updateBaiduUser({}).then(response => {
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
