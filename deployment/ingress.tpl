apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: main-ingress
  annotations:
    kubernetes.io/ingress.global-static-ip-name: read-it-later
    networking.gke.io/managed-certificates: read-it-later-certificate
spec:
  rules:
  - http:
      paths:
      - path: /*
        backend:
          serviceName: web
          servicePort: 8080
      - path: /api/*
        backend:
          serviceName: backend
          servicePort: 8000
      - path: /oauth/*
        backend:
          serviceName: backend
          servicePort: 8000
      - path: /confirm-email
        backend:
          serviceName: backend
          servicePort: 8000
