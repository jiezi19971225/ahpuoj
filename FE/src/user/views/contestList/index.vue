<template lang="pug">
  .content
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            .tags__wrapper
              .section__title 查找竞赛：
              .siderbar__searchbar__wrapper
                el-input(size="small",style="max-width:20em", placeholder="请输入竞赛名称", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                  el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
      .main.has__pagination
        h1.content__panel__title 竞赛列表
        el-table(size="small",:data="tableData",v-loading="loading")
          el-table-column(width="90")
            template(slot-scope="scope")
              el-tag(size="small",v-if="scope.row.status==1",type="success",effect="dark") 未开始
              el-tag(size="small",v-if="scope.row.status==2",type="primary",effect="dark") 进行中
              el-tag(size="small",v-if="scope.row.status==3",type="danger",effect="dark") 已结束
          el-table-column(label="名称", min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'contest',params:{id:scope.row.id}}") {{scope.row.name}}
          el-table-column(label="模式", min-width="150")
            template(slot-scope="scope")
              el-tag(size="small",:type="scope.row.private == 0 ? 'success':'danger'",effect="dark") {{ scope.row.private == 0?"公开赛":"私有赛" }}
              el-tag(size="small",:type="scope.row.team_mode == 0 ? 'success':'primary'",effect="dark",style="margin-left:3px;") {{ scope.row.team_mode == 0?"个人赛":"团体赛" }}
          el-table-column(label="开始时间", min-width="100",prop="start_time")
          el-table-column(label="结束时间", min-width="100",prop="end_time")
        el-pagination.user__pagination(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:pager-count="5",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total")
</template>

<script>
import { getContestList } from 'user/api/nologin';
import { mapState } from 'vuex';

export default {
  data() {
    return {
      loading: false,
      currentPage: 1,
      perpage: 20,
      queryParam: '',
      tableData: [],
      total: 0,
      tags: [],
    };
  },
  computed: {
    ...mapState({
      device: (state) => state.app.device,
    }),
  },
  activated() {
    this.fetchData();
  },
  methods: {
    async fetchData() {
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      this.loading = true;
      try {
        const res = await getContestList(
          this.currentPage,
          this.perpage,
          this.queryParam,
        );
        const { data } = res;
        this.tableData = data.data.filter((x) => x.defunct === 0);
        this.total = this.tableData.length;
        this.loading = false;
      } catch (err) {
        console.log(err);
      }
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.loading = true;
      this.fetchData();
    },
    spliteDate(dateTimeString) {
      return dateTimeString.split(' ')[0];
    },
    spliteTime(dateTimeString) {
      return dateTimeString.split(' ')[1];
    },
  },
};
</script>
<style lang="scss" scoped>
.contestlist__time__tag {
  display: inline-block;
}
</style>
