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
      - name: backend
        image: gcr.io/${PROJECT_ID}/mailme_app_backend:${BACKEND_VERSION}
        command: ["/app/mailmeapp"]
        args: ["--tw-consumer-key=${TW_KEY}", "--tw-consumer-secret=${TW_SECRET}", "--auth-key=${AUTH_KEY}", "--encrypt-key=${ENCRYPT_KEY}"]
        readinessProbe:
          initialDelaySeconds: 10
          httpGet:
            path: "/healthcheck"
            port: 8000
        livenessProbe:
          initialDelaySeconds: 10
          httpGet:
            path: "/healthcheck"
            port: 8000
        env:
         - name: MAILME_APP_PORT
           value: "8000"
         - name: MAILME_APP_DEBUG
           value: "0"
         - name: MAILME_APP_HOST
           value: "0.0.0.0"
         - name: MAILME_APP_HTTP_ONLY
           value: "1"
         - name: MAILME_APP_MAX_AGE
           value: "43200"
         - name: MAILME_APP_TW_CALLBACK_URL
           value: "https://localhost/oauth/tw/callback"
         - name: MAILME_APP_TW_PROXY_HOST
           value: "twproxy" 
         - name: MAILME_APP_DSN
           valueFrom:
             secretKeyRef:
               name: dsn
               key: dsn
        ports:
        - containerPort: 8000

      - name: frontend
        image: gcr.io/${PROJECT_ID}/mailme_app_frontend:${FRONTEND_VERSION}
        ports:
        - containerPort: 8080

      - name: cloudsql-proxy
        image: gcr.io/cloudsql-docker/gce-proxy:1.16
        command: ["/cloud_sql_proxy", "--dir=/cloudsql",
        	  "-instances=${CLOUD_SQL_CONNECTION}=tcp:5432",
        	  "-credential_file=/secrets/cloudsql/credentials.json"]
        volumeMounts:
          - name: cloudsql-oauth-credentials
            mountPath: /secrets/cloudsql
            readOnly: true
          - name: ssl-certs
            mountPath: /etc/ssl/certs
          - name: cloudsql
            mountPath: /cloudsql

      volumes:
       - name: cloudsql-oauth-credentials
         secret:
          secretName: cloudsql-oauth-credentials
       - name: ssl-certs
         hostPath:
          path: /etc/ssl/certs
       - name: cloudsql
         emptyDir:

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
    port: 80
    targetPort: 8080
