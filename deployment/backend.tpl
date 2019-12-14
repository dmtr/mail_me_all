apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: backend
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
         - name: MAILME_APP_DOMAIN
           value: "read-it-later.app"
         - name: MAILME_APP_PORT
           value: "8000"
         - name: MAILME_APP_DEBUG
           value: "0"
         - name: MAILME_APP_HOST
           value: "0.0.0.0"
         - name: MAILME_APP_HTTP_ONLY
           value: "1"
         - name: MAILME_APP_SECURE
           value: "1"
         - name: MAILME_APP_MAX_AGE
           value: "43200"
         - name: MAILME_APP_TW_CALLBACK_URL
           value: "https://read-it-later.app/oauth/tw/callback"
         - name: MAILME_APP_TW_PROXY_HOST
           value: "twproxy" 
         - name: MAILME_APP_PEM_FILE
           value: "/app/service.pem"
         - name: MAILME_APP_DSN
           valueFrom:
             secretKeyRef:
               name: dsn
               key: dsn
        ports:
        - containerPort: 8000
        resources:
          requests:
            memory: "64Mi"
            cpu: "80m"
          limits:
            memory: "128Mi"
            cpu: "100m"

      - name: crontasks
        image: gcr.io/${PROJECT_ID}/mailme_app_crontasks:${CRONTASKS_VERSION}
        env:
         - name: MAILME_APP_TW_PROXY_HOST
           value: "twproxy" 
         - name: MAILME_APP_PEM_FILE
           value: "/app/service.pem"
         - name: MAILME_APP_KEY_FILE
           value: "/app/service.key"
         - name: MAILME_APP_TEMPLATE_PATH
           value: "/app/templates/"
         - name: MAILME_APP_MG_DOMAIN
           valueFrom:
             secretKeyRef:
               name: mgdomain
               key: mgdomain
         - name: MAILME_APP_MG_APIKEY
           valueFrom:
             secretKeyRef:
               name: mgapikey
               key: mgapikey 
         - name: MAILME_APP_FROM
           valueFrom:
             secretKeyRef:
               name: mgfrom
               key: mgfrom 
         - name: MAILME_APP_DSN
           valueFrom:
             secretKeyRef:
               name: dsn
               key: dsn
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "110m"

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
  name: backend
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - name: http
    port: 8000
    targetPort: 8000
  type: NodePort
