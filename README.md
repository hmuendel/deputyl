# deputyl
Deputyl is your little helper, that runs inside a kubernetes cluster,
discovers all running pod images and tags and checks the docker registry for
newer semver tags. It then reports newer versions as a prometheus metric for 
each container.

```
deputyl_newer_patch_versions{app="prometheus",chart="prometheus-10.4.0",component="kube-state-metrics",container="auditwin-web",container_id="docker://d2bdfeff0d1db5aa3b5b02d6798aca0e27500e46a5992e536b3f91d1d565caa9",heritage="Helm",image="concourse/concourse:5.8.0",image_id="docker-pullable://concourse/concourse@sha256:b50f6207373cc671965ea0c4832bb394508b7eca337a80100d38d961be27fbb5",instance="10.5.9.48:8080",job="kubernetes-service-endpoints",kubernetes_name="prometheus-kube-state-metrics",kubernetes_namespace="kube-system",kubernetes_node="ip-10-5-34-71.eu-west-1.compute.internal",namespace="auditwin",pod="auditwin-web-867b696b5b-55tbf",release="prometheus"} 3

deputyl_newer_minor_verions{app="prometheus",chart="prometheus-10.4.0",component="kube-state-metrics",container="auditwin-web",container_id="docker://d2bdfeff0d1db5aa3b5b02d6798aca0e27500e46a5992e536b3f91d1d565caa9",heritage="Helm",image="concourse/concourse:5.8.0",image_id="docker-pullable://concourse/concourse@sha256:b50f6207373cc671965ea0c4832bb394508b7eca337a80100d38d961be27fbb5",instance="10.5.9.48:8080",job="kubernetes-service-endpoints",kubernetes_name="prometheus-kube-state-metrics",kubernetes_namespace="kube-system",kubernetes_node="ip-10-5-34-71.eu-west-1.compute.internal",namespace="auditwin",pod="auditwin-web-867b696b5b-55tbf",release="prometheus"} 1 

deputyl_newer_major_versions{app="prometheus",chart="prometheus-10.4.0",component="kube-state-metrics",container="auditwin-web",container_id="docker://d2bdfeff0d1db5aa3b5b02d6798aca0e27500e46a5992e536b3f91d1d565caa9",heritage="Helm",image="concourse/concourse:5.8.0",image_id="docker-pullable://concourse/concourse@sha256:b50f6207373cc671965ea0c4832bb394508b7eca337a80100d38d961be27fbb5",instance="10.5.9.48:8080",job="kubernetes-service-endpoints",kubernetes_name="prometheus-kube-state-metrics",kubernetes_namespace="kube-system",kubernetes_node="ip-10-5-34-71.eu-west-1.compute.internal",namespace="auditwin",pod="auditwin-web-867b696b5b-55tbf",release="prometheus"} 0



```
