provider "kubernetes" {
  config_path = "~/.kube/config"
  config_context = "microk8s"
}

resource "kubernetes_namespace" "namespace" {
  metadata {
    name = "discordle"
  }
}

data "kubernetes_secret" "discord_token" {
  metadata {
    name = "discord_token"
    namespace = kubernetes_namespace.namespace.id
  }
}

resource "kubernetes_stateful_set" "discordle" {
  metadata {
    name = "discordle"
    namespace = kubernetes_namespace.namespace.id
    labels = {
      app = "discordle"
    }
  }
  spec {
    service_name = "discordle"
    selector {
      match_labels = {
        app = "discordle"
      }
    }
    template {
      metadata {
        labels = {
          app = "discordle"
        }
      }

      spec {
        container {
          name = "discordle"
          image = "ghcr.io/sodle/discordle:main"
          env_from {
            secret_ref {
              name = data.kubernetes_secret.discord_token.id
            }
          }
        }
      }
    }
  }
}
