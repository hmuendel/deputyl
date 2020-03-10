# deputyl
Deputyl is your little helper, that runs inside a kubernetes cluster,
discovers all running pod images and tags and checks the docker registry for
newer semver tags. It then reports newer versions as a prometheus metric for 
each container.

```
deputyl_newer_patch_versions{namespace="default",pod="server",repository="nginx:16.0.0"} 3

deputyl_newer_minor_verions{namespace="default",pod="server",repository="nginx:16.0.0"} 1 

deputyl_newer_major_versions{namespace="default",pod="server",repository="nginx:16.0.0"} 0



```
