#!/bin/bash

set -e

# Go to the project root directory
cd "$(dirname "$0")/.."

echo "Moving Authpost proto files..."
mkdir -p pkg/types/proto/pb/authpost
cp -r pkg/types/proto/pb/authpost/github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/authpost/* pkg/types/proto/pb/authpost/

echo "Moving Newsfeed proto files..."
mkdir -p pkg/types/proto/pb/newsfeed
cp -r pkg/types/proto/pb/newsfeed/github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed/* pkg/types/proto/pb/newsfeed/

echo "Moving Newsfeed Publishing proto files..."
mkdir -p pkg/types/proto/pb/newsfeed_publishing
cp -r pkg/types/proto/pb/newsfeed_publishing/github.com/hoangNguyenDev3/WanderSphere/backend/pkg/types/proto/pb/newsfeed_publishing/* pkg/types/proto/pb/newsfeed_publishing/

# Clean up the deep directory structure
rm -rf pkg/types/proto/pb/authpost/github.com
rm -rf pkg/types/proto/pb/newsfeed/github.com
rm -rf pkg/types/proto/pb/newsfeed_publishing/github.com

echo "All proto files moved successfully!" 