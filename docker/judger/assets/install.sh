#!/bin/bash

#设置判题机容器运行环境的脚本 生产与部署环境下通用的脚本
#!/bin/bash

# 创建用户
/usr/sbin/useradd -m -u 1536 judge

# 创建评测机目录文件夹
cd /home/judge/hustoj_core
mkdir log data run0 run1 run2 run3
chown judge run0 run1 run2 run3

# 赋予评测机源码部署脚本执行权限并部署
chmod +x /home/judge/core/make.sh
cd /home/judge/core && ./make.sh
cd /usr/bin && rm awk && cp -s mawk awk

# 权限设置
chmod 775 -R /home/judge/hustoj_core/data
chown -r judge -R /home/judge/hustoj_core/data

# 修改配置文件
CPU=$(grep "cpu cores" /proc/cpuinfo | head -1 | awk '{print $4}')
sed -i "s/OJ_RUNNING=1/OJ_RUNNING=$CPU/g" /home/judge/hustoj_core/etc/judge.conf
