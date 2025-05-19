#!/bin/bash

# Install Java (Minecraft requires Java 17+)
sudo apt-get update
sudo apt-get install -y openjdk-21-jdk

# Create directories for server instances
mkdir -p servers/{user_id}/{server_id}  # Replace {user_id} and {server_id} dynamically in code
mkdir -p backups

# Download Minecraft server.jar (example: Vanilla 1.20.1)
wget https://piston-data.mojang.com/v1/objects/e6ec2f64e6080b9b5d9b471b291c33cc7f509733/server.jar \
     -O servers/_template/server.jar

# Set permissions (restrict access to server directories)
chmod -R 700 servers/