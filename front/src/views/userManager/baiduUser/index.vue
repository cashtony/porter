<template>
  <div class="content">
    <div class="filter-container">

      <el-select
        v-model="searchType"
        placeholder="类型"
        class="filter-item"
        style="width: 130px"
      >
        <el-option
          v-for="item in searchTypeOptions"
          :key="item.key"
          :label="item.display_name"
          :value="item.key"
        />
      </el-select>

      <el-input
        v-model="searchData"
        placeholder="内容"
        style="width: 200px;"
        class="filter-item"
        clearable
        @clear="onCleanFilter"
        @keyup.enter.native="onHandleFilter"
      />
      <el-button
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        @click="onHandleFilter"
      >
        搜索
      </el-button>

      <el-button
        type="primary"
        @click="onExportExcel"
      >
        导出excel
      </el-button>
    </div>

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
          <span>{{ scope.row.username }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="昵称"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.nickname }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="粉丝数"
        width="80"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.fansNum }}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="钻石数"
        width="80"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.diamond }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="视频数量"
        width="80"
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
        label="抖音UID"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.douyinUID }}</span>
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
        width="80"
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
        width="100"
        label="操作"
        align="center"
      >
        <template slot-scope="scope">
          <span>
            <el-button
              type="primary"
              size="mini"
              icon="el-icon-edit"
              circle
              @click="onEdit(scope.row)"
            />
          </span>
          <span>
            <el-popconfirm
              title="确定删除这个用户吗?"
              @onConfirm="onDelete(scope.row)"
            >
              <el-button
                slot="reference"
                type="danger"
                size="mini"
                icon="el-icon-delete"
                circle
              />
            </el-popconfirm>

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
            v-model="dialogForm.nickname"
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
import {
  getBaiduUserList,
  editBaiduUser,
  changeStatus,
  deleteUser
} from '@/api/baidu'
const searchTypeOptions = [
  { key: 'douyinUID', display_name: '抖音UID' },
  { key: 'douyinURL', display_name: '抖音分享链接' },
  { key: 'nickname', display_name: '昵称' }
]
export default {
  name: 'BaiduAccountList',
  data() {
    return {
      list: null,
      totalNum: 0,
      listLoading: false,
      listQuery: {
        page: 1,
        limit: 20,
        uid: '',
        douyinUID: '',
        douyinURL: '',
        nickname: ''
      },
      dialogForm: {
        uid: '',
        username: '',
        nickname: '',
        douyinURL: ''
      },
      dialogStatusMap: { update: '修改数据', create: '新增抖音用户' },
      dialogFormVisible: false,
      dialogStatus: 'update',
      searchTypeOptions,
      searchType: searchTypeOptions[0].key,
      searchData: ''
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    resetQuery() {
      this.listQuery = {
        page: 1,
        limit: 20,
        uid: '',
        douyinUID: '',
        douyinURL: '',
        nickname: ''
      }
    },
    fetchData() {
      this.listLoading = true
      getBaiduUserList(this.listQuery)
        .then(response => {
          this.list = response.users
          this.totalNum = response.totalNum
          this.listLoading = false
        })
        .catch(messge => {})
    },
    handleCurrentChange(num) {
      this.listQuery.page = num
      this.fetchData()
    },
    resetDialogForm() {
      this.dialogForm = {
        uid: '',
        username: '',
        nickname: '',
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

      editBaiduUser({ uid: tempData.uid, douyinURL: tempData.douyinURL })
        .then(response => {
          this.dialogFormVisible = false

          this.$notify({
            title: '用户数据变更',
            message: '更新成功',
            type: 'success',
            duration: 2000
          })

          this.fetchData()
        })
        .catch(message => {})
    },
    onStatusChange(row) {
      changeStatus({ uid: row.uid, status: row.status })
        .then(response => {
          this.$notify({
            title: '用户数据变化',
            message: '更新成功',
            type: 'success',
            duration: 2000
          })
        })
        .catch(message => {})
    },
    onHandleFilter() {
      this.resetQuery()

      switch (this.searchType) {
        case 'douyinUID':
          this.listQuery.douyinUID = this.searchData
          break
        case 'douyinURL':
          this.listQuery.douyinURL = this.searchData
          break
        case 'nickname':
          this.listQuery.nickname = this.searchData
          break
      }
      this.fetchData()
    },
    onExportExcel() {
      const aTag = document.createElement('a')
      aTag.href = '/baidu/user/excel'
      aTag.click()
    },
    onDelete(row) {
      deleteUser({ uid: row.uid })
        .then(response => {
          this.fetchData()
        })
        .catch(messge => {})
    },
    onCleanFilter() {
      this.resetQuery()
      this.fetchData()
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
