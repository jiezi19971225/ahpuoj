 <template lang="pug">
  .content
    .content__main
      .siderbar
        .tags__wrapper
          ul.siderbar__item__list
            li
              .section__title 查找系列赛：
              .siderbar__searchbar__wrapper
                el-input(style="max-width:20em", placeholder="请输入系列赛名称", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                  el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
      .main
        h1.content__panel__title 系列赛列表
        el-table(:data="tableData", style="width: 100%", class="dataTable")
          el-table-column(label="名称", min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'series',params:{id:scope.row.id}}") {{scope.row.name}}
          el-table-column(label="模式", min-width="150")
            template(slot-scope="scope")
              el-tag(:type="scope.row.team_mode == 0 ? 'success':'primary'",effect="dark") {{ scope.row.team_mode == 0?"个人系列赛":"团体系列赛" }}
        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total",:small="device === 'mobile'")
</template>

<script>
import { getSeriesList } from 'user/api/nologin';
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
      const self = this;
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      try {
        const res = await getSeriesList(
          self.currentPage,
          self.perpage,
          self.queryParam,
        );
        console.log(res);
        const { data } = res;
        self.tableData = data.data;
        self.total = data.total;
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
</style>
