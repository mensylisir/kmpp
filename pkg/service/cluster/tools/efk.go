package tools

import (
	"encoding/json"
	"fmt"

	"github.com/kmpp/pkg/constant"
	"github.com/kmpp/pkg/model"
)

type EFK struct {
	Cluster             *Cluster
	Tool                *model.ClusterTool
	LocalHostName       string
	LocalRepositoryPort int
}

func NewEFK(cluster *Cluster, tool *model.ClusterTool) (*EFK, error) {
	p := &EFK{
		Tool:                tool,
		Cluster:             cluster,
		LocalHostName:       constant.LocalRepositoryDomainName,
		LocalRepositoryPort: cluster.helmRepoPort,
	}
	return p, nil
}

func (e EFK) setDefaultValue(toolDetail model.ClusterToolDetail, isInstall bool) {
	imageMap := map[string]interface{}{}
	_ = json.Unmarshal([]byte(toolDetail.Vars), &imageMap)

	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(e.Tool.Vars), &values)
	values["fluentd-elasticsearch.image.repository"] = fmt.Sprintf("%s:%d/%s", e.LocalHostName, e.LocalRepositoryPort, imageMap["fluentd_image_name"])
	values["fluentd-elasticsearch.imageTag"] = imageMap["fluentd_image_tag"]
	values["elasticsearch.image"] = fmt.Sprintf("%s:%d/%s", e.LocalHostName, e.LocalRepositoryPort, imageMap["elasticsearch_image_name"])
	values["elasticsearch.imageTag"] = imageMap["elasticsearch_image_tag"]

	if isInstall {
		if _, ok := values["elasticsearch.esJavaOpts.item"]; !ok {
			values["elasticsearch.esJavaOpts.item"] = 1
		}
		values["elasticsearch.esJavaOpts"] = fmt.Sprintf("-Xmx%vg -Xms%vg", values["elasticsearch.esJavaOpts.item"], values["elasticsearch.esJavaOpts.item"])
		delete(values, "elasticsearch.esJavaOpts.item")

		if _, ok := values["elasticsearch.volumeClaimTemplate.resources.requests.storage"]; ok {
			values["elasticsearch.volumeClaimTemplate.resources.requests.storage"] = fmt.Sprintf("%vGi", values["elasticsearch.volumeClaimTemplate.resources.requests.storage"])
		}
	}
	str, _ := json.Marshal(&values)
	e.Tool.Vars = string(str)
}

func (e EFK) Install(toolDetail model.ClusterToolDetail) error {
	e.setDefaultValue(toolDetail, true)
	if err := installChart(e.Cluster.HelmClient, e.Tool, constant.LoggingChartName, toolDetail.ChartVersion); err != nil {
		return err
	}
	if err := createRoute(e.Cluster.Namespace, constant.DefaultLoggingIngressName, constant.DefaultLoggingIngress, constant.DefaultLoggingServiceName, 9200, e.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForStatefulSetsRunning(e.Cluster.Namespace, constant.DefaultLoggingStateSetsfulName, 1, e.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (e EFK) Upgrade(toolDetail model.ClusterToolDetail) error {
	e.setDefaultValue(toolDetail, false)
	return upgradeChart(e.Cluster.HelmClient, e.Tool, constant.LoggingChartName, toolDetail.ChartVersion)
}

func (e EFK) Uninstall() error {
	return uninstall(e.Cluster.Namespace, e.Tool, constant.DefaultLoggingIngressName, e.Cluster.HelmClient, e.Cluster.KubeClient)
}
