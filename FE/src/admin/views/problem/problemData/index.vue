<template lang="pug">
.admin-content
  .content-header 问题名称:{{problem?problem.title:''}}
  el-alert(title="上传文件大小限制为 80 M，上传数据文件前请先命名数据文件为 xxx.in 或 xxx.out 不支持其他扩展名文件的上传，windows不显示扩展名请先设置文件夹选项" ,type="info",style="margin-bottom:.1rem;")
  table-tools
    template(#tool)
      el-button(icon="el-icon-plus",@click="handleAdd") 新建
      el-upload(ref="upload", :action="`/api/admin/problem/${$route.params.id}/datafile`",
            :headers="{Authorization:$store.getters.token}", :filelist = "filelist",
            accept=".in,.out", :auto-upload="true", :on-success="handleUploadSuccess", :on-error="handleUploadError", :limit="1")
              el-button(slot="trigger", icon="el-icon-upload",type="primary") 上传
  el-table(:data="tableData", style="width: 100%;margin-bottom:30px;",border)
    el-table-column(label="文件名",prop="filename", width="180")
    el-table-column(label="文件大小(字节)",prop="size",width="180")
    el-table-column(label="修改时间",prop="mod_time",min-width="180")
    el-table-column(label="操作", width="180")
      template(slot-scope="scope")
        el-button(size="mini",@click="handleEditData(scope.row)") 编辑
        el-button(size="mini" type="danger" @click="deleteProblemData(scope.row)") 删除
  el-dialog(:title="dialogFormTitle",:visible.sync="dialogFormVisible",@closed="closeDialog",@opened="openDialog",width="600px",:close-on-click-modal="false")
    el-form(v-show="dialogType == 'addData'",:model="addDataForm",ref="addDataForm",:rules="addDataRules",@submit.native.prevent)
      el-form-item(abel="数据名称", prop="filename")
        el-input(v-model="addDataForm.filename",autocomplete="off",@keyup.enter.native="submit")
    el-form(v-show="dialogType == 'editData'",:model="editDataForm",ref="editDataForm",:rules="editDataRules",@submit.native.prevent)
      el-form-item(label="数据内容", prop="content")
        el-input(v-model="editDataForm.content", type="textarea", :rows="20", autocomplete="off",resize="none")
    .dialog-footer(slot="footer")
      el-button(@click="cancel") 取消
      el-button(type="primary",native-type="submit",@click="submit") 确定
</template>

<script>
import {
  getProblem,
  getProblemDataList,
  editProblemData,
  addProblemData,
  deleteProblemData,
  getProblemDataFile,
} from 'admin/api/problem';

export default {
  data() {
    return {
      problem: null,
      dialogFormTitle: '',
      dialogFormVisible: false,
      dialogType: '',
      currentRow: null,
      filelist: [],
      addDataForm: {
        filename: '',
      },
      editDataForm: {
        content: '',
      },
      addDataRules: {
        filename: [
          {
            required: true,
            message: '请输入数据名称',
            trigger: 'blur',
          },
          {
            max: 20,
            message: '超出长度限制',
            trigger: 'blur',
          },
        ],
      },
      editDataRules: {
        content: [
          {
            required: true,
            message: '请输入数据内容',
            trigger: 'blur',
          },
        ],
      },
      tableData: [],
    };
  },
  async activated() {
    const { id } = this.$route.params;
    try {
      const res = await getProblem(id);
      this.problem = res.data.problem;
      this.fetchDataList();
    } catch (err) {
      this.$store.dispatch('tagsView/delViewByRoute', this.$route);
      this.$router.replace({ name: 'admin404Page' });
      console.log(err);
    }
  },
  methods: {
    async fetchDataList() {
      const self = this;
      try {
        const res = await getProblemDataList(self.problem.id);
        console.log(res);
        self.tableData = res.data.files;
      } catch (err) {
        console.log(err);
      }
    },
    openDialog() {
      if (this.dialogType === 'addData') {
        this.$notify({
          title: '提示',
          message:
            '添加方法：输入数据名称点击确定，系统将生成对应的.in和.out文件，之后再编辑对应的数据文件。',
          duration: 6000,
        });
      }
      this.$refs.addDataForm.clearValidate();
      this.$refs.editDataForm.clearValidate();
    },
    closeDialog() {
      this.$refs.addDataForm.resetFields();
      this.$refs.editDataForm.resetFields();
      this.addDataForm.filename = '';
      this.editDataForm.content = '';
    },
    submit() {
      let refform = '';
      if (this.dialogType === 'addData') {
        refform = 'addDataForm';
      } else {
        refform = 'editDataForm';
      }
      this.$refs[refform].validate(async (valid) => {
        if (valid) {
          try {
            let res;
            if (this.dialogType === 'addData') {
              res = await addProblemData(this.problem.id, this.addDataForm);
            } else {
              res = await editProblemData(
                this.problem.id,
                this.currentRow.filename,
                this.editDataForm,
              );
            }
            this.$message({
              message: res.data.message,
              type: 'success',
            });
            res.data.info.forEach((x, index) => {
              setTimeout(() => {
                this.$notify({
                  title: '提示',
                  message: x,
                  duration: 6000,
                });
              }, 500 * index);
            });
            this.fetchDataList();
          } catch (err) {
            this.$message({
              message: err.response.data.message,
              type: 'error',
            });
          }
          this.dialogFormVisible = false;
        } else {
          return false;
        }
      });
    },
    cancel() {
      this.dialogFormVisible = false;
    },
    handleAdd() {
      this.dialogFormTitle = '添加数据';
      this.dialogType = 'addData';
      this.dialogFormVisible = true;
    },
    handleUploadSuccess(res, file, filelist) {
      this.filelist = filelist.shift();
      this.$message({
        message: res.message,
        type: 'success',
      });
      this.fetchDataList();
    },
    handleUploadError(err, file, filelist) {
      this.filelist = filelist.shift();
      this.$message({
        message: err.response.message ? err.response.message : err,
        type: 'error',
      });
    },
    async handleEditData(row) {
      this.dialogFormTitle = `编辑数据文件${row.filename}`;
      this.dialogType = 'editData';
      this.currentRow = row;
      const res = await getProblemDataFile(this.problem.id, row.filename);
      this.editDataForm.content = res.data.content;
      this.dialogFormVisible = true;
    },
    async deleteProblemData(row) {
      try {
        await this.$confirm(`确认要删除数据文件${row.filename}吗?`, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        });
        try {
          const res = await deleteProblemData(this.problem.id, row.filename);
          this.$message({
            type: 'success',
            message: res.data.message,
          });
          this.fetchDataList();
        } catch (err) {
          this.$message({
            type: 'error',
            message: err.response.data.message,
          });
        }
      } catch (err) {
        this.$message({
          type: 'info',
          message: '已取消删除',
        });
      }
    },
  },
};
</script>
