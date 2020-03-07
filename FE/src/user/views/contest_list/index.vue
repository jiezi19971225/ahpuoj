<template lang="pug">
  .content
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            .tags__wrapper
              .section__title 查找竞赛：
              .siderbar__searchbar__wrapper
                el-input(style="max-width:20em", placeholder="请输入竞赛名称", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                  el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
      .main
        h1.content__panel__title 竞赛列表
        el-table(:data="tableData", style="width: 100%", class="dataTable")
          el-table-column(width="90")
            template(slot-scope="scope")
              el-tag(v-if="scope.row.status==1",type="success",effect="dark") 未开始
              el-tag(v-if="scope.row.status==2",type="primary",effect="dark") 进行中
              el-tag(v-if="scope.row.status==3",type="danger",effect="dark") 已结束
          el-table-column(label="名称", min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'contest',params:{id:scope.row.id}}") {{scope.row.name}}
          el-table-column(label="模式", min-width="150")
            template(slot-scope="scope")
              el-tag(:type="scope.row.private == 0 ? 'success':'danger'",effect="dark") {{ scope.row.private == 0?"公开赛":"私有赛" }}
              el-tag(:type="scope.row.team_mode == 0 ? 'success':'primary'",effect="dark",style="margin-left:3px;") {{ scope.row.team_mode == 0?"个人赛":"团体赛" }}
          el-table-column(label="开始时间", min-width="100")
            template(slot-scope="scope")
              span(class="contestlist__time__tag") {{spliteDate(scope.row.start_time)}}&nbsp
              span(class="contestlist__time__tag") {{spliteTime(scope.row.start_time)}}
          el-table-column(label="结束时间", min-width="100")
            template(slot-scope="scope")
              span(class="contestlist__time__tag") {{spliteDate(scope.row.end_time)}}&nbsp
              span(class="contestlist__time__tag") {{spliteTime(scope.row.end_time)}}
        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total",:small="device === 'mobile'")
</template>

<script>
import { getContestList } from 'user/api/nologin';
import { mapState } from 'vuex';

export default {
  data() {
    return {
      currentPage: 1,
      perpage: 10,
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
      console.log('??');
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      try {
        const res = await getContestList(
          this.currentPage,
          this.perpage,
          this.queryParam,
        );
        console.log(res);
        const { data } = res;
        this.tableData = data.data.filter(x => x.defunct === 0)
        this.total = this.tableData.total;
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
