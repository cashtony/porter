<template>
  <div>
    <div class="filter-container">
      <el-date-picker
        v-model="listQuery.pickDate"
        type="date"
        placeholder="选择天数"
        value-format="yyyy-MM-dd"
        :picker-options="pickerOptions"
        @change="dateConfirm"
      />
    </div>
    <div>
      当日一共更新了{{totalNum}}个用户,共计{{totalVideos}}条视频
    </div>
    <div>
      <el-table
        v-loading="listLoading"
        :data="list"
        border
      >
        <el-table-column
          label="UID"
          align="center"
          width="180"
        >
          <template slot-scope="{row}">
            <span>{{ row.uid }}</span>
          </template>
        </el-table-column>

        <el-table-column
          label="昵称"
          align="center"
        >
          <template slot-scope="{row}">
            <span>{{ row.nickname }}</span>
          </template>
        </el-table-column>

        <el-table-column
          label="更新数量"
          align="center"
          width="120"
          prop="num"
        >
          <!-- <template slot-scope="{row}">
            <span>{{ row.num }}</span>
          </template> -->
        </el-table-column>
        <el-table-column
          label="日期"
          align="center"
          prop="date"
          :formatter="formatter"
        >

          <!-- <template slot-scope="{row}">
            <span>{{ row.date }}</span>
          </template> -->
        </el-table-column>
      </el-table>
      <el-pagination
        layout="prev, pager, next"
        :total="totalNum"
        :page-size="20"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script>
import { getStatistic } from '@/api/table'
export default {
  data() {
    return {
      list: null,
      totalNum: 0,
      totalVideos: 0,
      listLoading: false,

      listQuery: {
        page: 1,
        limit: 20,
        pickDate: ''
      },
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now()
        }
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      getStatistic(this.listQuery).then(response => {
        this.list = response.list
        this.totalNum = response.totalNum
        this.totalVideos = response.totalVideos
        this.listLoading = false
      })
    },
    handleCurrentChange(num) {
      this.listQuery.page = num
      this.fetchData()
    },
    dateConfirm(date) {
      this.fetchData()
    },
    formatter(row, column, cellValue, index) {
      var d = new Date(cellValue)

      return `${d.getFullYear()}年${d.getMonth()}月${d.getDate()}日`
    }
  }
}
</script>
