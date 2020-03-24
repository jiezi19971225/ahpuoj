#include "judge_client.h"
int after_equal(char *c)
{
	int i = 0;
	for (; c[i] != '\0' && c[i] != '='; i++)
		;
	return ++i;
}

void trim(char *c)
{
	char buf[BUFFER_SIZE];
	char *start, *end;
	strcpy(buf, c);
	start = buf;
	while (isspace(*start))
		start++;
	end = start;
	while (!isspace(*end))
		end++;
	*end = '\0';
	strcpy(c, start);
}

int read_buf(char *buf, const char *key, char *value)
{
	if (strncmp(buf, key, strlen(key)) == 0)
	{
		strcpy(value, buf + after_equal(buf));
		trim(value);
		if (DEBUG)
			printf("%s\n", value);
		return 1;
	}
	return 0;
}

void read_double(char *buf, const char *key, double *value)
{
	char buf2[BUFFER_SIZE];
	if (read_buf(buf, key, buf2))
		sscanf(buf2, "%lf", value);
}

void read_int(char *buf, const char *key, int *value)
{
	char buf2[BUFFER_SIZE];
	if (read_buf(buf, key, buf2))
		sscanf(buf2, "%d", value);
}

void init_conf()
{
	FILE *fp = NULL;
	char buf[BUFFER_SIZE];
	strcpy(java_xms, "-Xms32m");
	strcpy(java_xmx, "-Xmx256m");
	sprintf(buf, "%s/etc/judge.conf", oj_home);
	fp = fopen("./etc/judge.conf", "re");
	if (fp != NULL)
	{
		while (fgets(buf, BUFFER_SIZE - 1, fp))
		{
			read_int(buf, "OJ_JAVA_TIME_BONUS", &java_time_bonus);
			read_int(buf, "OJ_JAVA_MEMORY_BONUS", &java_memory_bonus);
			read_buf(buf, "OJ_JAVA_XMS", java_xms);
			read_buf(buf, "OJ_JAVA_XMX", java_xmx);
			read_int(buf, "OJ_OI_MODE", &oi_mode);
			read_int(buf, "OJ_SHM_RUN", &shm_run);
			read_int(buf, "OJ_USE_MAX_TIME", &use_max_time);
			read_int(buf, "OJ_USE_PTRACE", &use_ptrace);
			read_int(buf, "OJ_COMPILE_CHROOT", &compile_chroot);
			read_double(buf, "OJ_CPU_COMPENSATION", &cpu_compensation);
		}
		//fclose(fp);
	}
}

const char *getFileNameFromPath(const char *path)
{
	for (int i = strlen(path); i >= 0; i--)
	{
		if (path[i] == '/')
			return &path[i + 1];
	}
	return path;
}

int execute_cmd(const char *fmt, ...)
{
	printf("execute_cmd\n");
	char cmd[BUFFER_SIZE];

	int ret = 0;
	va_list ap;

	va_start(ap, fmt);
	vsprintf(cmd, fmt, ap);
	printf("%s\n", cmd);
	ret = system(cmd);
	va_end(ap);
	return ret;
}

void copy_shell_runtime(char *work_dir)
{
	execute_cmd("/bin/mkdir %s/lib", work_dir);
	execute_cmd("/bin/mkdir %s/lib64", work_dir);
	execute_cmd("/bin/mkdir %s/bin", work_dir);

#ifdef __i386
	execute_cmd("/bin/cp /lib/ld-linux* %s/lib/", work_dir);
	execute_cmd("/bin/cp -a /lib/i386-linux-gnu  %s/lib/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/i386-linux-gnu %s/lib/", work_dir);
#endif

	execute_cmd("/bin/cp -a /lib/x86_64-linux-gnu %s/lib/", work_dir);
	execute_cmd("/bin/cp /lib64/* %s/lib64/", work_dir);
	//	execute_cmd("/bin/cp /lib32 %s/", work_dir);
	execute_cmd("/bin/cp /bin/busybox %s/bin/", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/sh", work_dir);
	execute_cmd("/bin/cp /bin/bash %s/bin/bash", work_dir);
}
void copy_objc_runtime(char *work_dir)
{
	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir -p %s/proc", work_dir);
	execute_cmd("/bin/mount -o bind /proc %s/proc", work_dir);
	execute_cmd("/bin/mkdir -p %s/lib/", work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/libdbus-1.so.3                          %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/libgcc_s.so.1                           %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/libgcrypt.so.11                         %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/libgpg-error.so.0                       %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/libz.so.1                               %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/tls/i686/cmov/libc.so.6                 %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/tls/i686/cmov/libdl.so.2                %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/tls/i686/cmov/libm.so.6                 %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/tls/i686/cmov/libnsl.so.1               %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/tls/i686/cmov/libpthread.so.0           %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /lib/tls/i686/cmov/librt.so.1                %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libavahi-client.so.3                %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libavahi-common.so.3                %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libdns_sd.so.1                      %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libffi.so.5                         %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libgnustep-base.so.1.19             %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libgnutls.so.26                     %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libobjc.so.2                        %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libtasn1.so.3                       %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libxml2.so.2                        %s/lib/ ",
		work_dir);
	execute_cmd(
		"/bin/cp -aL /usr/lib/libxslt.so.1                        %s/lib/ ",
		work_dir);
}
void copy_bash_runtime(char *work_dir)
{
	//char cmd[BUFFER_SIZE];
	//const char * ruby_run="/usr/bin/ruby";
	copy_shell_runtime(work_dir);
	execute_cmd("/bin/cp `which bc`  %s/bin/", work_dir);
	execute_cmd("busybox dos2unix Main.sh", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/grep", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/awk", work_dir);
	execute_cmd("/bin/cp /bin/sed %s/bin/sed", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/cut", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/sort", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/join", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/wc", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/tr", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/dc", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/dd", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/cat", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/tail", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/head", work_dir);
	execute_cmd("/bin/ln -s /bin/busybox %s/bin/xargs", work_dir);
	execute_cmd("chmod +rx %s/Main.sh", work_dir);
}
void copy_ruby_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("mkdir -p %s/usr", work_dir);
	execute_cmd("mkdir -p %s/usr/lib", work_dir);
	execute_cmd("mkdir -p %s/usr/lib64", work_dir);
	execute_cmd("cp -a /usr/lib/libruby* %s/usr/lib/", work_dir);
	execute_cmd("cp -a /usr/lib/ruby* %s/usr/lib/", work_dir);
	execute_cmd("cp -a /usr/lib64/ruby* %s/usr/lib64/", work_dir);
	execute_cmd("cp -a /usr/lib64/libruby* %s/usr/lib64/", work_dir);
	execute_cmd("cp -a /usr/bin/ruby* %s/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/x86_64-linux-gnu/libruby* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/x86_64-linux-gnu/libgmp* %s/usr/lib/", work_dir);
}

void copy_guile_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir -p %s/proc", work_dir);
	execute_cmd("/bin/mount -o bind /proc %s/proc", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/lib", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/share", work_dir);
	execute_cmd("/bin/cp -a /usr/share/guile %s/usr/share/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libguile* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libgc* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libffi* %s/usr/lib/",
				work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libunistring* %s/usr/lib/",
				work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libgmp* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libgmp* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libltdl* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libltdl* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/bin/guile* %s/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/x86_64-linux-gnu/libguile* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/x86_64-linux-gnu/libgc* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/x86_64-linux-gnu/libffi* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/x86_64-linux-gnu/libunistring* %s/usr/lib/", work_dir);
}

void copy_python_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("mkdir -p %s/usr/include", work_dir);
	execute_cmd("mkdir -p %s/dev", work_dir);
	execute_cmd("mkdir -p %s/usr/lib", work_dir);
	execute_cmd("mkdir -p %s/usr/lib64", work_dir);
	execute_cmd("mkdir -p %s/usr/local/lib", work_dir);
	execute_cmd("cp /usr/bin/python* %s/", work_dir);
	execute_cmd("cp -a /usr/lib/python* %s/usr/lib/", work_dir);
	execute_cmd("cp -a /usr/lib64/python* %s/usr/lib64/", work_dir);
	execute_cmd("cp -a /usr/local/lib/python* %s/usr/local/lib/", work_dir);
	execute_cmd("cp -a /usr/include/python* %s/usr/include/", work_dir);
	execute_cmd("cp -a /usr/lib/libpython* %s/usr/lib/", work_dir);
	execute_cmd("/bin/mkdir -p %s/home/judge", work_dir);
	execute_cmd("/bin/chown judge %s", work_dir);
	execute_cmd("/bin/mkdir -p %s/etc", work_dir);
	execute_cmd("/bin/grep judge /etc/passwd>%s/etc/passwd", work_dir);
	execute_cmd("/bin/mount -o bind /dev %s/dev", work_dir);
}
void copy_php_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir %s/usr", work_dir);
	execute_cmd("/bin/mkdir %s/usr/lib", work_dir);
	execute_cmd("/bin/cp /usr/lib/libedit* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libdb* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libgssapi_krb5* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libkrb5* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libk5crypto* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libedit* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libdb* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libgssapi_krb5* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libkrb5* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/*/libk5crypto* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libxml2* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libxml2.so* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libicuuc.so* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libicudata.so* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libstdc++.so* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libssl* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libcrypto* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/bin/php* %s/", work_dir);
	execute_cmd("chmod +rx %s/Main.php", work_dir);
}
void copy_perl_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir %s/usr", work_dir);
	execute_cmd("/bin/mkdir %s/usr/lib", work_dir);
	execute_cmd("/bin/cp /usr/lib/libperl* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /usr/bin/perl* %s/", work_dir);
}
void copy_freebasic_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/local/lib", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/local/bin", work_dir);
	execute_cmd("/bin/cp /usr/local/lib/freebasic %s/usr/local/lib/", work_dir);
	execute_cmd("/bin/cp /usr/local/bin/fbc %s/", work_dir);
	execute_cmd("/bin/cp -a /lib32/* %s/lib/", work_dir);
}
void copy_mono_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir %s/usr", work_dir);
	execute_cmd("/bin/mkdir %s/proc", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/lib/mono/2.0", work_dir);
	execute_cmd("/bin/cp -a /usr/lib/mono %s/usr/lib/", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/lib64/mono/2.0", work_dir);
	execute_cmd("/bin/cp -a /usr/lib64/mono %s/usr/lib64/", work_dir);

	execute_cmd("/bin/cp /usr/lib/libgthread* %s/usr/lib/", work_dir);

	execute_cmd("/bin/mount -o bind /proc %s/proc", work_dir);
	execute_cmd("/bin/cp /usr/bin/mono* %s/", work_dir);

	execute_cmd("/bin/cp /usr/lib/libgthread* %s/usr/lib/", work_dir);
	execute_cmd("/bin/cp /lib/libglib* %s/lib/", work_dir);
	execute_cmd("/bin/cp /lib/tls/i686/cmov/lib* %s/lib/tls/i686/cmov/",
				work_dir);
	execute_cmd("/bin/cp /lib/libpcre* %s/lib/", work_dir);
	execute_cmd("/bin/cp /lib/ld-linux* %s/lib/", work_dir);
	execute_cmd("/bin/cp /lib64/ld-linux* %s/lib64/", work_dir);
	execute_cmd("/bin/mkdir -p %s/home/judge", work_dir);
	execute_cmd("/bin/chown judge %s/home/judge", work_dir);
	execute_cmd("/bin/mkdir -p %s/etc", work_dir);
	execute_cmd("/bin/grep judge /etc/passwd>%s/etc/passwd", work_dir);
}
void copy_lua_runtime(char *work_dir)
{

	copy_shell_runtime(work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/local/lib", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/local/bin", work_dir);
	execute_cmd("/bin/cp /usr/bin/lua %s/", work_dir);
}
void copy_js_runtime(char *work_dir)
{

	//copy_shell_runtime(work_dir);
	execute_cmd("mkdir -p %s/dev", work_dir);
	execute_cmd("/bin/mount -o bind /dev %s/dev", work_dir);
	execute_cmd("/bin/mkdir -p %s/usr/lib %s/lib/i386-linux-gnu/ %s/lib64/", work_dir, work_dir, work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libz.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libuv.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libicui18n.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libicuuc.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libicudata.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libtinfo.so.*  %s/lib/i386-linux-gnu/", work_dir);

	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libcares.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/libv8.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libssl.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libcrypto.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libdl.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/librt.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/i386-linux-gnu/libstdc++.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libpthread.so.*  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libc.so.6  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libm.so.6  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/i386-linux-gnu/libgcc_s.so.1  %s/lib/i386-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/ld-linux.so.*  %s/lib/", work_dir);

	execute_cmd("/bin/mkdir -p %s/usr/lib %s/lib/x86_64-linux-gnu/", work_dir, work_dir);

	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libz.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libuv.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/librt.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libpthread.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libdl.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libssl.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libcrypto.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libicui18n.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libicuuc.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libstdc++.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libm.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libgcc_s.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib/x86_64-linux-gnu/libc.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /lib64/ld-linux-x86-64.so.* %s/lib64/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libicudata.so.* %s/lib/x86_64-linux-gnu/", work_dir);

	// nodejs需要的共享库
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libhttp_parser.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libcares.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libnghttp2.so.* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libcrypto* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/lib/x86_64-linux-gnu/libssl* %s/lib/x86_64-linux-gnu/", work_dir);
	execute_cmd("/bin/cp /usr/bin/nodejs %s/", work_dir);
}

void copy_go_runtime(char *work_dir)
{
	char envbuff[1024] = {0};

	copy_shell_runtime(work_dir);

	execute_cmd("/bin/mkdir %s/usr", work_dir);
	execute_cmd("/bin/mkdir %s/usr/lib", work_dir);
	execute_cmd("/bin/mkdir %s/usr/bin", work_dir);
	putenv((char *)"GOROOT=/usr/lib/go");
	sprintf(envbuff, "GOPATH=%s", work_dir);
	putenv(envbuff);
}

int data_list_has(char *file)
{
	for (int i = 0; i < data_list_len; i++)
	{
		if (strcmp(data_list[i], file) == 0)
			return 1;
	}
	return 0;
}
int data_list_add(char *file)
{
	if (data_list_len < BUFFER_SIZE - 1)
	{
		strcpy(data_list[data_list_len], file);
		data_list_len++;
		return 0;
	}
	else
	{
		return 1;
	}
}
long get_file_size(const char *filename)
{
	struct stat f_stat;

	if (stat(filename, &f_stat) == -1)
	{
		return 0;
	}

	return (long)f_stat.st_size;
}

void write_log(const char *_fmt, ...)
{
	va_list ap;
	char fmt[4096];
	strncpy(fmt, _fmt, 4096);
	char buffer[4096];
	//      time_t          t = time(NULL);
	//int l;
	sprintf(buffer, "%s/log/client.log", oj_home);
	FILE *fp = fopen(buffer, "ae+");
	if (fp == NULL)
	{
		fprintf(stderr, "%s!\n", buffer);
		fprintf(stderr, "openfile error!\n");
		system("pwd");
	}
	va_start(ap, _fmt);
	//l =
	vsprintf(buffer, fmt, ap);
	fprintf(fp, "%s\n", buffer);
	if (DEBUG)
		printf("%s\n", buffer);
	va_end(ap);
	fclose(fp);
}

void write_chroot_log(const char *_fmt, ...)
{
	va_list ap;
	char fmt[4096];
	strncpy(fmt, _fmt, 4096);
	char buffer[4096];
	//      time_t          t = time(NULL);
	//int l;
	FILE *fp = fopen("/client.log", "ae+");
	if (fp == NULL)
	{
		fprintf(stderr, "%s!\n", buffer);
		fprintf(stderr, "openfile error!\n");
		system("pwd");
	}
	va_start(ap, _fmt);
	//l =
	vsprintf(buffer, fmt, ap);
	fprintf(fp, "%s\n", buffer);
	if (DEBUG)
		printf("%s\n", buffer);
	va_end(ap);
	fclose(fp);
}

#define call_array_size 512
unsigned int call_id = 0;
unsigned int call_counter[call_array_size] = {0};
static char LANG_NAME[BUFFER_SIZE];
void init_syscalls_limits(int lang)
{
	int i;
	memset(call_counter, 0, sizeof(call_counter));
	if (DEBUG)
		write_log("init_call_counter:%d", lang);
	if (record_call)
	{ // recording for debuging
		for (i = 0; i < call_array_size; i++)
		{
			call_counter[i] = 0;
		}
	}
	else if (lang <= 1 || lang == 13 || lang == 14)
	{ // C & C++
		for (i = 0; i == 0 || LANG_CV[i]; i++)
		{
			call_counter[LANG_CV[i]] = HOJ_MAX_LIMIT;
		}
	}
	else if (lang == 2)
	{ // Pascal
		for (i = 0; i == 0 || LANG_PV[i]; i++)
			call_counter[LANG_PV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 3)
	{ // Java
		for (i = 0; i == 0 || LANG_JV[i]; i++)
			call_counter[LANG_JV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 4)
	{ // Ruby
		for (i = 0; i == 0 || LANG_RV[i]; i++)
			call_counter[LANG_RV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 5)
	{ // Bash
		for (i = 0; i == 0 || LANG_BV[i]; i++)
			call_counter[LANG_BV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 6)
	{ // Python
		for (i = 0; i == 0 || LANG_YV[i]; i++)
			call_counter[LANG_YV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 7)
	{ // php
		for (i = 0; i == 0 || LANG_PHV[i]; i++)
			call_counter[LANG_PHV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 8)
	{ // perl
		for (i = 0; i == 0 || LANG_PLV[i]; i++)
			call_counter[LANG_PLV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 9)
	{ // mono c#
		for (i = 0; i == 0 || LANG_CSV[i]; i++)
			call_counter[LANG_CSV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 10)
	{ //objective c
		for (i = 0; i == 0 || LANG_OV[i]; i++)
			call_counter[LANG_OV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 11)
	{ //free basic
		for (i = 0; i == 0 || LANG_BASICV[i]; i++)
			call_counter[LANG_BASICV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 12)
	{ //scheme guile
		for (i = 0; i == 0 || LANG_SV[i]; i++)
			call_counter[LANG_SV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 15)
	{ //lua
		for (i = 0; i == 0 || LANG_LUAV[i]; i++)
			call_counter[LANG_LUAV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 16)
	{ //nodejs
		for (i = 0; i == 0 || LANG_JSV[i]; i++)
			call_counter[LANG_JSV[i]] = HOJ_MAX_LIMIT;
	}
	else if (lang == 17)
	{ //go
		for (i = 0; i == 0 || LANG_GOV[i]; i++)
			call_counter[LANG_GOV[i]] = HOJ_MAX_LIMIT;
	}
}

FILE *read_cmd_output(const char *fmt, ...)
{
	char cmd[BUFFER_SIZE];

	FILE *ret = NULL;
	va_list ap;

	va_start(ap, fmt);
	vsprintf(cmd, fmt, ap);
	va_end(ap);
	if (DEBUG)
		printf("%s\n", cmd);
	ret = popen(cmd, "r");

	return ret;
}

int isInFile(const char fname[])
{
	int l = strlen(fname);
	if (l <= 3 || strcmp(fname + l - 3, ".in") != 0)
		return 0;
	else
		return l - 3;
}
void find_next_nonspace(int &c1, int &c2, FILE *&f1, FILE *&f2, int &ret)
{
	// Find the next non-space character or \n.
	while ((isspace(c1)) || (isspace(c2)))
	{
		if (c1 != c2)
		{
			if (c2 == EOF)
			{
				do
				{
					c1 = fgetc(f1);
				} while (isspace(c1));
				continue;
			}
			else if (c1 == EOF)
			{
				do
				{
					c2 = fgetc(f2);
				} while (isspace(c2));
				continue;
#ifdef IGNORE_ESOL
			}
			else if (isspace(c1) && isspace(c2))
			{
				while (c2 == '\n' && isspace(c1) && c1 != '\n')
					c1 = fgetc(f1);
				while (c1 == '\n' && isspace(c2) && c2 != '\n')
					c2 = fgetc(f2);

#else
			}
			else if ((c1 == '\r' && c2 == '\n'))
			{
				c1 = fgetc(f1);
			}
			else if ((c2 == '\r' && c1 == '\n'))
			{
				c2 = fgetc(f2);
#endif
			}
			else
			{
				if (DEBUG)
					printf("%d=%c\t%d=%c", c1, c1, c2, c2);
				;
				ret = OJ_PE;
			}
		}
		if (isspace(c1))
		{
			c1 = fgetc(f1);
		}
		if (isspace(c2))
		{
			c2 = fgetc(f2);
		}
	}
}

void delnextline(char s[])
{
	int L;
	L = strlen(s);
	while (L > 0 && (s[L - 1] == '\n' || s[L - 1] == '\r'))
		s[--L] = 0;
}

/*
 * translated from ZOJ judger r367
 * http://code.google.com/p/zoj/source/browse/trunk/judge_client/client/text_checker.cc#25
 *
 */
int compare(const char *file1, const char *file2)
{
	int ret = OJ_AC;
	int c1, c2;
	FILE *f1, *f2;
	f1 = fopen(file1, "re");
	f2 = fopen(file2, "re");
	if (!f1 || !f2)
	{
		ret = OJ_RE;
	}
	else
		for (;;)
		{
			// Find the first non-space character at the beginning of line.
			// Blank lines are skipped.
			c1 = fgetc(f1);
			c2 = fgetc(f2);
			find_next_nonspace(c1, c2, f1, f2, ret);
			// Compare the current line.
			for (;;)
			{
				// Read until 2 files return a space or 0 together.
				while ((!isspace(c1) && c1) || (!isspace(c2) && c2))
				{
					if (c1 == EOF && c2 == EOF)
					{
						goto end;
					}
					if (c1 == EOF || c2 == EOF)
					{
						break;
					}
					if (c1 != c2)
					{
						// Consecutive non-space characters should be all exactly the same
						ret = OJ_WA;
						goto end;
					}
					c1 = fgetc(f1);
					c2 = fgetc(f2);
				}
				find_next_nonspace(c1, c2, f1, f2, ret);
				if (c1 == EOF && c2 == EOF)
				{
					goto end;
				}
				if (c1 == EOF || c2 == EOF)
				{
					ret = OJ_WA;
					goto end;
				}

				if ((c1 == '\n' || !c1) && (c2 == '\n' || !c2))
				{
					break;
				}
			}
		}
end:
	// 此处有修改，当答案错的的时候改为记录发生错误的数据点
	// if (ret == OJ_WA)
	// {
	// 	// if (full_diff)
	// 	// make_diff_out_full(f1, f2, c1, c2, file1);
	// 	// else
	// }
	if (f1)
		fclose(f1);
	if (f2)
		fclose(f2);
	return ret;
}

void umount(char *work_dir)
{
	execute_cmd("/bin/umount -f %s/proc", work_dir);
	execute_cmd("/bin/umount -f %s/dev ", work_dir);
	execute_cmd("/bin/umount -f %s/lib ", work_dir);
	execute_cmd("/bin/umount -f %s/lib64 ", work_dir);
	execute_cmd("/bin/umount -f %s/etc/alternatives ", work_dir);
	execute_cmd("/bin/umount -f %s/usr ", work_dir);
	execute_cmd("/bin/umount -f %s/bin ", work_dir);
	// 在docker下不能umount /proc目录
	execute_cmd("/bin/umount -f bin usr lib lib64 etc/alternatives dev ");
	execute_cmd("/bin/umount -f %s/* ", work_dir);
	execute_cmd("/bin/umount -f %s/log/* ", work_dir);
	execute_cmd("/bin/umount -f %s/log/etc/alternatives ", work_dir);
}
int compile(int lang, char *work_dir)
{
	int pid;

	const char *CP_C[] = {"gcc", "Main.c", "-o", "Main", "-O2", "-fmax-errors=10", "-Wall",
						  "-lm", "--static", "-std=c99", "-DONLINE_JUDGE", NULL};
	const char *CP_X[] = {"g++", "-fno-asm", "-O2", "-fmax-errors=10", "-Wall",
						  "-lm", "--static", "-std=c++17", "-DONLINE_JUDGE", "-o", "Main", "Main.cc", NULL};
	const char *CP_P[] =
		{"fpc", "Main.pas", "-Cs32000000", "-Sh", "-O2", "-Co", "-Ct", "-Ci", NULL};
	//      const char * CP_J[] = { "javac", "-J-Xms32m", "-J-Xmx256m","-encoding","UTF-8", "Main.java",NULL };

	const char *CP_R[] = {"ruby", "-c", "Main.rb", NULL};
	const char *CP_B[] = {"chmod", "+rx", "Main.sh", NULL};
	//const char * CP_Y[] = { "python", "-c",
	//		"import py_compile; py_compile.compile(r'Main.py')", NULL };
	const char *CP_PH[] = {"php", "-l", "Main.php", NULL};
	const char *CP_PL[] = {"perl", "-c", "Main.pl", NULL};
	const char *CP_CS[] = {"gmcs", "-warn:0", "Main.cs", NULL};
	const char *CP_OC[] = {"gcc", "-o", "Main", "Main.m",
						   "-fconstant-string-class=NSConstantString", "-I",
						   "/usr/include/GNUstep/", "-L", "/usr/lib/GNUstep/Libraries/",
						   "-lobjc", "-lgnustep-base", NULL};
	const char *CP_BS[] = {"fbc", "-lang", "qb", "Main.bas", NULL};
	const char *CP_CLANG[] = {"clang", "Main.c", "-o", "Main", "-ferror-limit=10", "-fno-asm", "-Wall",
							  "-lm", "--static", "-std=c99", "-DONLINE_JUDGE", NULL};
	const char *CP_CLANG_CPP[] = {"clang++", "Main.cc", "-o", "Main", "-ferror-limit=10", "-fno-asm", "-Wall",
								  "-lm", "--static", "-std=c++0x", "-DONLINE_JUDGE", NULL};
	const char *CP_LUA[] = {"luac", "-o", "Main", "Main.lua", NULL};
	//const char * CP_JS[] = { "js24","-c", "Main.js", NULL };
	const char *CP_GO[] = {"/usr/lib/go/bin/go", "build", "-o", "Main", "Main.go", NULL};

	char javac_buf[7][32];
	char *CP_J[7];

	for (int i = 0; i < 7; i++)
		CP_J[i] = javac_buf[i];

	sprintf(CP_J[0], "javac");
	sprintf(CP_J[1], "-J%s", java_xms);
	sprintf(CP_J[2], "-J%s", java_xmx);
	sprintf(CP_J[3], "-encoding");
	sprintf(CP_J[4], "UTF-8");
	sprintf(CP_J[5], "Main.java");
	CP_J[6] = (char *)NULL;

	pid = fork();
	if (pid == 0)
	{
		struct rlimit LIM;
		int cpu = 6;
		if (lang == 3)
			cpu = 30;
		LIM.rlim_max = cpu;
		LIM.rlim_cur = cpu;
		setrlimit(RLIMIT_CPU, &LIM);
		alarm(cpu);
		LIM.rlim_max = 40 * STD_MB;
		LIM.rlim_cur = 40 * STD_MB;
		setrlimit(RLIMIT_FSIZE, &LIM);

		if (lang == 3 || lang == 17)
		{
#ifdef __i386
			LIM.rlim_max = STD_MB << 11;
			LIM.rlim_cur = STD_MB << 11;
#else
			LIM.rlim_max = STD_MB << 12;
			LIM.rlim_cur = STD_MB << 12;
#endif
		}
		else
		{
			LIM.rlim_max = STD_MB * 512;
			LIM.rlim_cur = STD_MB * 512;
		}
		if (lang != 3)
			setrlimit(RLIMIT_AS, &LIM);
		if (lang != 2 && lang != 11)
		{
			freopen("ce.txt", "w", stderr);
			//freopen("/dev/null", "w", stdout);
		}
		else
		{
			freopen("ce.txt", "w", stdout);
		}

		if (compile_chroot && lang != 3 && lang != 9 && lang != 6 && lang != 11)
		{
			execute_cmd("mkdir -p bin usr lib lib64 etc/alternatives proc tmp dev");
			execute_cmd("chown judge *");
			execute_cmd("mount -o bind /bin bin");
			execute_cmd("mount -o bind /usr usr");
			execute_cmd("mount -o bind /lib lib");
#ifndef __i386
			execute_cmd("mount -o bind /lib64 lib64");
#endif
			execute_cmd("mount -o bind /etc/alternatives etc/alternatives");
			execute_cmd("mount -o bind /proc proc");
			if (lang > 2 && lang != 10 && lang != 13 && lang != 14)
				execute_cmd("mount -o bind /dev dev");
			chroot(work_dir);
		}
		while (setgid(1536) != 0)
			sleep(1);
		while (setuid(1536) != 0)
			sleep(1);
		while (setresuid(1536, 1536, 1536) != 0)
			sleep(1);

		switch (lang)
		{
		case 0:
			execvp(CP_C[0], (char *const *)CP_C);
			break;
		case 1:
			execvp(CP_X[0], (char *const *)CP_X);
			break;
		case 2:
			execvp(CP_P[0], (char *const *)CP_P);
			break;
		case 3:
			execvp(CP_J[0], (char *const *)CP_J);
			break;
		case 4:
			execvp(CP_R[0], (char *const *)CP_R);
			break;
		case 5:
			execvp(CP_B[0], (char *const *)CP_B);
			break;
		//case 6:
		//	execvp(CP_Y[0], (char * const *) CP_Y);
		//	break;
		case 7:
			execvp(CP_PH[0], (char *const *)CP_PH);
			break;
		case 8:
			execvp(CP_PL[0], (char *const *)CP_PL);
			break;
		case 9:
			execvp(CP_CS[0], (char *const *)CP_CS);
			break;

		case 10:
			execvp(CP_OC[0], (char *const *)CP_OC);
			break;
		case 11:
			execvp(CP_BS[0], (char *const *)CP_BS);
			break;
		case 13:
			execvp(CP_CLANG[0], (char *const *)CP_CLANG);
			break;
		case 14:
			execvp(CP_CLANG_CPP[0], (char *const *)CP_CLANG_CPP);
			break;
		case 15:
			execvp(CP_LUA[0], (char *const *)CP_LUA);
			break;
		//case 16:
		//	execvp(CP_JS[0], (char * const *) CP_JS);
		//	break;
		case 17:
			execvp(CP_GO[0], (char *const *)CP_GO);
			break;
		default:
			printf("nothing to do!\n");
		}
		if (DEBUG)
			printf("compile end!\n");
		// 子线程推出
		exit(0);
	}
	else
	{
		int status = 0;
		printf("status=111");
		waitpid(pid, &status, 0);
		if (lang > 3 && lang < 7)
			status = get_file_size("ce.txt");
		if (DEBUG)
			printf("status=%d\n", status);
		// docker中不能umount proc
		execute_cmd("/bin/umount -f bin usr lib lib64 etc/alternatives dev 2>&1 >/dev/null");
		execute_cmd("/bin/umount -f %s/* 2>&1 >/dev/null", work_dir);
		umount(work_dir);
		printf("status=%d", status);
		return status;
	}
}
/*
 int read_proc_statm(int pid){
 FILE * pf;
 char fn[4096];
 int ret;
 sprintf(fn,"/proc/%d/statm",pid);
 pf=fopen(fn,"r");
 fscanf(pf,"%d",&ret);
 fclose(pf);
 return ret;
 }
 */
int get_proc_status(int pid, const char *mark)
{
	FILE *pf;
	char fn[BUFFER_SIZE], buf[BUFFER_SIZE];
	int ret = 0;
	sprintf(fn, "/proc/%d/status", pid);
	pf = fopen(fn, "re");
	int m = strlen(mark);
	while (pf && fgets(buf, BUFFER_SIZE - 1, pf))
	{

		buf[strlen(buf) - 1] = 0;
		if (strncmp(buf, mark, m) == 0)
		{
			sscanf(buf + m + 1, "%d", &ret);
		}
	}
	if (pf)
		fclose(pf);
	return ret;
}

// void _get_solution_mysql(int solution_id, char *work_dir, int lang)
// {
// 	char sql[BUFFER_SIZE], src_pth[BUFFER_SIZE];
// 	sprintf(sql, "SELECT source FROM source_code WHERE solution_id=%d",
// 					solution_id);
// 	// create the src file
// 	sprintf(src_pth, "Main.%s", lang_ext[lang]);
// 	if (DEBUG)
// 		printf("Main=%s", src_pth);
// 	FILE *fp_src = fopen(src_pth, "we");
// 	fprintf(fp_src, "%s", row[0]);
// 	if (res != NULL)
// 	{
// 		mysql_free_result(res); // free the memory
// 		res = NULL;
// 	}
// 	fclose(fp_src);
// }

// void get_solution(int solution_id, char *work_dir, int lang)
// {
// 	char src_pth[BUFFER_SIZE];
// 	sprintf(src_pth, "Main.%s", lang_ext[lang]);

// 	_get_solution_mysql(solution_id, work_dir, lang);
// 	execute_cmd("chown judge %s/%s", work_dir, src_pth);
// }

char *escape(char s[], char t[])
{
	int i, j;
	for (i = j = 0; t[i] != '\0'; ++i)
	{
		if (t[i] == '\'')
		{
			s[j++] = '\'';
			s[j++] = '\\';
			s[j++] = '\'';
			s[j++] = '\'';
			continue;
		}
		else
		{
			s[j++] = t[i];
		}
	}
	s[j] = '\0';
	return s;
}

void prepare_files(char *filename, int namelen, char *infile, int p_id,
				   char *work_dir, char *outfile, char *userfile, int runid)
{
	//              printf("ACflg=%d %d check a file!\n",ACflg,solution_id);

	char fname0[BUFFER_SIZE];
	char fname[BUFFER_SIZE];
	strncpy(fname0, filename, namelen);
	fname0[namelen] = 0;
	escape(fname, fname0);
	printf("%s\n%s\n", fname0, fname);
	sprintf(infile, "%s/data/%d/%s.in", oj_home, p_id, fname);
	execute_cmd("/bin/cp '%s' %s/data.in", infile, work_dir);
	execute_cmd("/bin/cp %s/data/%d/*.dic %s/", oj_home, p_id, work_dir);

	sprintf(outfile, "%s/data/%d/%s.out", oj_home, p_id, fname0);
	sprintf(userfile, "%s/run%d/user.out", oj_home, runid);
}

void run_solution(int &lang, char *work_dir, int &time_lmt, int &usedtime, int &mem_lmt)
{
	nice(19);
	int py2 = execute_cmd("/bin/grep 'python2' Main.py");
	chdir(work_dir);
	// 打开文件
	freopen("data.in", "r", stdin);
	freopen("user.out", "w", stdout);
	freopen("error.out", "a+", stderr);
	// trace me
	if (use_ptrace)
		ptrace(PTRACE_TRACEME, 0, NULL, NULL);
	// 语言不是 JAVA
	if (lang != 3)
		chroot(work_dir);
	/*	
		可以这样使用 stuid() 函数：
		开始时，某个程序需要 root 权限玩成一些工作，但后续的工作不需要 root 权限。
		可以将该可执行程序文件设置 set_uid 位，并使得该文件的属主为 root。
		这样，普通用户执行这个程序时，进程就具有了 root 权限，当不再需要 root 权限时，
		调用 setuid( getuid() ) 恢复进程的实际用户 ID 和有效用户 ID 为执行该程序的普通用户的 ID 。
		对于一些提供网络服务的程序，这样做是非常有必要的，否则就可能被攻击者利用，使攻击者控制整个系统。
		对于设置了 set_uid 位的可执行程序也要注意，尤其是对那些属主是 root 的更要注意。因为 Linux 系统中 root 用户拥有最高权力。黑客们往往喜欢寻找设置了 set_uid 位的可执行程序的漏洞。这样的程序如果存在缓冲区溢出漏洞，并且该程序是一个网络程序，那么黑客就可以从远程的地方轻松地利用该漏洞获得运行该漏洞程序的主机的 root 权限。即使这样的成不是网络程序，那么也可以使本机上的恶意普通用户提升为 root 权限。
 		*/
	while (setgid(1536) != 0)
		sleep(1);
	while (setuid(1536) != 0)
		sleep(1);
	while (setresuid(1536, 1536, 1536) != 0)
		sleep(1);
	// 设置资源限制
	struct rlimit LIM;
	if (oi_mode)
	{
		LIM.rlim_cur = time_lmt / cpu_compensation + 1;
	}
	else
	{
		LIM.rlim_cur = (time_lmt / cpu_compensation - usedtime / 1000) + 1;
	}
	LIM.rlim_max = LIM.rlim_cur;
	// 设置CPU时间的最大量值（秒），当超过此软限制时向该进程发送SIGXCPU信号
	setrlimit(RLIMIT_CPU, &LIM);
	alarm(0);
	alarm(time_lmt * 5 / cpu_compensation);

	// 文件限制
	LIM.rlim_max = STD_F_LIM + STD_MB;
	LIM.rlim_cur = STD_F_LIM;
	// 可以创建的文件的最大字节长度，当超过此软限制时向进程发送SIGXFSZ
	setrlimit(RLIMIT_FSIZE, &LIM);
	switch (lang)
	{
	case 5: //bash
		LIM.rlim_cur = LIM.rlim_max = 3;
		break;
	case 9: //C#
	case 17:
		LIM.rlim_cur = LIM.rlim_max = 280;
		break;
	case 3: //java
	case 4: //ruby
	//case 6:  //python
	case 12:
	case 16:
		LIM.rlim_cur = LIM.rlim_max = 80;
		break;
	default:
		LIM.rlim_cur = LIM.rlim_max = 1;
	}
	// 每个实际用户ID所拥有的最大子进程数，更改此限制将影响到sysconf函数在参数_SC_CHILD_MAX中返回的值
	setrlimit(RLIMIT_NPROC, &LIM);
	LIM.rlim_cur = STD_MB << 7;
	LIM.rlim_max = STD_MB << 7;
	// 栈的最大长度
	setrlimit(RLIMIT_STACK, &LIM);
	// set the memory
	LIM.rlim_cur = STD_MB * mem_lmt / 2 * 3;
	LIM.rlim_max = STD_MB * mem_lmt * 2;
	if (lang < 3 || lang == 10 || lang == 13 || lang == 14)
	{
		// 进程总共可用的内存大小的最大值
		setrlimit(RLIMIT_AS, &LIM);
	}
	switch (lang)
	{
	case 0:
	case 1:
	case 2:
	case 10:
	case 11:
	case 13:
	case 14:
	case 17:
		execl("./Main", "./Main", (char *)NULL);
		break;
	case 3:
		sprintf(java_xms, "-Xmx%dM", mem_lmt);
		//sprintf(java_xmx, "-XX:MaxPermSize=%dM", mem_lmt);
		execl("/usr/bin/java", "/usr/bin/java", java_xms, java_xmx,
			  "-Djava.security.manager",
			  "-Djava.security.policy=./java.policy", "Main", (char *)NULL);
		break;
	case 4:
		//system("/ruby Main.rb<data.in");
		execl("/ruby", "/ruby", "Main.rb", (char *)NULL);
		break;
	case 5: //bash
		execl("/bin/bash", "/bin/bash", "Main.sh", (char *)NULL);
		break;
	case 6: //Python
		if (!py2)
		{
			execl("/python2", "/python2", "Main.py", (char *)NULL);
		}
		else
		{
			execl("/python3.7", "/python3.7", "Main.py", (char *)NULL);
		}
		break;
	case 7: //php
		execl("/php", "/php", "Main.php", (char *)NULL);
		break;
	case 8: //perl
		execl("/perl", "/perl", "Main.pl", (char *)NULL);
		break;
	case 9: //Mono C#
		execl("/mono", "/mono", "--debug", "Main.exe", (char *)NULL);
		break;
	case 12: //guile
		execl("/guile", "/guile", "Main.scm", (char *)NULL);
		break;
	case 15: //guile
		execl("/lua", "/lua", "Main", (char *)NULL);
		break;
	case 16: //Node.js
		execl("/nodejs", "/nodejs", "Main.js", (char *)NULL);
		break;
	}
	//sleep(1);
	fflush(stderr);
	write_chroot_log("fflush end");
	// 子线程退出
	exit(0);
}
int fix_python_mis_judge(char *work_dir, int &ACflg, int &topmemory,
						 int mem_lmt)
{
	int comp_res = OJ_AC;

	comp_res = execute_cmd(
		"/bin/grep 'MemoryError'  %s/error.out", work_dir);

	if (!comp_res)
	{
		printf("Python need more Memory!");
		ACflg = OJ_ML;
		topmemory = mem_lmt * STD_MB;
	}

	return comp_res;
}

int fix_java_mis_judge(char *work_dir, int &ACflg, int &topmemory,
					   int mem_lmt)
{
	int comp_res = OJ_AC;
	execute_cmd("chmod 700 %s/error.out", work_dir);
	if (DEBUG)
		execute_cmd("cat %s/error.out", work_dir);
	comp_res = execute_cmd("/bin/grep 'Exception'  %s/error.out", work_dir);
	if (!comp_res)
	{
		printf("Exception reported\n");
		ACflg = OJ_RE;
	}
	execute_cmd("cat %s/error.out", work_dir);

	comp_res = execute_cmd(
		"/bin/grep 'java.lang.OutOfMemoryError'  %s/error.out", work_dir);

	if (!comp_res)
	{
		printf("JVM need more Memory!");
		ACflg = OJ_ML;
		topmemory = mem_lmt * STD_MB;
	}

	if (!comp_res)
	{
		printf("JVM need more Memory or Threads!");
		ACflg = OJ_ML;
		topmemory = mem_lmt * STD_MB;
	}
	comp_res = execute_cmd("/bin/grep 'Could not create'  %s/error.out",
						   work_dir);
	if (!comp_res)
	{
		printf("jvm need more resource,tweak -Xmx(OJ_JAVA_BONUS) Settings");
		ACflg = OJ_RE;
		//topmemory=0;
	}
	return comp_res;
}
int special_judge(char *oj_home, int problem_id, char *infile, char *outfile,
				  char *userfile)
{

	pid_t pid;
	printf("pid=%d\n", problem_id);
	pid = fork();
	int ret = 0;
	if (pid == 0)
	{

		while (setgid(1536) != 0)
			sleep(1);
		while (setuid(1536) != 0)
			sleep(1);
		while (setresuid(1536, 1536, 1536) != 0)
			sleep(1);

		struct rlimit LIM; // time limit, file limit& memory limit

		LIM.rlim_cur = 5;
		LIM.rlim_max = LIM.rlim_cur;
		setrlimit(RLIMIT_CPU, &LIM);
		alarm(0);
		alarm(10);

		// file limit
		LIM.rlim_max = STD_F_LIM + STD_MB;
		LIM.rlim_cur = STD_F_LIM;
		setrlimit(RLIMIT_FSIZE, &LIM);

		ret = execute_cmd("%s/data/%d/spj '%s' '%s' %s", oj_home, problem_id,
						  infile, outfile, userfile);
		if (DEBUG)
			printf("spj1=%d\n", ret);
		if (ret)
			exit(1);
		else
			exit(0);
	}
	else
	{
		int status;

		waitpid(pid, &status, 0);
		ret = WEXITSTATUS(status);
		if (DEBUG)
			printf("spj2=%d\n", ret);
	}
	return ret;
}

void judge_solution(int &ACflg, int &usedtime, int time_lmt, int isspj,
					int p_id, char *infile, char *outfile, char *userfile, int &PEflg,
					int lang, char *work_dir, int &topmemory, int mem_lmt, int num_of_test, char *error_point)
{
	//usedtime-=1000;
	int comp_res;
	if (!oi_mode)
		num_of_test = 1.0;
	if (ACflg == OJ_AC && usedtime > time_lmt * 1000 * (use_max_time ? 1 : num_of_test))
		ACflg = OJ_TL;
	if (topmemory > mem_lmt * STD_MB)
		ACflg = OJ_ML; //issues79
					   // compare
	// 在TL和ML时候也记录错误点
	if (ACflg == OJ_TL || ACflg == OJ_ML)
	{
		// 将错误点信息保存
		strcpy(error_point, getFileNameFromPath(outfile));
	}
	// TOADD
	if (ACflg == OJ_AC)
	{
		if (isspj)
		{
			comp_res = special_judge(oj_home, p_id, infile, outfile, userfile);

			if (comp_res == 0)
				comp_res = OJ_AC;
			else
			{
				if (DEBUG)
					printf("fail test %s\n", infile);
				comp_res = OJ_WA;
			}
		}
		else
		{
			comp_res = compare(outfile, userfile);
		}
		if (comp_res == OJ_WA)
		{
			ACflg = OJ_WA;
			// 在 WA 时记录错误点
			strcpy(error_point, getFileNameFromPath(outfile));
			if (DEBUG)
			{
				printf("fail test %s\n", infile);
			}
		}
		else if (comp_res == OJ_PE)
		{
			PEflg = OJ_PE;
		}
		ACflg = comp_res;
	}
	//jvm popup messages, if don't consider them will get miss-WrongAnswer
	if (lang == 3)
	{
		comp_res = fix_java_mis_judge(work_dir, ACflg, topmemory, mem_lmt);
	}
	if (lang == 6)
	{
		comp_res = fix_python_mis_judge(work_dir, ACflg, topmemory, mem_lmt);
	}
}

int get_page_fault_mem(struct rusage ruse, pid_t pidApp)
{
	//java use pagefault
	int m_vmpeak, m_vmdata, m_minflt;
	m_minflt = ruse.ru_minflt * getpagesize();
	if (0 && DEBUG)
	{
		m_vmpeak = get_proc_status(pidApp, "VmPeak:");
		m_vmdata = get_proc_status(pidApp, "VmData:");
		printf("VmPeak:%d KB VmData:%d KB minflt:%d KB\n", m_vmpeak, m_vmdata,
			   m_minflt >> 10);
	}
	return m_minflt;
}
void print_runtimeerror(char *err)
{
	FILE *ferr = fopen("error.out", "a+");
	fprintf(ferr, "Runtime Error:%s\n", err);
	fclose(ferr);
}

void watch_solution(pid_t pidApp, char *infile, int &ACflg, int isspj,
					char *userfile, char *outfile, int lang,
					int &topmemory, int mem_lmt, int &usedtime, int time_lmt, int p_id,
					int PEflg, char *work_dir)
{
	int tempmemory = 0;
	if (DEBUG)
		printf("pid=%d judging filep:%p %s \n", pidApp, infile, infile);
	int status, sig, exitcode;
	if (DEBUG)
		printf("step1");
	struct user_regs_struct reg;
	if (DEBUG)
		printf("step2");
	// 进程资源信息
	struct rusage ruse;
	if (DEBUG)
		printf("step3");
	int first = 1;
	while (1)
	{
		wait4(pidApp, &status, __WALL, &ruse);
		if (first)
		{
			ptrace(PTRACE_SETOPTIONS, pidApp, NULL, PTRACE_O_TRACESYSGOOD | PTRACE_O_TRACEEXIT
				   //	|PTRACE_O_EXITKILL
				   //	|PTRACE_O_TRACECLONE
				   //	|PTRACE_O_TRACEFORK
				   //	|PTRACE_O_TRACEVFORK
			);
		}
		//jvm gc ask VM before need,so used kernel page fault times and page size
		if (lang == 3 || lang == 7 || lang == 9 || lang == 13 || lang == 14 || lang == 16 || lang == 17)
		{
			tempmemory = get_page_fault_mem(ruse, pidApp);
		}
		else
		{ //other use VmPeak
			tempmemory = get_proc_status(pidApp, "VmPeak:") << 10;
		}
		if (tempmemory > topmemory)
			topmemory = tempmemory;
		if (topmemory > mem_lmt * STD_MB)
		{
			if (DEBUG)
				printf("out of memory %d\n", topmemory);
			if (ACflg == OJ_AC)
				ACflg = OJ_ML;
			ptrace(PTRACE_KILL, pidApp, NULL, NULL);
			break;
		}
		//sig = status >> 8;/*status >> 8 EXITCODE*/
		if (WIFEXITED(status))
			break;
		if ((lang < 4 || lang == 9) && get_file_size("error.out") && !oi_mode)
		{
			ACflg = OJ_RE;
			ptrace(PTRACE_KILL, pidApp, NULL, NULL);
			break;
		}

		if (!isspj && get_file_size(userfile) > get_file_size(outfile) * 2 + 1024)
		{
			ACflg = OJ_OL;
			ptrace(PTRACE_KILL, pidApp, NULL, NULL);
			break;
		}

		exitcode = WEXITSTATUS(status);
		/*exitcode == 5 waiting for next CPU allocation          * ruby using system to run,exit 17 ok
		 *  Runtime Error:Unknown signal xxx need be added here  
                 */
		if ((lang >= 3 && exitcode == 17) || exitcode == 0x05 || exitcode == 0 || exitcode == 133)
		{
			;
		}
		else
		{

			if (DEBUG)
			{
				printf("status>>8=%d\n", exitcode);
			}
			//psignal(exitcode, NULL);

			if (ACflg == OJ_AC)
			{
				switch (exitcode)
				{
				case SIGCHLD:
				case SIGALRM:
					alarm(0);
					if (DEBUG)
						printf("alarm:%d\n", time_lmt);
				case SIGKILL:
				case SIGXCPU:
					ACflg = OJ_TL;
					usedtime = time_lmt * 1000;
					if (DEBUG)
						printf("TLE:%d\n", usedtime);
					break;
				case SIGXFSZ:
					ACflg = OJ_OL;
					break;
				default:
					ACflg = OJ_RE;
				}
				print_runtimeerror(strsignal(exitcode));
			}
			ptrace(PTRACE_KILL, pidApp, NULL, NULL);

			break;
		}
		if (WIFSIGNALED(status))
		{
			/*  WIFSIGNALED: if the process is terminated by signal
			 *
			 *  psignal(int sig, char *s)，like perror(char *s)，print out s, with error msg from system of sig  
			 * sig = 5 means Trace/breakpoint trap
			 * sig = 11 means Segmentation fault
			 * sig = 25 means File size limit exceeded
			 */
			sig = WTERMSIG(status);

			if (DEBUG)
			{
				printf("WTERMSIG=%d\n", sig);
				psignal(sig, NULL);
			}
			if (ACflg == OJ_AC)
			{
				switch (sig)
				{
				case SIGCHLD:
				case SIGALRM:
					alarm(0);
				case SIGKILL:
				case SIGXCPU:
					ACflg = OJ_TL;
					break;
				case SIGXFSZ:
					ACflg = OJ_OL;
					break;

				default:
					ACflg = OJ_RE;
				}
				print_runtimeerror(strsignal(sig));
			}
			break;
		}
		/*     comment from http://www.felix021.com/blog/read.php?1662

		 WIFSTOPPED: return true if the process is paused or stopped while ptrace is watching on it
		 WSTOPSIG: get the signal if it was stopped by signal
		 */

		// check the system calls
		ptrace(PTRACE_GETREGS, pidApp, NULL, &reg);
		call_id = (unsigned int)reg.REG_SYSCALL % call_array_size;
		if (call_counter[call_id])
		{
			//call_counter[reg.REG_SYSCALL]--;
		}
		else if (record_call)
		{
			call_counter[call_id] = 1;
		}
		else
		{ //do not limit JVM syscall for using different JVM
			ACflg = OJ_RE;
			char error[BUFFER_SIZE];
			sprintf(error,
					"[错误] 一个不被允许的系统调用: 系统调用ID:%u\n"
					"要解决这个问题 , 询问管理员将系统调用ID添加到相应的 LANG_XXV[]中（位于okcalls32/64.h 文件中）,\n"
					"然后重新编译判题程序 \n"
					"如果你是管理员，不知道如何处理 ,\n"
					"可以在这里找到中文解释 https://zhuanlan.zhihu.com/p/24498599\n",
					call_id);

			write_log(error);
			print_runtimeerror(error);
			ptrace(PTRACE_KILL, pidApp, NULL, NULL);
		}

		ptrace(PTRACE_SYSCALL, pidApp, NULL, NULL);
		first = 0;
	}
	usedtime += (ruse.ru_utime.tv_sec * 1000 + ruse.ru_utime.tv_usec / 1000) * cpu_compensation;
	usedtime += (ruse.ru_stime.tv_sec * 1000 + ruse.ru_stime.tv_usec / 1000) * cpu_compensation;
}

void mk_shm_workdir(char *work_dir)
{
	char shm_path[BUFFER_SIZE];
	sprintf(shm_path, "/dev/shm/hustoj/%s", work_dir);
	execute_cmd("/bin/mkdir -p %s", shm_path);
	execute_cmd("/bin/ln -s %s %s/", shm_path, oj_home);
	execute_cmd("/bin/chown judge %s ", shm_path);
	execute_cmd("chmod 755 %s ", shm_path);
	//sim need a soft link in shm_dir to work correctly
	sprintf(shm_path, "/dev/shm/hustoj/%s/", oj_home);
	execute_cmd("/bin/ln -s %s/data %s", oj_home, shm_path);
	printf("mk_shm_workdir");
}
int count_in_files(char *dirpath)
{
	const char *cmd = "ls -l %s/*.in|wc -l";
	int ret = 0;
	FILE *fjobs = read_cmd_output(cmd, dirpath);
	fscanf(fjobs, "%d", &ret);
	pclose(fjobs);

	return ret;
}

void print_call_array()
{
	printf("int LANG_%sV[256]={", LANG_NAME);
	int i = 0;
	for (i = 0; i < call_array_size; i++)
	{
		if (call_counter[i])
		{
			printf("%d,", i);
		}
	}
	printf("0};\n");

	printf("int LANG_%sC[256]={", LANG_NAME);
	for (i = 0; i < call_array_size; i++)
	{
		if (call_counter[i])
		{
			printf("HOJ_MAX_LIMIT,");
		}
	}
	printf("0};\n");
}

// ./judge_client /home/judge/hustoj_core 0 1 1 128 0
// ./judge_client /home/judge/hustoj_core 4 1 1 128 0
// 入参 第二个参数 ojhome 第三个参数 问题id 第四个参数 语言 第五个参数 时间限制 第六个参数内存限制 第七个参数 runid
void init_parameters(int argc, char **argv, int &p_id, int &lang, int &time_lmt, int &mem_lmt, int &runid)
{
	DEBUG = (argc > 5);
	strcpy(oj_home, argv[1]);
	p_id = atol(argv[2]);
	lang = atol(argv[3]);
	time_lmt = atol(argv[4]);
	mem_lmt = atol(argv[5]);
	runid = atol(argv[6]);
}

void write_result(int result, int time_used, int memeory_used, double pass_rate, char *error_point)
{
	char src_pth[BUFFER_SIZE];
	sprintf(src_pth, "result.txt");
	FILE *fp_src = fopen(src_pth, "w");
	fprintf(fp_src, "%d\n", result);
	fprintf(fp_src, "%d\n", time_used);
	fprintf(fp_src, "%d\n", memeory_used);
	fprintf(fp_src, "%lf\n", pass_rate);
	fprintf(fp_src, "%s\n", error_point);
}

int main(int argc, char **argv)
{
	int p_id, lang, time_lmt, mem_lmt, runid;
	int isspj, sim, sim_s_id, max_case_time = 0, num_of_test, namelen;
	int usedtime = 0, topmemory = 0;
	int ACflg, PEflg;
	ACflg = PEflg = OJ_AC;
	int finalACflg = ACflg;

	double pass_rate = 0.0;

	init_conf();
	init_parameters(argc, argv, p_id, lang, time_lmt, mem_lmt, runid);
	char work_dir[BUFFER_SIZE];
	sprintf(work_dir, "%s/run%d/", oj_home, runid);
	chdir(work_dir);

	char error_point[BUFFER_SIZE];
	memset(error_point, 0, BUFFER_SIZE);
	// 设置临时工作目录路径
	if (shm_run)
		mk_shm_workdir(work_dir);
	// 如果是测试运行
	if (p_id == 0)
	{
		time_lmt = 5;
		mem_lmt = 128;
		isspj = 0;
	}
	// 对 java 增大运行时间限制
	if (lang >= 3 && lang != 10 && lang != 13 && lang != 14)
	{
		// Clang Clang++ not VM or Script
		// the limit for java
		time_lmt = time_lmt + java_time_bonus;
		mem_lmt = mem_lmt + java_memory_bonus;
		// copy java.policy
		if (lang == 3)
		{
			execute_cmd("/bin/cp %s/etc/java0.policy %s/java.policy", oj_home, work_dir);
			execute_cmd("chmod 755 %s/java.policy", work_dir);
			execute_cmd("chown judge %s/java.policy", work_dir);
		}
	}
	// 时间和内存限制有上限
	if (time_lmt > 300 || time_lmt < 1)
		time_lmt = 300;
	if (mem_lmt > 1024 || mem_lmt < 1)
		mem_lmt = 1024;

	// 编译
	int Compile_OK = compile(lang, work_dir);
	// 如果编译失败了
	if (Compile_OK != 0)
	{
		write_log("compile error");
		// 处理返回值
		write_result(OJ_CE, 0, 0, 0, error_point);
		return 0;
	}
	else
		umount(work_dir);

	char fullpath[BUFFER_SIZE], infile[BUFFER_SIZE], outfile[BUFFER_SIZE], userfile[BUFFER_SIZE];
	sprintf(fullpath, "%s/data/%d", oj_home, p_id); // 数据文件所在路径
	// 打开目录
	DIR *dp;
	struct dirent *dirp;
	// 获取数据文件
	if (p_id > 0 && (dp = opendir(fullpath)) == NULL)
	{
		write_log("No such dir:%s!\n", fullpath);
		return 1;
	}

	//create chroot for ruby bash python
	if (lang == 4)
		copy_ruby_runtime(work_dir);
	if (lang == 5)
		copy_bash_runtime(work_dir);
	if (lang == 6)
		copy_python_runtime(work_dir);
	if (lang == 7)
		copy_php_runtime(work_dir);
	if (lang == 8)
		copy_perl_runtime(work_dir);
	if (lang == 9)
		copy_mono_runtime(work_dir);
	if (lang == 10)
		copy_objc_runtime(work_dir);
	if (lang == 11)
		copy_freebasic_runtime(work_dir);
	if (lang == 12)
		copy_guile_runtime(work_dir);
	if (lang == 15)
		copy_lua_runtime(work_dir);
	if (lang == 16)
		copy_js_runtime(work_dir);
	if (lang == 17)
		copy_go_runtime(work_dir);

	// 如果是测试运行
	if (p_id == 0)
	{
		printf("runid%d,running a custom input...\n", runid);
		init_syscalls_limits(lang);
		pid_t pidApp = fork();
		// 如果是 fork 的子线程
		if (pidApp == 0)
		{
			run_solution(lang, work_dir, time_lmt, usedtime, mem_lmt);
		}
		else
		{
			// 主线程等待线程同步
			watch_solution(pidApp, infile, ACflg, isspj, userfile, outfile, lang, topmemory, mem_lmt, usedtime, time_lmt, p_id, PEflg, work_dir);
		}
		if (finalACflg == OJ_TL)
		{
			usedtime = time_lmt * 1000;
		}
		return 0;
	}
	// 循环每一个测试数据运行测试
	// 保存错误点信息
	for (; (oi_mode || ACflg == OJ_AC || ACflg == OJ_PE) && (dirp = readdir(dp)) != NULL;)
	{
		namelen = isInFile(dirp->d_name); // check if the file is *.in or not
		if (namelen == 0)
			continue;

		prepare_files(dirp->d_name, namelen, infile, p_id, work_dir, outfile,
					  userfile, runid);
		init_syscalls_limits(lang);
		pid_t pidApp = fork();
		if (pidApp == 0)
		{
			run_solution(lang, work_dir, time_lmt, usedtime, mem_lmt);
		}
		else
		{
			num_of_test++;
			watch_solution(pidApp, infile, ACflg, isspj, userfile, outfile, lang, topmemory, mem_lmt, usedtime, time_lmt, p_id, PEflg, work_dir);

			judge_solution(ACflg, usedtime, time_lmt, isspj, p_id, infile, outfile, userfile, PEflg, lang, work_dir, topmemory, mem_lmt, num_of_test, error_point);
			if (use_max_time)
			{
				max_case_time = usedtime > max_case_time ? usedtime : max_case_time;
				usedtime = 0;
			}
		}
		if (oi_mode)
		{
			if (ACflg == OJ_AC)
			{
				++pass_rate;
			}
			if (finalACflg < ACflg)
			{
				finalACflg = ACflg;
			}
			ACflg = OJ_AC;
		}
	}
	if (ACflg == OJ_AC && PEflg == OJ_PE)
		ACflg = OJ_PE;
	// TODO 这里应该做判重的处理

	if ((oi_mode && finalACflg == OJ_RE) || ACflg == OJ_RE)
	{
		if (DEBUG)
			printf("add RE info");
	}

	if (use_max_time)
	{
		usedtime = max_case_time;
	}
	if (finalACflg == OJ_TL)
	{
		usedtime = time_lmt * 1000;
		if (DEBUG)
			printf("usedtime:%d\n", usedtime);
	}
	if (oi_mode)
	{
		if (num_of_test > 0)
			pass_rate /= num_of_test;
	}
	else
	{
		pass_rate = 0;
	}
	// 此处有修改，修改为在TL和ML的时候也输出
	if ((oi_mode && finalACflg == OJ_WA) || (ACflg >= OJ_WA && ACflg <= OJ_ML))
	{
		if (DEBUG)
			printf("add diff info\n");
	}
	if (DEBUG)
		write_log("result=%d", oi_mode ? finalACflg : ACflg);
	if (record_call)
	{
		print_call_array();
	}
	closedir(dp);
	write_log("time_used=%d,memory_used=%d", usedtime, topmemory);
	int finalResult = oi_mode ? finalACflg : ACflg;
	write_result(finalResult, usedtime, topmemory >> 10, pass_rate, error_point);
	return 0;
}
