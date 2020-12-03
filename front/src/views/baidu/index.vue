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
      <el-table-column label="UID" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.uid }}</span>
        </template>
      </el-table-column>

      <el-table-column label="账号" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.userName }}</span>
        </template>
      </el-table-column>

      <el-table-column label="昵称" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.nickName }}</span>
        </template>
      </el-table-column>

      <el-table-column label="粉丝数" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.fansNum }}</span>
        </template>
      </el-table-column>

      <el-table-column label="钻石数" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.diamond }}</span>
        </template>
      </el-table-column>
      <el-table-column label="视频数量" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.videoCount }}</span>
        </template>
      </el-table-column>
      <el-table-column label="抖音绑定" width="150" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.douyinUID }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.createTime }}</span>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      layout="prev, pager, next"
      :total="totalNum"
      :page-size="20"
      @current-change="handleCurrentChange"
    >
    </el-pagination>
  </div>
</template>

<script>
import { getBaiduUserList } from "@/api/table";

export default {
  name: "BaiduAccountList",
  data() {
    return {
      list: null,
      totalNum: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
      },
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      getBaiduUserList(this.listQuery).then((response) => {
        this.list = response.users;
        this.totalNum = response.totalNum;
        this.listLoading = false;
      });
    },
    handleCurrentChange(num) {
      this.listQuery.page = num;
      this.fetchData();
    },
  },
};
</script>

<style lang="scss" scoped>
</style>
