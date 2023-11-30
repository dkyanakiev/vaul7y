#!/bin/bash

# Set Vault address
export VAULT_ADDR='http://127.0.0.1:8200'

# Login to Vault
vault login $VAULT_TOKEN

# # Create multiple KV v2 stores with random names
# for i in {1..10}
# do
#   kv_store_name="kv$(uuidgen | cut -c1-8)"
#   vault secrets enable -version=2 -path=$kv_store_name kv
# done

# Create random secrets in each KV store
for kv_store_name in $(vault secrets list -format=json | jq -r 'to_entries[] | select(.value.type == "kv") | .key')
do
  for j in {1..100}
  do
    # Create a secret at the root of the KV store
    vault kv put $kv_store_name/data/secret$j key1=$(openssl rand -base64 12) key2=$(openssl rand -base64 12)

    # Create a nested object in the KV store
    nested_object_name="object$(uuidgen | cut -c1-8)"
    for k in {1..5}
    do
      # Create a secret in the nested object
      vault kv put $kv_store_name/data/$nested_object_name/secret$k key1=$(openssl rand -base64 12) key2=$(openssl rand -base64 12)
    done
  done
done

# Number of policies to create
num_policies=10

# Base name for the policies
policy_base_name="random_policy"

# Base path for the policies
policy_base_path="secret/data/random"

for ((i=1; i<=num_policies; i++)); do
    # Generate a random path
    random_path="$policy_base_path/$RANDOM"

    # Create a temporary file for the policy
    policy_file=$(mktemp)

    # Write the policy to the temporary file
    echo "path \"$random_path\" { capabilities = [\"read\"] }" > $policy_file

    # Create the policy in Vault
    vault policy write "$policy_base_name$i" $policy_file

    # Remove the temporary file
    rm $policy_file
done