package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sre/jmx_exporter/config"
	"sre/jmx_exporter/logic"
	"strings"
)

func HandleJmx(c *gin.Context) {
	jmx_host_port := c.Request.URL.Query().Get("target")
	jmx_gc_args := c.Request.URL.Query().Get("gc")
	if jmx_gc_args == "" && jmx_host_port == "" {
		c.JSON(500, gin.H{
			"msg": "请求参数缺少target或者gc",
		})
		return
	}
	if config.C.GetString(fmt.Sprintf("gc.%s", jmx_gc_args)) == "" {
		c.JSON(500, gin.H{
			"msg": "配置文件中没有找到对应的gc算法",
		})
		return
	}

	// 组装获取fullgc post数据
	jmxadd := fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s/jmxrmi", jmx_host_port)
	target := make(map[string]string)
	target["url"] = jmxadd
	postdatamap := make(map[string]interface{})
	postdatamap["type"] = "read"
	postdatamap["mbean"] = config.C.GetString(fmt.Sprintf("gc.%s", jmx_gc_args))
	postdatamap["target"] = target
	r, err := logic.GetJmxData(&postdatamap)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		return
	}

	metrics_string_gc := `# HELP full_gc_count jvm暴露的jmx 统计自jvm运行起来后full gc的次数.
# TYPE full_gc_count counter
`
	var gc strings.Builder
	gc.WriteString(metrics_string_gc)
	gc.WriteString(fmt.Sprintf("full_gc_count %d\n", r.Value.CollectionCount))

	// 组装获取thread post数据
	postdatamap["mbean"] = config.C.GetString("thread.mbean")

	thread, err := logic.GetJmxData(&postdatamap)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		return
	}
	metrics_string_th := `# HELP jvm_ThreadCount jvm业务运行中的总线程，不包含gc线程.
# TYPE jvm_ThreadCount counter
`
	var th strings.Builder
	th.WriteString(metrics_string_th)
	th.WriteString(fmt.Sprintf("jvm_ThreadCount %d\n", thread.Value.ThreadCount))
	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.Status(200)

	_, _ = c.Writer.WriteString(gc.String())
	_, _ = c.Writer.WriteString(th.String())

	c.Writer.Flush()
	return
}
