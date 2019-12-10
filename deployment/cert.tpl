apiVersion: networking.gke.io/v1beta1
kind: ManagedCertificate
metadata:
  name: read-it-later-certificate
spec:
  domains:
    - read-it-later.app
