apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: spire-server
  namespace: {{ .Release.Namespace }}
  labels:
    app: spire-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: spire-server
  serviceName: spire-server
  template:
    metadata:
      namespace: spire
      labels:
        app: spire-server
    spec:
      serviceAccountName: spire-server
      containers:
        - name: spire-server
          image: {{ .Values.registry }}/{{ .Values.org }}/spire-registration:{{ .Values.tag }}
          imagePullPolicy: {{ .Values.pullPolicy }}
          ports:
            - containerPort: 8081
          volumeMounts:
            - name: spire-config
              mountPath: /run/spire/config
              readOnly: true
            - name: spire-entries
              mountPath: /run/spire/entries
              readOnly: true
            - name: spire-secrets
              mountPath: /run/spire/secrets
              readOnly: true
            - name: spire-data
              mountPath: /run/spire/data
              readOnly: false
          livenessProbe:
            tcpSocket:
              port: 8081
            failureThreshold: 2
            initialDelaySeconds: 15
            periodSeconds: 60
            timeoutSeconds: 3
      volumes:
        - name: spire-config
          configMap:
            name: spire-server
        - name: spire-entries
          configMap:
            name: spire-entries
        - name: spire-secrets
          secret:
            secretName: spire-server
        - name: spire-data
          hostPath:
            path: /var/spire-data
            type: DirectoryOrCreate
