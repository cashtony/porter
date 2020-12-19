<template>
  <div class="content">
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column
        label="UID"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.uid }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="账号"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.userName }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="昵称"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.nickName }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="粉丝数"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.fansNum }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="钻石数"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.diamond }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="视频数量"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.videoCount }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="抖音链接"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.douyinURL }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="创建时间"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.createTime }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="状态"
        align="center"
      >
        <template slot-scope="{row}">

          <el-switch
            v-model="row.status"
            active-color="#13ce66"
            inactive-color="#ff4949"
            :active-value=1
            :inactive-value=0
            @change="onStatusChange(row)"
          />

        </template>
      </el-table-column>

      <el-table-column
        label="操作"
        align="center"
      >
        <template slot-scope="scope">
          <span>
            <el-button
              type="primary"
              size="mini"
              @click="onEdit(scope.row)"
            >修改绑定</el-button>
          </span>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      layout="prev, pager, next"
      :total="totalNum"
      :page-size="20"
      @current-change="handleCurrentChange"
    />

    <el-dialog
      :title="dialogStatusMap[dialogStatus]"
      :visible.sync="dialogFormVisible"
    >
      <el-form :model="dialogForm">

        <el-form-item label="uid">
          <el-input
            v-model="dialogForm.uid"
            :disabled="true"
          />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input
            v-model="dialogForm.nickName"
            :disabled="true"
          />
        </el-form-item>
        <el-form-item label="抖音链接">
          <el-input v-model="dialogForm.douyinURL" />
        </el-form-item>
      </el-form>
      <div
        slot="footer"
        class="dialog-footer"
      >
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="dialogStatus==='create'?createData():updateData()"
        >
          确定
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getBaiduUserList, editBaiduUser, changeStatus } from '@/api/baidu'

export default {
  name: 'BaiduAccountList',
  data() {
    return {
      list: null,
      totalNum: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20
      },
      dialogForm: {
        uid: '',
        userName: '',
        nickName: '',
        douyinURL: ''
      },
      dialogStatusMap: { update: '修改数据', create: '新增抖音用户' },
      dialogFormVisible: false,
      dialogStatus: 'update'
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      getBaiduUserList(this.listQuery).then(response => {
        this.list = response.users
        this.totalNum = response.totalNum
        this.listLoading = false
      })
    },
    handleCurrentChange(num) {
      this.listQuery.page = num
      this.fetchData()
    },
    resetDialogForm() {
      this.dialogForm = {
        uid: '',
        userName: '',
        nickName: '',
        douyinURL: ''
      }
    },
    onEdit(row) {
      this.dialogFormVisible = true
      this.dialogStatus = 'update'
      this.dialogForm = Object.assign({}, row) // copy obj
    },
    updateData() {
      const tempData = Object.assign({}, this.dialogForm)

      editBaiduUser(tempData).then(response => {
        if (response.code !== 1000) {
          this.$notify({
            title: '操作失败了,请稍后再试:' + response.code,
            type: 'failed',
            duration: 2000
          })
          return
        }

        this.dialogFormVisible = false

        this.$notify({
          title: '用户数据变更',
          message: '更新成功',
          type: 'success',
          duration: 2000
        })

        this.fetchData()
      })
    },
    onStatusChange(row) {
      console.log('row.status:', row.uid)
      changeStatus({ uid: row.uid, status: row.status }).then(response => {
        if (response.code !== 1000) {
          this.$notify({
            title: '操作失败了,请稍后再试:' + response.code,
            type: 'failed',
            duration: 2000
          })
          return
        }

        this.$notify({
          title: '数据变更',
          message: '更新成功',
          type: 'success',
          duration: 2000
        })
      })

      // row.status = 1
      // var item = this.list[index]
      // console.log('row.status:', item.status)
      // if (item === undefined) {
      //   return
      // }
      // item.status = item.status === '0' ? '1' : '0'
      // this.list[index] = item
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
