#/bin/bash
scp dist/* root@172.16.0.3:/root/ahpuojv2docker/web/
scp dist/admin_index.html root@172.16.0.3:/root/ahpuojv2docker/web/
scp -r dist root@172.16.0.3:/root/ahpuojv2docker/web/
# scp -r static root@172.16.0.3:/root/ahpuojv2docker/web/
# scp -r config root@172.16.0.3:/root/ahpuojv2docker/web/
# scp -r core root@172.16.0.3:/root/ahpuojv2docker/
# scp ahpuoj root@172.16.0.3:/root/ahpuojv2docker/web/ahpuoj
