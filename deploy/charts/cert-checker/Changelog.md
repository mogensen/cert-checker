# Change Log

## 0.0.6 

**Release date:** 2021-06-21

![AppVersion: v0.0.6](https://img.shields.io/static/v1?label=AppVersion&message=v0.0.6&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Version v0.0.6 
* Update Documentation 
* Adding readinessProbe and livenessProbe to Kubernetes deployment artifacts 

### Default value changes

```diff
diff --git a/deploy/charts/cert-checker/values.yaml b/deploy/charts/cert-checker/values.yaml
index 1462119..ac999ec 100644
--- a/deploy/charts/cert-checker/values.yaml
+++ b/deploy/charts/cert-checker/values.yaml
@@ -23,6 +23,24 @@ serviceAccount:
   # If not set and create is true, a name is generated using the fullname template
   name: ""
 
+## Monitoring
+livenessProbe:
+  enabled: true
+  httpGet:
+    path: /
+    port: 8081
+  initialDelaySeconds: 2
+  periodSeconds: 3
+
+readinessProbe:
+  enabled: true
+  httpGet:
+    path: /
+    port: 8081
+  initialDelaySeconds: 2
+  periodSeconds: 3
+
+## Config
 certchecker:
   loglevel: info
   intervalminutes: 1
```

## 0.0.5 

**Release date:** 2021-06-15

![AppVersion: v0.0.5](https://img.shields.io/static/v1?label=AppVersion&message=v0.0.5&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Version v0.0.5 
* Update servicemonitor.yaml 

### Default value changes

```diff
# No changes in this release
```

## 0.0.4 

**Release date:** 2021-05-21

![AppVersion: v0.0.4](https://img.shields.io/static/v1?label=AppVersion&message=v0.0.4&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Upgrade version to v0.0.4 
* Update Documentation 
* Moving html and css into embedded files 
* Adding Skaffold for quick helm development 
* Use UID over 10.000 to not clash with host's UID 

### Default value changes

```diff
diff --git a/deploy/charts/cert-checker/values.yaml b/deploy/charts/cert-checker/values.yaml
index c2961a5..1462119 100644
--- a/deploy/charts/cert-checker/values.yaml
+++ b/deploy/charts/cert-checker/values.yaml
@@ -53,9 +53,11 @@ podAnnotations:
   prometheus.io/port: "8080"
   prometheus.io/scrape: "true"
   enable.cert-checker.io/cert-checker: "true"
+  # If you want apparmor security
+  # container.apparmor.security.beta.kubernetes.io/cert-checker: runtime/default
 
 podSecurityContext:
-  fsGroup: 2000
+  fsGroup: 35212
 
 securityContext:
   privileged: false
@@ -64,13 +66,34 @@ securityContext:
     - ALL
   readOnlyRootFilesystem: true
   runAsNonRoot: true
-  runAsUser: 1000
+  runAsUser: 35212
   allowPrivilegeEscalation: false
 
 service:
   type: ClusterIP
   port: 8080
 
+uiService:
+  type: ClusterIP
+  port: 8081
+
+ingress:
+  enabled: false
+  annotations: {}
+    # kubernetes.io/ingress.class: nginx
+    # kubernetes.io/tls-acme: "true"
+  hosts:
+    - host: cert-checker.local
+      paths:
+      - path: /
+        backend:
+          serviceName: cert-checker.local
+          servicePort: 8081
+  tls: []
+  #  - secretName: chart-example-tls
+  #    hosts:
+  #      - cert-checker.local
+
 resources: {}
   # limits:
   #   cpu: 100m
```

## 0.0.3 

**Release date:** 2021-03-25

![AppVersion: v0.0.3](https://img.shields.io/static/v1?label=AppVersion&message=v0.0.3&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Release version v0.0.3 to also release helm chart 
* Update Documentation 

### Default value changes

```diff
# No changes in this release
```

## 0.0.2 

**Release date:** 2021-03-25

![AppVersion: v0.0.2](https://img.shields.io/static/v1?label=AppVersion&message=v0.0.2&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Fix servicemonitor in k8s/yaml. Bump to v0.0.2 
* Update Documentation 
* Started on generating k8s yaml from helm chart 
* Adding minimum TLS to metrics and dashboards 
* Update Documentation 
* Updating and cleanup of Kubernetes deployments 

### Default value changes

```diff
diff --git a/deploy/charts/cert-checker/values.yaml b/deploy/charts/cert-checker/values.yaml
index 59f56e4..c2961a5 100644
--- a/deploy/charts/cert-checker/values.yaml
+++ b/deploy/charts/cert-checker/values.yaml
@@ -8,7 +8,7 @@ image:
   repository: mogensen/cert-checker
   pullPolicy: IfNotPresent
   # Overrides the image tag whose default is the chart appVersion.
-  tag: "v0.0.1"
+  # tag: ""
 
 imagePullSecrets: []
 nameOverride: ""
@@ -29,6 +29,7 @@ certchecker:
   certificates:
     - dns: google.com
     - dns: example.com
+    - dns: twitter.com
     - dns: expired.badssl.com
     - dns: wrong.host.badssl.com
     - dns: untrusted-root.badssl.com
@@ -39,8 +40,13 @@ certchecker:
     - dns: null.badssl.com
     - dns: rc4-md5.badssl.com
     - dns: rc4.badssl.com
-  
-serviceMonitor: false
+
+serviceMonitor:
+  enabled: false
+  additionalLabels: {}
+
+grafanaDashboard:
+  enabled: false
 
 podAnnotations:
   prometheus.io/path: /metrics
@@ -48,11 +54,11 @@ podAnnotations:
   prometheus.io/scrape: "true"
   enable.cert-checker.io/cert-checker: "true"
 
-
 podSecurityContext:
   fsGroup: 2000
 
 securityContext:
+  privileged: false
   capabilities:
     drop:
     - ALL
```

## 0.0.1 

**Release date:** 2021-01-29

![AppVersion: 0.0.1](https://img.shields.io/static/v1?label=AppVersion&message=0.0.1&color=success&logo=)
![Helm: v3](https://img.shields.io/static/v1?label=Helm&message=v3&color=informational&logo=helm)


* Add helm chart. 

### Default value changes

```diff
# Default values for cert-checker.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 1

image:
  repository: mogensen/cert-checker
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: "v0.0.1"

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # Annotations to add to the service account
  annotations: {}
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name: ""

certchecker:
  loglevel: info
  intervalminutes: 1
  certificates:
    - dns: google.com
    - dns: example.com
    - dns: expired.badssl.com
    - dns: wrong.host.badssl.com
    - dns: untrusted-root.badssl.com
    - dns: self-signed.badssl.com
    - dns: revoked.badssl.com
    - dns: dh480.badssl.com
    - dns: dh512.badssl.com
    - dns: null.badssl.com
    - dns: rc4-md5.badssl.com
    - dns: rc4.badssl.com
  
serviceMonitor: false

podAnnotations:
  prometheus.io/path: /metrics
  prometheus.io/port: "8080"
  prometheus.io/scrape: "true"
  enable.cert-checker.io/cert-checker: "true"


podSecurityContext:
  fsGroup: 2000

securityContext:
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
  runAsNonRoot: true
  runAsUser: 1000
  allowPrivilegeEscalation: false

service:
  type: ClusterIP
  port: 8080

resources: {}
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}
```

---
Autogenerated from Helm Chart and git history using [helm-changelog](https://github.com/mogensen/helm-changelog)
