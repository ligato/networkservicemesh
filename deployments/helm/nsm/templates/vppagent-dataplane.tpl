{{ $fp := .Values.forwardingPlane }}

apiVersion: apps/v1
kind: DaemonSet
spec:
  selector:
    matchLabels:
      app: nsm-{{ $fp }}-plane
  template:
    metadata:
      labels:
        app: nsm-{{ $fp }}-plane
    spec:
      hostPID: true
      hostNetwork: true
      serviceAccount: forward-plane-acc
      containers:
        - name: {{ (index .Values $fp).image }}
          securityContext:
            privileged: true
          image: {{ .Values.registry }}/{{ .Values.org }}/{{ (index .Values $fp).image }}:{{ .Values.tag }}
          imagePullPolicy: {{ .Values.pullPolicy }}
          env:
            - name: INSECURE
{{- if .Values.insecure }}
              value: "true"
{{- else }}
              value: "false"
{{- end }}
            - name: NSM_DATAPLANE_SRC_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          volumeMounts:
            - name: workspace
              mountPath: /var/lib/networkservicemesh/
              mountPropagation: Bidirectional
            - name: spire-agent-socket
              mountPath: /run/spire/sockets
              readOnly: true
          livenessProbe:
            httpGet:
              path: /liveness
              port: 5555
            initialDelaySeconds: 10
            periodSeconds: 10
            timeoutSeconds: 3
          readinessProbe:
            httpGet:
              path: /readiness
              port: 5555
            initialDelaySeconds: 10
            periodSeconds: 10
            timeoutSeconds: 3
          {{- if (index .Values $fp).resources }}
            resources:
              limits:
                cpu: {{ (index .Values $fp).resources.limitCPU }}
              requests:
                cpu: {{ (index .Values $fp).resources.requestsCPU }}
          {{- end }}
      volumes:
        - hostPath:
            path: /var/lib/networkservicemesh
            type: DirectoryOrCreate
          name: workspace
        - hostPath:
            path: /run/spire/sockets
            type: DirectoryOrCreate
          name: spire-agent-socket
metadata:
  name: nsm-{{ $fp }}-forwarder
  namespace: {{ .Release.Namespace }}
