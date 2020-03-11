<template lang="pug">
.admin-content
  table-tools
    template(#tool)
      router-link(:to="{name:'adminAddProblem'}")
        el-button(icon="el-icon-plus") 新建
    template(#search)
      el-input(placeholder="搜索问题名称", style="max-width:20em", v-model="queryParam", @keyup.enter.native="handleSearchByParam",maxlength="20", clearable)
        el-button(slot="append",icon="el-icon-search",@click="handleSearchByParam")
  .content__table__wrapper
    el-table(:data="tableData" style="width: 100%", v-loading="loading",border)
      el-table-column(label="ID", prop="id", width="180")
      el-table-column(label="标题", prop="title", min-width="300")
        template(slot-scope="scope")
          a(:href="`/problem/${scope.row.id}`" target="_blank") {{scope.row.title}}
      el-table-column(label="标签", min-width="300")
        template(slot-scope="scope")
          el-tag(v-for="tag in scope.row.tags", :key="tag.id",type="success",effect="dark",style="margin-left:.04rem;") {{tag.name}}
      el-table-column(label="状态", width="180")
        template(slot-scope="scope")
          el-tag(:type="scope.row.defunct == 0 ? 'success':'danger'",effect="dark") {{scope.row.defunct == 0?'启用':'保留'}}
      el-table-column(label="操作", width="300")
        template(slot-scope="scope")
          el-button(size="mini", @click="$router.push({name:'adminProblemData',params:{id:scope.row.id}})") 数据
          el-button(size="mini", type="primary", @click="$router.push({name:'adminEditProblem',params:{id:scope.row.id}})") 编辑
          el-button(size="mini", :type="scope.row.defunct == 0?'danger':'success'", @click="handleToggleProblemStatus(scope.row)") {{scope.row.defunct == 0?'保留':'启用'}}
          el-button(size="mini", type="danger", @click="handleDeleteProblem(scope.row)") 删除
  .content__pagination__wrapper
    el-pagination(@size-change="handleSizeChange",@current-change="fetchDataList",:current-page.sync="currentPage",:page-sizes="[10, 20, 30, 40,50]",:page-size="10",layout="total, sizes, prev, pager, next, jumper",:total="total")
</template>

<script>
import {
  getProblemList,
  deleteProblem,
  toggleProblemStatus,
} from 'admin/api/problem';

export default {
  name: 'adminProblemList',
  data() {
    return {
      loading: true,
      currentPage: 1,
      currentRowId: 0,
      perpage: 10,
      total: 0,
      queryParam: '',
      tableData: [],
    };
  },
  activated() {
    this.fetchDataList();
  },
  methods: {
    async fetchDataList() {
      const self = this;
      self.loading = true;
      try {
        const res = await getProblemList(
          self.currentPage,
          self.perpage,
          self.queryParam,
        );
        const { data } = res;
        setTimeout(() => {
          self.tableData = data.data;
          self.total = data.total;
          self.loading = false;
        }, 200);
      } catch (err) {
        console.log(err);
      }
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchDataList();
    },
    handleSizeChange(val) {
      this.perpage = val;
      this.fetchDataList();
    },
    async handleToggleProblemStatus(row) {
      const self = this;
      const msg = `确认要${row.defunct === 0 ? '保留' : '启用'}问题${
        row.title
      }吗?`;
      try {
        await self.$confirm(msg, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        });
        try {
          const res = await toggleProblemStatus(row.id);
          self.$message({
            type: 'success',
            message: res.data.message,
          });

          row.defunct = 1 - row.defunct;
        } catch (err) {
          self.$message({
            type: 'error',
            message: err.response.data.message,
          });
        }
      } catch (err) {
        self.$message({
          type: 'info',
          message: '已取消操作',
        });
      }
    },
    async handleDeleteProblem(row) {
      const self = this;
      try {
        await self.$confirm(`确认要删除问题${row.title}吗?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        });
        try {
          const res = await deleteProblem(row.id);
          self.$message({
            type: 'success',
            message: res.data.message,
          });
          // 删除最后一页最后一条记录，如果不是第一页，则当前页码-1
          if (self.tableData.length === 1) {
            if (self.currentPage > 1) {
              self.currentPage -= 1;
            }
          }
          self.fetchDataList();
        } catch (err) {
          self.$message({
            type: 'error',
            message: err.response.data.message,
          });
        }
      } catch (err) {
        self.$message({
          type: 'info',
          message: '已取消删除',
        });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
</style>
