helm repo add bitnami https://charts.bitnami.com/bitnami
helm install user-list-db bitnami/postgresql
# export POSTGRES_PASSWORD=$(kubectl get secret --namespace default user-list-db-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)