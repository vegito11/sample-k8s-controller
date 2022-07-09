## How to Run Code

1. Create the Configmap and deployment

```bash
kubectl apply -f manifests
```

2. Check the file inside pod
```bash
kubectl exec -it deploy/nginx-reload -- cat /data/config.ini
kubectl exec -it deploy/nginx-reload2 -- cat /data/config.ini
kubectl exec -it deploy/nginx-normal -- cat /data/config.ini
```
3. Update the Confimap and see the values
```bash
kubectl apply -f manifests/common_cm.yaml
kubectl exec -it deploy/nginx-reload -- cat /data/config.ini
kubectl exec -it deploy/nginx-reload2 -- cat /data/config.ini
kubectl exec -it deploy/nginx-normal -- cat /data/config.ini
```