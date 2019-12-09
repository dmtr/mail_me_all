apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: web
  labels:
    app: web
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: frontend
        image: gcr.io/${PROJECT_ID}/mailme_app_frontend:${FRONTEND_VERSION}
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: web
spec:
  type: ClusterIP
  selector:
    app: web
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  type: NodePort
