#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <dirent.h>
#include <unistd.h>
#include <time.h>
#include <stdarg.h>
#include <ctype.h>
#include <sys/wait.h>
#include <sys/ptrace.h>
#include <sys/types.h>
#include <sys/user.h>
#include <sys/syscall.h>
#include <sys/time.h>
#include <sys/resource.h>
#include <sys/signal.h>
#include <sys/stat.h>
#include <unistd.h>
#include <assert.h>
#include "okcalls.h"
#include <errno.h>
#define BUFFER_SIZE 5120

struct task
{
	int problem_id;
	int language;
	int time_limit;
	int memory_limit;
	char *source;
	char *input_text;
};

struct result
{
	int result;
	int time_used;
	int memory_used;
	double pass_rate;
	char *ceinfo;
	char *reinfo;
	char *custom_out;
	char *error_point;
};

int DEBUG = 1;

#define IGNORE_ESOL //ignore the ending space char of lines while comparing
#define STD_MB 1048576LL
#define STD_T_LIM 2
#define STD_F_LIM (STD_MB << 5) //default file size limit 32m ,2^5=32
#define STD_M_LIM (STD_MB << 7) //default memory limit 128m ,2^7=128

#define OJ_WT0 0
#define OJ_WT1 1
#define OJ_CI 2
#define OJ_RI 3
#define OJ_AC 4
#define OJ_PE 5
#define OJ_WA 6
#define OJ_TL 7
#define OJ_ML 8
#define OJ_OL 9
#define OJ_RE 10
#define OJ_CE 11
#define OJ_CO 12
#define OJ_TR 13
/*copy from ZOJ
 http://code.google.com/p/zoj/source/browse/trunk/judge_client/client/tracer.cc?spec=svn367&r=367#39
 */
#ifdef __i386
#define REG_SYSCALL orig_eax
#define REG_RET eax
#define REG_ARG0 ebx
#define REG_ARG1 ecx
#else
#define REG_SYSCALL orig_rax
#define REG_RET rax
#define REG_ARG0 rdi
#define REG_ARG1 rsi

#endif

static char oj_home[BUFFER_SIZE];
static char data_list[BUFFER_SIZE][BUFFER_SIZE];
static int data_list_len = 0;
static int java_time_bonus = 5;
static int java_memory_bonus = 512;
static char java_xms[BUFFER_SIZE];
static char java_xmx[BUFFER_SIZE];
static int oi_mode = 0;
static int use_max_time = 0;
static double cpu_compensation = 1.0;
static int shm_run = 0;
static char record_call = 0;
static int use_ptrace = 1;
static int compile_chroot = 1;

static char lang_ext[18][8] = {"c", "cc", "pas", "java", "rb", "sh", "py", "php", "pl", "cs", "m", "bas", "scm", "c", "cc", "lua", "js", "go"};

int after_equal(char *c);
void trim(char *c);
int read_buf(char *buf, const char *key, char *value);
void read_double(char *buf, const char *key, double *value);
void read_int(char *buf, const char *key, int *value);
void init_conf();
const char *getFileNameFromPath(const char *path);
int execute_cmd(const char *fmt, ...);
void copy_shell_runtime(char *work_dir);
void copy_objc_runtime(char *work_dir);
void copy_bash_runtime(char *work_dir);
void copy_ruby_runtime(char *work_dir);
void copy_guile_runtime(char *work_dir);
void copy_python_runtime(char *work_dir);
void copy_php_runtime(char *work_dir);
void copy_perl_runtime(char *work_dir);
void copy_freebasic_runtime(char *work_dir);
void copy_mono_runtime(char *work_dir);
void copy_lua_runtime(char *work_dir);
void copy_js_runtime(char *work_dir);
void copy_go_runtime(char *work_dir);
int data_list_has(char *file);
int data_list_add(char *file);
long get_file_size(const char *filename);
void write_log(const char *_fmt, ...);
void init_syscalls_limits(int lang);
FILE *read_cmd_output(const char *fmt, ...);
int isInFile(const char fname[]);
void find_next_nonspace(int *c1, int *c2, FILE *f1, FILE *f2, int *ret);
void delnextline(char s[]);
int compare(const char *file1, const char *file2);
void umount(char *work_dir);
int compile(int lang, char *work_dir);
char *escape(char s[], char t[]);
void prepare_files(char *filename, int namelen, char *infile, int p_id, char *work_dir, char *outfile, char *userfile, int runner_id);
void run_solution(int lang, char *work_dir, int time_lmt, int usedtime, int mem_lmt);
int fix_python_mis_judge(char *work_dir, int *ACflg, int *topmemory, int mem_lmt);
int fix_java_mis_judge(char *work_dir, int *ACflg, int *topmemory, int mem_lmt);
int special_judge(char *oj_home, int problem_id, char *infile, char *outfile, char *userfile);
void judge_solution(int *ACflg, int usedtime, int time_lmt, int isspj, int p_id, char *infile, char *outfile, char *userfile, int *PEflg, int lang, char *work_dir, int topmemory, int mem_lmt, int num_of_test, char *error_point);
int get_page_fault_mem(struct rusage ruse, pid_t pidApp);
void print_runtimeerror(char *err);
void watch_solution(pid_t pidApp, char *infile, int *ACflg, int isspj, char *userfile, char *outfile, int lang, int *topmemory, int mem_lmt, int *usedtime, int time_lmt, int p_id, int PEflg, char *work_dir);
void clean_workdir(char *work_dir);
void mk_shm_workdir(char *work_dir);
int count_in_files(char *dirpath);
void print_call_array();
void read_from_file(char *dist, char *filename);
struct result judge(struct task *t, char *home, int runid);