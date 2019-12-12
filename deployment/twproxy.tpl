apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: twproxy
  labels:
    app: twproxy 
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: twproxy
    spec:
      containers:
      - name: twproxy
        image: gcr.io/${PROJECT_ID}/mailme_app_twproxy:${TWPROXY_VERSION}
        command: ["/app/mailmeapp"]
        args: ["run-tw-proxy", "--tw-consumer-key=${TW_KEY}", "--tw-consumer-secret=${TW_SECRET}"]
        env:
         - name: MAILME_APP_PEM_FILE
           value: "/app/service.pem"
         - name: MAILME_APP_KEY_FILE
           value: "/app/service.key"
        ports:
        - name: grpc
          containerPort: 5000
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "80Mi"
            cpu: "110m"

---
apiVersion: v1
kind: Service
metadata:
  name: twproxy
spec:
  type: ClusterIP
  selector:
    app: twproxy
  ports:
  - name: grpc
    port: 5000
    targetPort: 5000
