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
        @click="onAddDouyinuUser"
      >
        增加用户
      </el-button>

      <el-checkbox
        v-model="listQuery.hideSimilar"
        @change="fetchData"
      >只显示未搬运</el-checkbox>
      <el-checkbox
        v-model="listQuery.hideBinded"
        @change="fetchData"
      >只显示未绑定</el-checkbox>

    </div>
    <el-table
      ref="multipleTable"
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
      @selection-change="handleSelectionChange"
    >
      <el-table-column
        type="selection"
        width="55"
      />

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
        label="昵称"
        width="150"
        align="center"
        show-overflow-tooltip="true"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.nickName }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="性别"
        width="50"
        align="center"
      >
        <template slot-scope="scope">
          <span v-if="scope.row.gender == 0">未知</span>
          <span v-if="scope.row.gender == 1">男</span>
          <span v-if="scope.row.gender == 2">女</span>
        </template>
      </el-table-column>
      <el-table-column
        label="粉丝数"
        width="100"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ (scope.row.follower_count / 10000) | rounding }}万</span>
        </template>
      </el-table-column>

      <el-table-column
        label="视频数量"
        width="80"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.aweme_count }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="地区"
        width="150"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.province }}{{scope.row.location}}</span>
        </template>
      </el-table-column>

      <el-table-column
        label="签名"
        width="300"
        align="center"
        show-overflow-tooltip="true"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.signature }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="加入后台时间"
        align="center"
      >
        <template slot-scope="scope">
          <span>{{ scope.row.createTime }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="搬运状态"
        width="100"
        align="center"
      >
        <template slot-scope="scope">
          <span v-if="scope.row.similarSign === 0">未被搬运</span>
          <span v-if="scope.row.similarSign === 1">全民小视频</span>
        </template>
      </el-table-column>
      <el-table-column
        width="80"
        label="操作"
        align="center"
      >
        <template slot-scope="scope">
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
    <div style="margin-top: 20px">
      <el-button @click="submitSelection()">批量绑定</el-button>
      <el-button @click="cancelSelection()">取消选择</el-button>
    </div>
    <el-pagination
      layout="prev, pager, next"
      :total="totalNum"
      :page-size="20"
      @current-change="handleCurrentChange"
    />
    <el-dialog
      title="增加抖音账号"
      :visible.sync="AddDouyinUserDialogFormVisible"
    >

      <el-form>
        <el-form-item label="账号信息">
          <el-input
            v-model="addUserContent"
            type="textarea"
            :autosize="{ minRows: 5, maxRows:20 }"
            placeholder="一行一个抖音分享链接,可批量输入"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            @click="onSubmitAddDouyinUser"
          >提交</el-button>
        </el-form-item>
      </el-form>

    </el-dialog>
  </div>
</template>

<script>
import {
  getDouyinUserList,
  addDouyinUser,
  bind,
  deleteDouyinUser
} from '@/api/douyin'
const searchTypeOptions = [
  { key: 'nickname', display_name: '昵称' },
  { key: 'uid', display_name: '抖音UID' }
]
export default {
  name: 'DouyinAccountList',
  filters: {
    rounding(value) {
      return value.toFixed(1)
    }
  },
  data() {
    return {
      list: null,
      totalNum: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        douyinUID: '',
        nickname: '',
        hideSimilar: false,
        hideBinded: false
      },
      addUserContent: '',
      AddDouyinUserDialogFormVisible: false,
      searchTypeOptions,
      searchType: searchTypeOptions[0].key,
      searchData: '',
      multipleSelection: []
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      getDouyinUserList(this.listQuery)
        .then(response => {
          this.list = response.users
          this.totalNum = response.totalNum
          this.listLoading = false
        })
        .catch(message => {
          this.listLoading = false
        })
    },
    handleCurrentChange(num) {
      this.listQuery.page = num
      this.fetchData()
    },
    onAddDouyinuUser() {
      this.AddDouyinUserDialogFormVisible = true
    },
    onSubmitAddDouyinUser() {
      addDouyinUser({ content: this.addUserContent })
        .then(response => {
          this.$notify({
            title: '增加成功',
            type: 'success',
            duration: 2000
          })

          this.AddDouyinUserDialogFormVisible = false
          this.fetchData()
        })
        .catch(message => {
          this.AddDouyinUserDialogFormVisible = false
        })
    },
    resetQuery() {
      this.listQuery = {
        page: 1,
        limit: 20,
        douyinUID: '',
        nickname: ''
      }
    },
    onHandleFilter() {
      this.resetQuery()

      switch (this.searchType) {
        case 'nickname':
          this.listQuery.nickname = this.searchData
          break
        case 'douyinUID':
          this.listQuery.douyinUID = this.searchData
          break
      }
      this.fetchData()
    },
    onCleanFilter() {
      this.resetQuery()
      this.fetchData()
    },

    onDelete(row) {
      deleteDouyinUser({ uid: row.uid })
        .then(resp => {
          this.$notify({
            title: '删除成功',
            type: 'success',
            duration: 2000
          })

          this.fetchData()
        })
        .catch(message => {})
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    submitSelection() {
      const content = []
      this.multipleSelection.forEach(data => {
        content.push(data.uid)
      })

      if (content.length === 0) {
        return
      }
      this.listLoading = true
      bind({ content: content })
        .then(resp => {
          if (resp.resultMsg !== '') {
            this.$alert(resp.resultMsg, '操作结果', {
              confirmButtonText: '确定'
            })
          }

          this.fetchData()
          this.listLoading = false
        })
        .catch(message => {
          this.listLoading = false
        })
    },
    cancelSelection() {
      this.$refs.multipleTable.clearSelection()
    }
  }
}
</script>

<style lang="scss" scoped>
.el-checkbox {
  padding-left: 20px;
}
</style>
