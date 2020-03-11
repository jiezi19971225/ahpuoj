 <template lang="pug">
  .content
    .content__main
      .siderbar
        .tags__wrapper
          ul.siderbar__item__list
            li
              .section__title 查找系列赛：
              .siderbar__searchbar__wrapper
                el-input(size="small",style="max-width:20em", placeholder="请输入系列赛名称", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                  el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
      .main.has__pagination
        h1.content__panel__title 系列赛列表
        el-table(size="small",:data="tableData")
          el-table-column(label="名称", min-width="180")
            template(slot-scope="scope")
              router-link(:to="{name:'series',params:{id:scope.row.id}}") {{scope.row.name}}
          el-table-column(label="模式", min-width="150")
            template(slot-scope="scope")
              el-tag(size="small",:type="scope.row.team_mode == 0 ? 'success':'primary'",effect="dark") {{ scope.row.team_mode == 0?"个人系列赛":"团体系列赛" }}
        el-pagination.user__pagination(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:pager-count="5",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total")
</template>

<script>
import { getSeriesList } from 'user/api/nologin';
import { mapState } from 'vuex';

export default {
  data() {
    return {
      loading: false,
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
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      this.loading = true;
      try {
        const res = await getSeriesList(
          this.currentPage,
          this.perpage,
          this.queryParam,
        );
        const { data } = res;
        this.tableData = data.data;
        this.total = data.total;
        this.loading = false;
      } catch (err) {
        console.log(err);
      }
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.fetchData();
    },
  },
};
</script>
