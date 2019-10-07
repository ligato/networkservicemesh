apiVersion: v1
kind: ConfigMap
metadata:
  name: spire-bundle
  namespace: {{ .Release.Namespace }}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: spire-server
  namespace: {{ .Release.Namespace }}
data:
  server.conf: |
    server {
      bind_address = "0.0.0.0"
      bind_port = "8081"
      trust_domain = "test.com"
      data_dir = "/run/spire/data"
      log_level = "DEBUG"
      svid_ttl = "1h"
      upstream_bundle = true
      ca_subject = {
        Country = ["US"],
        Organization = ["SPIFFE"],
        CommonName = "",
      }
    }
    plugins {
      DataStore "sql" {
        plugin_data {
          database_type = "sqlite3"
          connection_string = "/run/spire/data/datastore.sqlite3"
        }
      }
      NodeAttestor "k8s_sat" {
        plugin_data {
          clusters = {
            # NOTE: Change this to your cluster name
            "kubernetes" = {
              use_token_review_api_validation = true
              service_account_whitelist = ["{{ .Release.Namespace }}:spire-agent"]
            }
          }
        }
      }
      NodeResolver "noop" {
        plugin_data {}
      }
      KeyManager "disk" {
        plugin_data {
          keys_path = "/run/spire/data/keys.json"
        }
      }
      Notifier "k8sbundle" {
        plugin_data {
          # This plugin updates the bundle.crt value in the spire:spire-bundle
          # ConfigMap by default, so no additional configuration is necessary.
        }
      }
    }
