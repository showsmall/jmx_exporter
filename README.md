#### 使用文档详细说明
###### 功能说明
- 支持java的full gc运行次数，并且暴露metrics
- 支持java的线程数量，并且暴露metrics

###### 部署说明
- 使用systemd进行管理启动停止以及重启
- 日志默认输出到控制台，使用systemd管理后，日志输出到系统日志中，centos7日志在： /var/log/messages
- jmx_exporter 需要依赖jolokia 部署，jolokia的部署不再说明
- 配置文件：config.yaml放在和jxm_exporter同级目录
- 如果需要采集其他指标需要提供对应的mbean
- jmx_exporter.service 文件说明
```shell
vim /etc/systemd/system/jmx_exporter.service
[Unit]
Description=jmx_exporter
After=network.target
[Service]
Type=simple
WorkingDirectory=/Data/monitor/jmx_exporter 
User=root
Group=root
ExecStart=/Data/monitor/jmx_exporter/jmx_exporter
Restart=on-failure
RestartSec=5
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target

systemctl daemon-reload
# 启动jmx_exporter
systemctl start jmx_exporter.service
# 开启开机自动启动
systemctl enable jmx_exporter.service
# 查看运行状态
systemctl status jmx_exporter.service
```
##### prometheus 配置
- prometheus采集指标时候默认发出的是GET请求，参数携带在url中：http://192.168.8.88:9118/jmx?token=full@gc12AO&gc=a&target=172.16.5.30:1099
- jmx_exporter接口返回metrics指标如下：如果需要支持其它指标，请自行修代码或者提issues
```text
# HELP full_gc_count jvm暴露的jmx 统计自jvm运行起来后full gc的次数.
# TYPE full_gc_count counter
full_gc_count 469
# HELP jvm_ThreadCount jvm业务运行中的总线程，不包含gc线程.
# TYPE jvm_ThreadCount counter
jvm_ThreadCount 512
```
- targets填写需要监控的jvm节点的ip和jmx的端口1099，如：172.16.10.3:1099
```yaml
  - job_name: 'jvm_gc'
    metrics_path: /jmx
    scrape_interval: 5s
    scrape_timeout: 5s
    params:
     token: [full@gc12AO]
     gc: [a]
    static_configs:
      - targets: ['172.16.5.10:1099', '172.16.5.30:1099', '172.16.5.9:1099']
        labels:
          instance: 'jvm_gc'
          group: 'port'
    relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      replacement: 192.168.8.88:9118
```
- prometheus报警规则根据自己业务需要进行配置：如1小时内full gc次数增长不可以超过1次，jvm的线程不能超过指定值等
- 出现异常报警的时候可以调用远程dump然后执行相应的策率进行服务修复，dump是为了捕捉现场，在dump前可以切掉异常jvm节点业务流量，也可以dump后切，因为有可能在没有业务流量后不能捕获到现场
- 远程dump和执行jvm节点故障策率不在这个组件实现
##### 使用的框架和指标说明

- 代码使用了gin框架实现暴露metrics
- 未使用prometheus官方的client_golang

##### 特别注意
- 使用前请认真阅读配置文件内容：config.yaml，并且根据配置文件说明修改


#### ====================================================English=======================================================

#### Detailed Documentation
###### Feature Description
- Supports the number of times Java's Full GC runs, and exposes metrics.
- Supports the number of Java threads, and exposes metrics.

###### Deployment Instructions
- Managed by systemd for starting, stopping, and restarting.
- Logs are output to the console by default. When managed by systemd, logs are output to the system logs, in CentOS 7 located at: /var/log/messages.
- JMX Exporter requires Jolokia for deployment, the deployment of Jolokia is not covered here.
- Configuration file: config.yaml is placed in the same directory as jmx_exporter.
- To collect additional metrics, the corresponding MBean needs to be provided.
- jmx_exporter.service file description
```shell
vim /etc/systemd/system/jmx_exporter.service
[Unit]
Description=jmx_exporter
After=network.target
[Service]
Type=simple
WorkingDirectory=/Data/monitor/jmx_exporter 
User=root
Group=root
ExecStart=/Data/monitor/jmx_exporter/jmx_exporter
Restart=on-failure
RestartSec=5
LimitNOFILE=65536
[Install]
WantedBy=multi-user.target

systemctl daemon-reload
# To start jmx_exporter
systemctl start jmx_exporter.service
# To enable auto-start on boot
systemctl enable jmx_exporter.service
# To check the running status
systemctl status jmx_exporter.service
```
##### Prometheus Configuration
- Prometheus sends a GET request by default when collecting metrics, with parameters carried in the URL: http://192.168.8.88:9118/jmx?token=full@gc12AO&gc=a&target=172.16.5.30:1099
- JMX Exporter returns metrics as follows: To support other metrics, please modify the code or submit issues.
```text
# HELP full_gc_count Exposes JVM's JMX counting the number of Full GCs since the JVM started.
# TYPE full_gc_count counter
full_gc_count 469
# HELP jvm_ThreadCount Total threads in JVM business operations, excluding GC threads.
# TYPE jvm_ThreadCount counter
jvm_ThreadCount 512
```
- Fill in the IP and JMX port 1099 of the JVM nodes to be monitored in targets, e.g.: 172.16.10.3:1099
```yaml
  - job_name: 'jvm_gc'
    metrics_path: /jmx
    scrape_interval: 5s
    scrape_timeout: 5s
    params:
     token: [full@gc12AO]
     gc: [a]
    static_configs:
      - targets: ['172.16.5.10:1099', '172.16.5.30:1099', '172.16.5.9:1099']
        labels:
          instance: 'jvm_gc'
          group: 'port'
    relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      replacement: 192.168.8.88:9118
```
- Prometheus alert rules should be configured according to your own business needs: e.g., the number of Full GCs should not increase more than once in an hour, the number of JVM threads should not exceed a specified value, etc.
- In the event of an abnormal alert, you can call remote dump and then execute the appropriate strategy for service repair. Dumping is to capture the scene. Before dumping, you can cut off the abnormal JVM node business traffic, or you can cut after dumping because it may not be possible to capture the scene without business traffic.
- Remote dump and execution of JVM node fault strategy are not implemented in this component.
##### Framework and Metric Description

- The code uses the Gin framework to expose metrics.
- It does not use the official Prometheus client_golang.

##### Special Attention
- Please read the configuration file content: config.yaml carefully before use, and modify according to the configuration file instructions.
