#!/bin/sh
set -e

# mlcgo_mcp FULL ECOSYSTEM Installer
# This script installs ALL components:
# 1. Hub Servers (d2mcp, memory, osm, tasks)
# 2. Artifact Server & CLI
# 3. Wollmilchsau Execution Engine

echo "🌟 Starting Full mlcgo_mcp Ecosystem Installation..."

# 1. Install Hub
echo "--- Step 1: Installing Hub Servers ---"
curl -sfL https://raw.githubusercontent.com/hmsoft0815/mlcgo_mcp/main/scripts/install.sh | sh

# 2. Install Artifacts
echo "--- Step 2: Installing Artifact Store ---"
curl -sfL https://raw.githubusercontent.com/hmsoft0815/mlcartifact/main/scripts/install.sh | sh

# 3. Install Wollmilchsau
echo "--- Step 3: Installing Wollmilchsau ---"
curl -sfL https://raw.githubusercontent.com/hmsoft0815/wollmilchsau/main/scripts/install.sh | sh

echo "✨✨✨✨✨✨✨✨✨✨✨✨✨✨✨"
echo "✅ FULL ECOSYSTEM INSTALLED SUCCESSFULLY!"
echo "✨✨✨✨✨✨✨✨✨✨✨✨✨✨✨"
echo ""
echo "Installed binaries in /usr/local/bin:"
echo "- d2mcp, memory-server, openstreetmap_mcp, task-manager"
echo "- artifact-server, artifact-cli"
echo "- wollmilchsau"
echo ""
echo "🚀 Next step: Configure your Claude Desktop to use these tools."
