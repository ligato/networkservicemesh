---
apiVersion: extensions/v1beta1
kind: Deployment
spec:
  replicas: 1
  template:
    metadata:
      labels:
        networkservicemesh.io/app: "vpn-gateway"
        networkservicemesh.io/impl: "secure-intranet-connectivity"
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: networkservicemesh.io/app
                    operator: In
                    values:
                      - vpn-gateway
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: vpn-gateway
          image: {{ .Values.registry }}/networkservicemesh/icmp-responder-nse:{{ .Values.tag }}/
          imagePullPolicy: IfNotPresent
          env:
            - name: ADVERTISE_NSE_NAME
              value: "secure-intranet-connectivity"
            - name: ADVERTISE_NSE_LABELS
              value: "app=vpn-gateway"
            - name: TRACER_ENABLED
              value: "true"
            - name: IP_ADDRESS
              value: "10.60.1.0/24"
          resources:
            limits:
              networkservicemesh.io/socket: 1
        - name: nginx
          image: {{ .Values.registry }}/networkservicemesh/nginx:{{ .Values.tag }}/
metadata:
  name: vpn-gateway-nse
  namespace: default
