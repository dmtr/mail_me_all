apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: main-ingress
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
